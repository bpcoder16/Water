package middlewares

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bpcoder16/Water/env"
	"github.com/bpcoder16/Water/logit"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	LogIdKey = "logId"
	fuzzyStr = "***"
)

var needFilterMap = map[string]struct{}{
	"password": {},
	"token":    {},
}

// 自定义一个结构体，实现 gin.ResponseWriter interface
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func ApiLogger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		begin := time.Now()

		ctx.Set(logit.DefaultMessageKey, "HTTP")
		ctx.Set(LogIdKey, uuid.New().String())

		reqBody := generateRequestBody(ctx)

		writer := &responseWriter{
			ResponseWriter: ctx.Writer,
			body:           bytes.NewBuffer([]byte{}),
		}
		ctx.Writer = writer

		ctx.Next()

		elapsed := time.Since(begin)

		logit.Context(ctx).InfoW(
			"costTime", fmt.Sprintf("%.3fms", float64(elapsed.Nanoseconds())/1e6),
			"clientIP", ctx.ClientIP(),
			"method", ctx.Request.Method,
			"uri", ctx.Request.URL.Path,
			"rawQuery", ctx.Request.URL.RawQuery,
			"header", filterHeader(ctx.Request.Header),
			"reqBody", filterBody(reqBody),
			"statusCode", ctx.Writer.Status(),
			"response", filterBody(writer.body.Bytes()),
		)
	}
}

func WebsocketLogger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		begin := time.Now()

		ctx.Set(logit.DefaultMessageKey, "WebSocket")
		ctx.Set(LogIdKey, uuid.New().String())

		ctx.Next()

		elapsed := time.Since(begin)

		logit.Context(ctx).InfoW(
			"connDuration", fmt.Sprintf("%.3fms", float64(elapsed.Nanoseconds())/1e6),
			"clientIP", ctx.ClientIP(),
			"action", "connect",
			"header", filterHeader(ctx.Request.Header),
			"statusCode", ctx.Writer.Status(),
		)
	}
}

func filterHeader(header http.Header) http.Header {
	if env.RunMode() != env.RunModeRelease {
		return header
	}
	for key := range header {
		if _, ok := needFilterMap[strings.ToLower(key)]; ok {
			header.Set(key, fuzzyStr)
		}
	}
	return header
}

func filterBody(b []byte) interface{} {
	var reqBody map[string]interface{}
	err := json.Unmarshal(b, &reqBody)
	if err != nil {
		return string(b)
	}
	if env.RunMode() != env.RunModeRelease {
		return reqBody
	}
	for key := range reqBody {
		if _, ok := needFilterMap[strings.ToLower(key)]; ok {
			reqBody[key] = fuzzyStr
		}
	}
	return reqBody
}

func generateRequestBody(ctx *gin.Context) []byte {
	body, _ := ctx.GetRawData()                            // 读取 request body 的内容
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body)) // 创建 io.ReadCloser 对象传给 request body
	return body                                            // 返回 request body 的值
}
