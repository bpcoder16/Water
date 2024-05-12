package server

import (
	"context"
	"github.com/bpcoder16/Water/conf"
	"github.com/bpcoder16/Water/logit"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Server struct {
	ctx    context.Context
	cancel context.CancelFunc

	appConfig *conf.AppConfig

	Engine      *gin.Engine
	routers     []Router
	middlewares []gin.HandlerFunc

	logger *logit.Helper
}

var (
	DefaultServer *Server
)

func SetDefaultServer(s *Server) {
	DefaultServer = s
}

func NewServer(_ context.Context, appConfig *conf.AppConfig) *Server {
	s := &Server{
		appConfig:   appConfig,
		Engine:      gin.New(),
		routers:     make([]Router, 0),
		middlewares: make([]gin.HandlerFunc, 0),
	}

	if DefaultServer == nil {
		SetDefaultServer(s)
	}

	return s
}

func (s *Server) AppName() string {
	return s.appConfig.Env.AppName
}

func (s *Server) RunMode() string {
	return s.appConfig.Env.RunMode
}

func (s *Server) AddRouter(router Router) {
	s.routers = append(s.routers, router)
}

func (s *Server) AddMiddleware(m gin.HandlerFunc) {
	s.middlewares = append(s.middlewares, m)
}

func (s *Server) Start(ctx context.Context) error {
	s.ctx, s.cancel = context.WithCancel(ctx)

	if len(s.middlewares) > 0 {
		s.Engine.Use(s.middlewares...)
	}

	for _, router := range s.routers {
		router.RegisterHandler(s)
	}

	return (&http.Server{
		Addr:    ":" + config.Server.Port,
		Handler: s.Engine,
	}).ListenAndServe()
}
