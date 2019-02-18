package framework

import (
	"fmt"
	"gonote/framework/logger"
	"gonote/framework/route"
	"net/http"
)

type Server struct {
	BindIP string
	Port   int
	server *http.Server
	router route.Router

	configReadHandlerFunc   func(ctx *Context)
	preAccessHandlerFunc    func(ctx *Context)
	accessHandlerFunc       func(ctx *Context)
	postAccessHandlerFunc   func(ctx *Context)
	contentHandlerFunc      func(ctx *Context)
	afterRequestHandlerFunc func(ctx *Context)
	logHandlerFunc          func(ctx *Context)
}

func (this *Server) Initialize(ip string, port int) {
	this.router = route.NewBaseRouter()
	this.server = &http.Server{
		Addr:           fmt.Sprintf(":%v", port),
		MaxHeaderBytes: 1 << 30,
	}
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		var handler func(ctx *Context) = nil
		logger.Infof("收到请求 %q", request.RequestURI)
		defer func() { println("结束处理") }()
		handler = this.router.MatchRoute(request.Method, request.URL.Path)
		if handler == nil {
			handler = handler404
		}

		requestCtx := Context{writer, request}
		handler(&requestCtx)
	})
}

func (this *Server) Get(pattern string, handler func(ctx *Context)) {
	this.router.AddRoute("GET", pattern, handler)
}

func (this *Server) Post(pattern string, handler func(ctx *Context)) {
	this.router.AddRoute("POST", pattern, handler)
}

func (this *Server) Put(pattern string, handler func(ctx *Context)) {
	this.router.AddRoute("Put", pattern, handler)
}

func (this *Server) Delete(pattern string, handler func(ctx *Context)) {
	this.router.AddRoute("Delete", pattern, handler)
}

func (this *Server) Run() {
	this.server.ListenAndServe()
}
