package http

import (
	"github.com/bpcoder16/Water/env"
	"github.com/bpcoder16/Water/middlewares"
	"github.com/gin-gonic/gin"
	"net/http"
)

func init() {
	loadHttpConfig()
	switch env.RunMode() {
	case env.RunModeRelease:
		gin.SetMode(gin.ReleaseMode)
	case env.RunModeTest:
		gin.SetMode(gin.TestMode)
	case env.RunModeDebug:
		gin.SetMode(gin.DebugMode)
	}
}

type Handler struct {
	*gin.Engine
}

func NewHandler() *Handler {
	h := &Handler{
		Engine: gin.New(),
	}
	h.Engine.Use(middlewares.Logger())
	return h
}

func (h *Handler) StartServer() error {
	return (&http.Server{
		Addr:    ":" + config.Server.Port,
		Handler: h,
	}).ListenAndServe()
}
