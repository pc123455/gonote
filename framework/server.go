package framework

import (
	"fmt"
	"gonote/framework/context"
	"gonote/framework/logger"
	"gonote/framework/route"
	"net/http"
)

type Server struct {
	BindIP string
	Port   int
	server *http.Server
	router route.Router

	configReadHandlerFunc   func(ctx *context.Context)
	preAccessHandlerFunc    func(ctx *context.Context)
	accessHandlerFunc       func(ctx *context.Context)
	postAccessHandlerFunc   func(ctx *context.Context)
	contentHandlerFunc      func(ctx *context.Context)
	afterRequestHandlerFunc func(ctx *context.Context)
	logHandlerFunc          func(ctx *context.Context)
}

func (this *Server) Initialize(ip string, port int) {
	this.router = route.NewBaseRouter()
	this.server = &http.Server{
		Addr:           fmt.Sprintf(":%v", port),
		MaxHeaderBytes: 1 << 30,
	}
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		var handler func(ctx *context.Context) = nil
		var param context.Param
		logger.Infof("收到请求 %q", request.RequestURI)
		defer func() { println("结束处理") }()
		handler, param = this.router.MatchRoute(request.Method, request.URL.Path)
		fmt.Print(param)
		if handler == nil {
			handler = handler404
		}

		requestCtx := context.Context{writer, request, &param}
		handler(&requestCtx)
	})
}

func (this *Server) Get(pattern string, handler func(ctx *context.Context)) {
	this.router.AddRoute("GET", pattern, handler)
}

func (this *Server) Post(pattern string, handler func(ctx *context.Context)) {
	this.router.AddRoute("POST", pattern, handler)
}

func (this *Server) Put(pattern string, handler func(ctx *context.Context)) {
	this.router.AddRoute("Put", pattern, handler)
}

func (this *Server) Delete(pattern string, handler func(ctx *context.Context)) {
	this.router.AddRoute("Delete", pattern, handler)
}

func (this *Server) Run() {
	this.server.ListenAndServe()
}
