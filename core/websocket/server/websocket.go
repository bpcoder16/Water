package server

import (
	"encoding/json"
	"errors"
	"github.com/bpcoder16/Water/core/http/server"
	"github.com/bpcoder16/Water/logit"
	"github.com/bpcoder16/Water/middlewares"
	"github.com/bpcoder16/Water/module/gtask"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"reflect"
	"sync"
)

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

type WebSocketRouter struct {
	path                   string
	textMessageControllers map[string]TextMessageController
	middlewares            []gin.HandlerFunc
	mu                     sync.RWMutex
}

func NewWebSocketRouter(path string) *WebSocketRouter {
	return &WebSocketRouter{
		path:                   path,
		textMessageControllers: make(map[string]TextMessageController),
		middlewares: []gin.HandlerFunc{
			middlewares.WebsocketLogger(),
		},
	}
}

func (r *WebSocketRouter) AddMiddleware(m gin.HandlerFunc) {
	r.middlewares = append(r.middlewares, m)
}

func (r *WebSocketRouter) RegisterHandler(s *server.Server) {
	s.Engine.GET(r.path, append(r.middlewares, func(ctx *gin.Context) {
		r.handle(ctx)
	})...)
}

func (r *WebSocketRouter) OnTextMessageController(path string, controller TextMessageController) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.textMessageControllers == nil {
		r.textMessageControllers = make(map[string]TextMessageController)
	}
	r.textMessageControllers[path] = controller
}

func (r *WebSocketRouter) GetTextMessageController(path string) (controller TextMessageController, err error) {
	var exist bool
	var controllerTemplate TextMessageController
	controllerTemplate, exist = r.textMessageControllers[path]
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

	client := NewClient(ctx, conn, r)
	ClientManager().Store(client)
	defer func() {
		client.Close()
		ClientManager().Delete(client)
	}()

	var g *gtask.Group
	g, _ = gtask.WithContext(client.ctx)

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

	c.logger.WithContext(c.ctx).DebugW("process", "parse text success", "receiveMessage", receiveMessage)

	if len(receiveMessage.Path) == 0 {
		c.WarnLog(
			"receive",
			websocket.TextMessage,
			messageBytes,
			errors.New("receiveMessage.Path is empty"),
		)
		return
	}

	controller, errC := r.GetTextMessageController(receiveMessage.Path)
	if errC != nil {
		c.WarnLog(
			"receive",
			websocket.TextMessage,
			messageBytes,
			errC,
		)
		return
	}
	errP := controller.ParsePayload(c, receiveMessage.Payload)
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
