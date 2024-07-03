package server

import (
	"context"
	"github.com/bpcoder16/Water/middlewares"
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
	"path"
)

const abortIndex int8 = math.MaxInt8 >> 1

type apiRouterRegistry struct {
	method  string
	path    string
	handler gin.HandlersChain
}

type ApiRouter struct {
	RouterGroup

	registries  []apiRouterRegistry
	middlewares []gin.HandlerFunc
}

func (r *ApiRouter) AddMiddleware(m gin.HandlerFunc) {
	r.middlewares = append(r.middlewares, m)
}

func (r *ApiRouter) RegisterHandler(s *Server) {
	apiGroup := s.Engine.Group("/api", r.middlewares...)

	for _, registry := range r.registries {
		switch registry.method {
		case http.MethodPost, http.MethodGet, http.MethodDelete, http.MethodPatch, http.MethodPut, http.MethodHead, http.MethodOptions:
			apiGroup.Handle(registry.method, registry.path, registry.handler...)
		}
	}
}

// GetTasks TODO 后续添加实现
func (r *ApiRouter) GetTasks() []func(context.Context) func() error {
	return []func(context.Context) func() error{}
}

func NewApiRouter() *ApiRouter {
	r := &ApiRouter{
		RouterGroup: RouterGroup{
			Handlers: nil,
			basePath: "/",
		},
		registries: make([]apiRouterRegistry, 0),
		middlewares: []gin.HandlerFunc{
			middlewares.ApiLogger(),
		},
	}
	r.RouterGroup.apiRouter = r
	return r
}

func assert1(guard bool, text string) {
	if !guard {
		panic(text)
	}
}

func lastChar(str string) uint8 {
	if str == "" {
		panic("The length of the string can't be 0")
	}
	return str[len(str)-1]
}

func joinPaths(absolutePath, relativePath string) string {
	if relativePath == "" {
		return absolutePath
	}

	finalPath := path.Join(absolutePath, relativePath)
	if lastChar(relativePath) == '/' && lastChar(finalPath) != '/' {
		return finalPath + "/"
	}
	return finalPath
}

type RouterGroup struct {
	Handlers  gin.HandlersChain
	basePath  string
	apiRouter *ApiRouter
}

func (group *RouterGroup) combineHandlers(handlers gin.HandlersChain) gin.HandlersChain {
	finalSize := len(group.Handlers) + len(handlers)
	assert1(finalSize < int(abortIndex), "too many handlers")
	mergedHandlers := make(gin.HandlersChain, finalSize)
	copy(mergedHandlers, group.Handlers)
	copy(mergedHandlers[len(group.Handlers):], handlers)
	return mergedHandlers
}

func (group *RouterGroup) calculateAbsolutePath(relativePath string) string {
	return joinPaths(group.basePath, relativePath)
}

func (group *RouterGroup) Group(relativePath string, middlewares ...gin.HandlerFunc) *RouterGroup {
	return &RouterGroup{
		Handlers:  group.combineHandlers(middlewares),
		basePath:  group.calculateAbsolutePath(relativePath),
		apiRouter: group.apiRouter,
	}
}

func (group *RouterGroup) POST(relativePath string, handlerFunc ...gin.HandlerFunc) {
	group.apiRouter.on(http.MethodPost, group.calculateAbsolutePath(relativePath), group.combineHandlers(handlerFunc))
}

func (group *RouterGroup) GET(relativePath string, handlerFunc ...gin.HandlerFunc) {
	group.apiRouter.on(http.MethodGet, group.calculateAbsolutePath(relativePath), group.combineHandlers(handlerFunc))
}

func (group *RouterGroup) DELETE(relativePath string, handlerFunc ...gin.HandlerFunc) {
	group.apiRouter.on(http.MethodDelete, group.calculateAbsolutePath(relativePath), group.combineHandlers(handlerFunc))
}

func (group *RouterGroup) PATCH(relativePath string, handlerFunc ...gin.HandlerFunc) {
	group.apiRouter.on(http.MethodPatch, group.calculateAbsolutePath(relativePath), group.combineHandlers(handlerFunc))
}

func (group *RouterGroup) PUT(relativePath string, handlerFunc ...gin.HandlerFunc) {
	group.apiRouter.on(http.MethodPut, group.calculateAbsolutePath(relativePath), group.combineHandlers(handlerFunc))
}

func (group *RouterGroup) HEAD(relativePath string, handlerFunc ...gin.HandlerFunc) {
	group.apiRouter.on(http.MethodHead, group.calculateAbsolutePath(relativePath), group.combineHandlers(handlerFunc))
}

func (group *RouterGroup) OPTIONS(relativePath string, handlerFunc ...gin.HandlerFunc) {
	group.apiRouter.on(http.MethodOptions, group.calculateAbsolutePath(relativePath), group.combineHandlers(handlerFunc))
}

func (r *ApiRouter) on(method, path string, handler gin.HandlersChain) {
	r.registries = append(r.registries, apiRouterRegistry{
		method:  method,
		path:    path,
		handler: handler,
	})
}
