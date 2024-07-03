package server

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/bpcoder16/Water/logit"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"sync"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type Client struct {
	Ctx    *gin.Context
	router *WebSocketRouter

	Conn      *websocket.Conn
	textMsgCh chan []byte
	isClosed  bool
	manager   *Manager
	uuid      string

	// 客户端状态信息
	State State

	mu     sync.RWMutex
	Logger *logit.Helper
}

type State struct {
	Scene       string
	SceneParams map[string]interface{}
}

func NewClient(ctx *gin.Context, conn *websocket.Conn, r *WebSocketRouter, uuid string) *Client {
	c := &Client{
		Ctx:       ctx,
		router:    r,
		Conn:      conn,
		textMsgCh: make(chan []byte, 256),
		Logger: logit.Context(ctx).WithValues(
			"clientIP", ctx.ClientIP(),
			"header", ctx.Request.Header,
		),
		State: State{
			SceneParams: make(map[string]interface{}),
		},
		uuid: uuid,
	}
	return c
}

func (c *Client) log(level, action string, messageType int, message []byte, err error, keyValues ...interface{}) {
	newKeyValues := []interface{}{
		"action", action,
		"messageType", func() string {
			switch messageType {
			case websocket.BinaryMessage:
				return "Binary"
			case websocket.CloseMessage:
				return "Close"
			case websocket.PingMessage:
				return "Ping"
			case websocket.PongMessage:
				return "Pong"
			default:
				return "Text"
			}
		}(),
		"message", string(message),
		"err", err,
	}
	newKeyValues = append(keyValues, newKeyValues...)

	switch level {
	case "DEBUG":
		c.Logger.WithContext(c.Ctx).DebugW(newKeyValues...)
	case "INFO":
		c.Logger.WithContext(c.Ctx).InfoW(newKeyValues...)
	case "WARN":
		c.Logger.WithContext(c.Ctx).WarnW(newKeyValues...)
	case "ERROR":
		c.Logger.WithContext(c.Ctx).ErrorW(newKeyValues...)
	}
	return
}

func (c *Client) DebugLog(action string, messageType int, message []byte, err error, keyValues ...interface{}) {
	c.log("DEBUG", action, messageType, message, err, keyValues...)
}

func (c *Client) InfoLog(action string, messageType int, message []byte, err error, keyValues ...interface{}) {
	c.log("INFO", action, messageType, message, err, keyValues...)
}

func (c *Client) WarnLog(action string, messageType int, message []byte, err error, keyValues ...interface{}) {
	c.log("WARN", action, messageType, message, err, keyValues...)
}

func (c *Client) ErrorLog(action string, messageType int, message []byte, err error, keyValues ...interface{}) {
	c.log("ERROR", action, messageType, message, err, keyValues...)
}

func (c *Client) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if false == c.isClosed {
		_ = c.Conn.Close()
		c.isClosed = true
		close(c.textMsgCh)
		if c.router.clientCloseFunc != nil {
			c.router.clientCloseFunc(c.uuid)
		}
	}
	c.manager.Delete(c.uuid)
}

func (c *Client) ReadMessage() (messageType int, message []byte, err error) {
	messageType, message, err = c.Conn.ReadMessage()
	c.DebugLog("receive", messageType, message, err)
	if err != nil {
		c.Close()
	}
	return
}

func (c *Client) WriteTextMessage(message []byte) (err error) {
	defer func() {
		if r := recover(); r != nil {
			c.ErrorLog("sendTextMsgCh", websocket.TextMessage, message, errors.New("panic: textMsgCh is closed"))
			c.Close()
			err = errors.New("c.textMsgCh is closed")
		}
	}()
	if !c.isClosed {
		c.textMsgCh <- message
		c.DebugLog("sendWriteCh", websocket.TextMessage, message, nil)
	} else {
		err = errors.New("c.textMsgCh is closed")
	}

	return
}

func (c *Client) writeMessage(messageType int, message []byte) (err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.isClosed {
		_ = c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
		err = c.Conn.WriteMessage(messageType, message)
		c.DebugLog("send", messageType, message, err)
	}
	return
}

func (c *Client) sendPing() error {
	c.DebugLog("send", websocket.PingMessage, nil, nil)
	return c.Conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(writeWait))
}

func (c *Client) writePump() {
	// 维持心跳
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		if r := recover(); r != nil {
			c.Logger.WithContext(c.Ctx).ErrorW(
				"function", "client.writePump",
				"recover", r,
			)
		}
		ticker.Stop()
		c.Close()
	}()
	for {
		select {
		case message, ok := <-c.textMsgCh:
			if !ok {
				_ = c.Conn.WriteControl(websocket.CloseMessage, []byte{}, time.Now().Add(writeWait))
				return
			}

			c.mu.Lock()
			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			_, _ = w.Write(message)

			n := len(c.textMsgCh)
			for i := 0; i < n; i++ {
				_, _ = w.Write(newline)
				_, _ = w.Write(<-c.textMsgCh)
			}
			if errW := w.Close(); errW != nil {
				c.mu.Unlock()
				return
			}
			c.mu.Unlock()
		case <-ticker.C:
			if err := c.sendPing(); err != nil {
				return
			}
			// 鉴权与心跳一个频次校验
			if c.router.authorization != nil {
				if isOk := c.router.authorization(c.Ctx); !isOk {
					c.Close()
					return
				}
			}
		}
	}
}

func (c *Client) readPump() {
	defer func() {
		if r := recover(); r != nil {
			c.Logger.WithContext(c.Ctx).ErrorW(
				"function", "client.readPump",
				"recover", r,
			)
		}
		c.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	_ = c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) (err error) {
		c.DebugLog("receive", websocket.PongMessage, nil, nil)
		_ = c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return
	})

	for {
		mt, message, errR := c.ReadMessage()
		if errR != nil {
			if websocket.IsUnexpectedCloseError(errR, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.Logger.WithContext(c.Ctx).WarnW("WebsocketReadMsgFail", errR)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

		begin := time.Now()
		switch mt {
		case websocket.TextMessage:
			errR = c.router.receiveTextMessage(c, message)
		case websocket.BinaryMessage:
			errR = c.router.receiveBinaryMessage(c, message)
		case websocket.CloseMessage:
			errR = c.router.receiveCloseMessage(c, message)
		case websocket.PingMessage:
			errR = c.router.receivePingMessage(c, message)
		case websocket.PongMessage:
			errR = c.router.receivePongMessage(c, message)
		}

		elapsed := time.Since(begin)
		c.InfoLog("receive", mt, message, errR, "costTime", fmt.Sprintf("%.3fms", float64(elapsed.Nanoseconds())/1e6))

		if errR != nil {
			c.Logger.WithContext(c.Ctx).WarnW("WebsocketHandleFail", errR)
			break
		}
	}
}
