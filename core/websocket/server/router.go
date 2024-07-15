package server

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/bpcoder16/Water/core/http/server"
	"github.com/bpcoder16/Water/env"
	"github.com/bpcoder16/Water/logit"
	"github.com/bpcoder16/Water/middlewares"
	"github.com/bpcoder16/Water/module/gtask"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"net/http"
	"reflect"
	"sync"
)

const WebSocketConnUUIDCTXKey = "WebSocketConnUUIDCTXKey"

var upgrader = websocket.Upgrader{
	ReadBufferSize:  256,
	WriteBufferSize: 256,
	WriteBufferPool: &sync.Pool{},
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// TODO 临时解决方案，由于 gorilla/websocket 不支持 Sec-WebSocket-Extensions Header
func filterHeader(h http.Header) http.Header {
	h.Del("Sec-WebSocket-Extensions")
	return h
}

type Task func(*Manager) (err error)

type WebSocketRouter struct {
	path                   string
	textMessageControllers map[string]TextMessageController
	middlewares            []gin.HandlerFunc
	authorization          func(ctx *gin.Context) bool
	clientCloseFunc        func(uuidStr string)
	mu                     sync.RWMutex
	ClientManager          *Manager
	tasks                  []Task
}

func (r *WebSocketRouter) GetTasks() []func(context.Context) func() error {
	tasks := make([]func(context.Context) func() error, 0)
	if env.RunMode() != env.RunModeRelease {
		tasks = append(tasks, func(_ context.Context) func() error {
			return func() error {
				return WebSocketMonitor(r.ClientManager)
			}
		})
	}
	for _, task := range r.tasks {
		tasks = append(tasks, func(_ context.Context) func() error {
			return func() error {
				return task(r.ClientManager)
			}
		})
	}
	return tasks
}

func NewWebSocketRouter(path string) *WebSocketRouter {
	return &WebSocketRouter{
		path:                   path,
		textMessageControllers: make(map[string]TextMessageController),
		middlewares: []gin.HandlerFunc{
			middlewares.WebsocketLogger(),
		},
		ClientManager: NewManager(),
		tasks:         make([]Task, 0),
	}
}

func (r *WebSocketRouter) SetAuthorization(f func(ctx *gin.Context) bool) {
	r.authorization = f
}

func (r *WebSocketRouter) SetClientCloseFunc(f func(uuidStr string)) {
	r.clientCloseFunc = f
}

func (r *WebSocketRouter) AddMiddleware(m gin.HandlerFunc) {
	r.middlewares = append(r.middlewares, m)
}

func (r *WebSocketRouter) RegisterHandler(s *server.Server) {
	if r.authorization != nil {
		r.middlewares = append(r.middlewares, func(ctx *gin.Context) {
			if !r.authorization(ctx) {
				ctx.String(http.StatusForbidden, "")
				ctx.Abort()
				return
			}
			ctx.Next()
		})
	}
	s.Engine.GET(r.path, append(r.middlewares, func(ctx *gin.Context) {
		r.handle(ctx)
	})...)
}

func (r *WebSocketRouter) AddTask(t Task) {
	r.tasks = append(r.tasks, t)
}

func (r *WebSocketRouter) OnTextMessageController(scene string, controller TextMessageController) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.textMessageControllers == nil {
		r.textMessageControllers = make(map[string]TextMessageController)
	}
	r.textMessageControllers[scene] = controller
}

func (r *WebSocketRouter) GetTextMessageController(scene string) (controller TextMessageController, err error) {
	var exist bool
	var controllerTemplate TextMessageController
	controllerTemplate, exist = r.textMessageControllers[scene]
	if !exist {
		err = errors.New("textMessageController not register")
		return
	}
	controller, _ = reflect.New(reflect.TypeOf(controllerTemplate).Elem()).Interface().(TextMessageController)
	controller.Init(controllerTemplate)
	return
}

func (r *WebSocketRouter) handle(ctx *gin.Context) {
	writer, req := ctx.Writer, ctx.Request
	conn, err := upgrader.Upgrade(writer, req, filterHeader(ctx.Request.Header))
	if err != nil {
		logit.Context(ctx).Warn("websocket upgrade fail:", err)
		return
	}

	uuidStr := ctx.GetString(WebSocketConnUUIDCTXKey)
	if len(uuidStr) == 0 {
		uuidStr = uuid.New().String()
	}

	client := NewClient(ctx, conn, r, uuidStr)
	client.manager = r.ClientManager
	client.manager.Store(uuidStr, client)
	defer func() {
		client.Close()
	}()

	var g *gtask.Group
	g, _ = gtask.WithContext(client.Ctx)

	g.Go(func() (err error) {
		client.writePump()
		return
	})
	g.Go(func() (err error) {
		client.readPump()
		return
	})

	_ = g.Wait()
}

func (r *WebSocketRouter) receiveTextMessage(c *Client, messageBytes []byte) (err error) {
	if len(messageBytes) == 0 {
		c.WarnLog("receive", websocket.TextMessage, messageBytes, errors.New("text message is empty"))
		return
	}

	var receiveMessage ReceiveMessage
	errJ := json.Unmarshal(messageBytes, &receiveMessage)
	if errJ != nil {
		c.WarnLog(
			"receive",
			websocket.TextMessage,
			messageBytes,
			errors.New("parse text message failed["+errJ.Error()+"]"),
		)
		return
	}

	c.Logger.WithContext(c.Ctx).DebugW("process", "parse text success", "receiveMessage", receiveMessage)
	if len(receiveMessage.Scene) == 0 {
		if len(c.State.Scene) == 0 {
			c.WarnLog(
				"receive",
				websocket.TextMessage,
				messageBytes,
				errors.New("receiveMessage.Scene is empty"),
			)
			return
		}
		// 如果没有传递，说明用户停留在当前场景
		receiveMessage.Scene = c.State.Scene
		receiveMessage.SceneParams = c.State.SceneParams
		receiveMessage.SID = c.State.SID
	}

	controller, errC := r.GetTextMessageController(receiveMessage.Scene)
	if errC != nil {
		c.WarnLog(
			"receive",
			websocket.TextMessage,
			messageBytes,
			errC,
		)
		return
	}
	errP := controller.ParsePayload(c, receiveMessage)
	if errP != nil {
		c.WarnLog(
			"receive",
			websocket.TextMessage,
			messageBytes,
			errP,
		)
		return
	}
	err = controller.Process()

	return
}

func (r *WebSocketRouter) receiveBinaryMessage(_ *Client, _ []byte) (err error) {
	return
}

// 已经被劫持，正常收不到
func (r *WebSocketRouter) receiveCloseMessage(_ *Client, _ []byte) (err error) {
	return
}

// 已经被劫持，正常收不到
func (r *WebSocketRouter) receivePingMessage(_ *Client, _ []byte) (err error) {
	return
}

// 已经被劫持，正常收不到
func (r *WebSocketRouter) receivePongMessage(_ *Client, _ []byte) (err error) {
	return
}
