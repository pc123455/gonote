package framework

import (
	"fmt"
	"gonote/framework/context"
	"gonote/framework/logger"
	"gonote/framework/route"
	"gonote/framework/utils"
	"net/http"
	"runtime/debug"
	"strings"
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
	beforeRouteHandlerFunc  func(ctx *context.Context)
	contentHandlerFunc      func(ctx *context.Context)
	afterRequestHandlerFunc func(ctx *context.Context)
	logHandlerFunc          func(ctx *context.Context)
}

func queryParse(raw string) (param context.Param) {
	param = make(context.Param)
	if raw == "" {
		return
	}
	querylist := strings.Split(raw, "&")
	for _, q := range querylist {
		kvPair := strings.Split(q, "=")
		key := kvPair[0]
		value := ""
		if len(kvPair) > 1 {
			value = kvPair[1]
		}
		param[key] = value
	}
	return
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
		logger.Infof("%q %q %q ", request.Proto, request.Method, request.RequestURI)
		defer func() {
			if err := recover(); err != nil {
				logger.Errorf(string(debug.Stack()))
				writer.WriteHeader(http.StatusInternalServerError)
			}
		}()

		handler, param = this.router.MatchRoute(request.Method, request.URL.Path)
		if handler == nil {
			handler = handler404
		}

		queryParam := queryParse(request.URL.RawQuery)
		utils.Merge(param, queryParam)
		ctx := context.Context{
			Input: context.Request{
				Request:    request,
				Args:       param,
				RawContent: nil,
			},
			Output: context.Response{
				Writer: writer,
			},
		}

		contentType := request.Header.Get("Content-Type")
		if strings.ToLower(contentType) == "application/json" {
			request.Body.Read(ctx.Input.RawContent)
			//if n > 0 && err == nil {
			//	body := make(map[string]interface{})
			//	json.Unmarshal(rowContent, body)
			//
			//	//merge args
			//	utils.Merge(ctx.Input.Args, body)
			//}
		}

		handler(&ctx)
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
