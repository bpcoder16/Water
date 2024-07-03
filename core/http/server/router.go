package server

import (
	"context"
	"github.com/bpcoder16/Water/conf"
	"github.com/bpcoder16/Water/logit"
	"github.com/bpcoder16/Water/middlewares"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

type Router interface {
	RegisterHandler(*Server)
	GetTasks() []func(context.Context) func() error
}

type Server struct {
	ctx    context.Context
	cancel context.CancelFunc

	appConfig *conf.AppConfig

	Engine      *gin.Engine
	routers     []Router
	middlewares []gin.HandlerFunc
	tasks       []func(context.Context) func() error

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
		appConfig: appConfig,
		Engine:    gin.New(),
		routers:   make([]Router, 0),
		middlewares: []gin.HandlerFunc{
			middlewares.RecoveryWithWriter(os.Stderr),
		},
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
	tasks := router.GetTasks()
	if len(tasks) > 0 {
		s.AddTasks(tasks...)
	}
}

func (s *Server) AddMiddleware(m gin.HandlerFunc) {
	s.middlewares = append(s.middlewares, m)
}

func (s *Server) AddTasks(f ...func(context.Context) func() error) {
	s.tasks = append(s.tasks, f...)
}

func (s *Server) GetServerTasks(ctx context.Context) []func() error {
	tasks := make([]func() error, 0, len(s.tasks)+1)

	// 添加 HttpServer，异常抛出异常
	tasks = append(tasks, func() (err error) {
		s.ctx, s.cancel = context.WithCancel(ctx)

		if len(s.middlewares) > 0 {
			s.Engine.Use(s.middlewares...)
		}

		for _, router := range s.routers {
			router.RegisterHandler(s)
		}
		err = (&http.Server{
			Addr:    ":" + config.Server.Port,
			Handler: s.Engine,
		}).ListenAndServe()

		panic(err)
	})

	for _, task := range s.tasks {
		tasks = append(tasks, task(ctx))
	}

	return tasks
}
