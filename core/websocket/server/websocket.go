package server

import (
	"github.com/bpcoder16/Water/core/http/server"
	"github.com/bpcoder16/Water/logit"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  256,
	WriteBufferSize: 256,
	WriteBufferPool: &sync.Pool{},
}

// TODO 临时解决方案，由于 gorilla/websocket 不支持 Sec-WebSocket-Extensions Header
func filterHeader(h http.Header) http.Header {
	h.Del("Sec-WebSocket-Extensions")
	return h
}

//	{
//		"action": "test",
//		"payload": {
//			"key1": 1234,
//			"key2": "value"
//		}
//	}
type WebSocketRouter struct {
	path        string
	controllers map[string]Controller
	middlewares []gin.HandlerFunc
}

func NewWebSocketRouter(path string) *WebSocketRouter {
	return &WebSocketRouter{
		path:        path,
		controllers: make(map[string]Controller),
		middlewares: make([]gin.HandlerFunc, 0),
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

func (r *WebSocketRouter) handle(ctx *gin.Context) {
	ctx.Set(logit.DefaultMessageKey, "WebSocket")
	writer, req := ctx.Writer, ctx.Request
	conn, err := upgrader.Upgrade(writer, req, filterHeader(ctx.Request.Header))
	if err != nil {
		logit.Context(ctx).Warn("websocket upgrade fail:", err)
		return
	}
	defer func() { _ = conn.Close() }()

	for {
		mt, message, errR := conn.ReadMessage()
		if errR != nil {
			logit.Context(ctx).Warn("websocket read message fail:", errR)
			break
		}

		logit.Context(ctx).InfoW(
			"clientIP", ctx.ClientIP(),
			"uri", ctx.Request.URL.Path,
			"header", filterHeader(ctx.Request.Header),
			"action", "receive",
			"msg", string(message),
		)
		err = conn.WriteMessage(mt, message)
		logit.Context(ctx).InfoW(
			"clientIP", ctx.ClientIP(),
			"uri", ctx.Request.URL.Path,
			"header", filterHeader(ctx.Request.Header),
			"action", "send",
			"msg", string(message),
		)
		if err != nil {
			logit.Context(ctx).Warn("websocket write message fail:", err)
			break
		}
	}
}
