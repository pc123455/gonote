package framework

import (
	"fmt"
	"go/types"
	"gonote/framework/context"
	"gonote/framework/logger"
	"gonote/framework/route"
	"net/http"
	"runtime/debug"
	"strings"
)

type HandlerFunc func(ctx *context.Context)

type ErrorHandlerMap map[int]HandlerFunc

type Server struct {
	BindIP string
	Port   int
	server *http.Server
	router route.Router

	errHandlerMap           ErrorHandlerMap
	defaultErrorHandlerFunc HandlerFunc

	configReadHandlerFunc   HandlerFunc
	preAccessHandlerFunc    HandlerFunc
	accessHandlerFunc       HandlerFunc
	postAccessHandlerFunc   HandlerFunc
	beforeRouteHandlerFunc  HandlerFunc
	contentHandlerFunc      HandlerFunc
	afterRequestHandlerFunc HandlerFunc
	logHandlerFunc          HandlerFunc
}

func (this *Server) Initialize(ip string, port int) {
	this.router = route.NewBaseRouter()
	this.server = &http.Server{
		Addr:           fmt.Sprintf(":%v", port),
		MaxHeaderBytes: 1 << 30,
	}
	this.errHandlerMap = make(ErrorHandlerMap)
	this.defaultErrorHandlerFunc = handlerOtherError
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		var handler HandlerFunc = nil
		logger.Infof("%q %q %q ", request.Proto, request.Method, request.RequestURI)
		defer func() {
			err := recover()
			switch err.(type) {
			case types.Nil:
			case context.HttpError:
				httpError := err.(context.HttpError)
				errHandler := this.errHandlerMap[httpError.Status]
				if errHandler == nil {
					errHandler = this.defaultErrorHandlerFunc
				}
				ctx := httpError.GetContext()
				if ctx != nil {
					ctx.Output.Write()
				}
			default:
				logger.Errorf(string(debug.Stack()))
				writer.WriteHeader(http.StatusInternalServerError)
			}
		}()

		ctx := context.Context{
			Input: context.Request{
				Request:    request,
				Args:       make(context.Param),
				RawContent: nil,
			},
			Output: context.Response{
				Writer: writer,
			},
		}
		//config phase
		if this.configReadHandlerFunc != nil {
			this.configReadHandlerFunc(&ctx)
		}

		handler = this.handlerRouteFunc(&ctx)

		this.handlerParseParamFunc(&ctx)

		handler(&ctx)

		ctx.Output.Write()
	})
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

func (this *Server) handlerParseParamFunc(ctx *context.Context) {
	queryParam := queryParse(ctx.Input.URL.RawQuery)
	context.MergeParam(ctx.Input.Args, queryParam)
	contentType := ctx.Input.Header.Get("Content-Type")
	if strings.ToLower(contentType) == "application/json" {
		ctx.Input.Body.Read(ctx.Input.RawContent)
		//if n > 0 && err == nil {
		//	body := make(map[string]interface{})
		//	json.Unmarshal(rowContent, body)
		//
		//	//merge args
		//	utils.Merge(ctx.Input.Args, body)
		//}
	}
}

func (this *Server) handlerRouteFunc(ctx *context.Context) (handler HandlerFunc) {
	handler, param := this.router.MatchRoute(ctx.Input.Method, ctx.Input.URL.Path)
	if handler == nil {
		ctx.Abort(context.HttpError{
			Status:  http.StatusNotFound,
			Message: "not found",
		})
	}
	if param != nil {
		context.MergeParam(ctx.Input.Args, param)
	}
	return
}

func (this *Server) Get(pattern string, handler HandlerFunc) {
	this.router.AddRoute("GET", pattern, handler)
}

func (this *Server) Post(pattern string, handler HandlerFunc) {
	this.router.AddRoute("POST", pattern, handler)
}

func (this *Server) Put(pattern string, handler HandlerFunc) {
	this.router.AddRoute("Put", pattern, handler)
}

func (this *Server) Delete(pattern string, handler HandlerFunc) {
	this.router.AddRoute("Delete", pattern, handler)
}

func (this *Server) Run() {
	this.server.ListenAndServe()
}
