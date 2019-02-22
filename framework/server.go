package framework

import (
	"errors"
	"fmt"
	"go/types"
	"gonote/framework/context"
	"gonote/framework/logger"
	"gonote/framework/route"
	"net/http"
	"runtime/debug"
	"strings"
)

const (
	//frequency limit, concurrency limit, etc.
	PreAccessStage = iota
	//authentication
	AccessStage
	//filter before route
	BeforeRouteStage
	//route can not be customized
	RouteStage
	//filter before execute handler
	BeforeContentProcessStage
	//content process can not be customized
	ContentProcessStage
	//filter after execute handler
	AfterContentProcessStage
	//log statistics info
	LogStage
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

	preAccessHandlers     []HandlerFunc
	accessHandlers        []HandlerFunc
	beforeRouteHandlers   []HandlerFunc
	beforeContentHandlers []HandlerFunc
	afterContentHandlers  []HandlerFunc
	logHandlers           []HandlerFunc
}

func (this *Server) Initialize(ip string, port int) {
	this.router = route.NewBaseRouter()
	this.server = &http.Server{
		Addr:           fmt.Sprintf(":%v", port),
		MaxHeaderBytes: 1 << 30,
	}

	//initialize error handler
	this.errHandlerMap = make(ErrorHandlerMap)
	this.defaultErrorHandlerFunc = handlerOtherError

	//initialize handlers
	this.preAccessHandlers = make([]HandlerFunc, 0)
	this.accessHandlers = make([]HandlerFunc, 0)
	this.beforeRouteHandlers = make([]HandlerFunc, 0)
	this.beforeContentHandlers = []HandlerFunc{handlerParseParamFunc}
	this.afterContentHandlers = make([]HandlerFunc, 0)
	this.logHandlers = make([]HandlerFunc, 0)

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		var handler HandlerFunc = nil
		logger.Infof("%q %q %q ", request.Proto, request.Method, request.RequestURI)
		ctx := context.NewContext(writer, request)
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
				errHandler(ctx)
			default:
				logger.Errorf(string(debug.Stack()))
				writer.WriteHeader(http.StatusInternalServerError)
			}
			ctx.Output.Write()
		}()

		currentStage := PreAccessStage
		handlerIndex := 0
		var currentHandlers []HandlerFunc
		//config phase
		for {
			switch currentStage {
			case PreAccessStage:
				currentHandlers = this.preAccessHandlers
			case AccessStage:
				currentHandlers = this.accessHandlers
			case BeforeRouteStage:
				currentHandlers = this.beforeRouteHandlers
			case RouteStage:
				handler = this.handlerRouteFunc(ctx)
				currentStage++
				continue
			case BeforeContentProcessStage:
				currentHandlers = this.beforeContentHandlers
			case ContentProcessStage:
				handler(ctx)
				currentStage++
				continue
			case AfterContentProcessStage:
				currentHandlers = this.afterContentHandlers
				ctx.Output.Write()
			case LogStage:
				currentHandlers = this.logHandlers
			default:
				break
			}

			length := len(currentHandlers)
			if handlerIndex < length {
				currentHandlers[handlerIndex](ctx)
				handlerIndex++
			} else {
				handlerIndex = 0
				currentStage++
			}
			//is current stage terminated
			if ctx.IsStageOver() {
				currentStage++
				handlerIndex = 0
				ctx.ResetStageOver()
			}
		}
	})
}

func (this *Server) AppendFilterHandler(stage int, handler HandlerFunc) error {
	switch stage {
	case PreAccessStage:
		this.preAccessHandlers = append(this.preAccessHandlers, handler)
	case AccessStage:
		this.accessHandlers = append(this.accessHandlers, handler)
	case BeforeRouteStage:
		this.beforeRouteHandlers = append(this.beforeRouteHandlers, handler)
	case BeforeContentProcessStage:
		this.beforeRouteHandlers = append(this.beforeContentHandlers, handler)
	case AfterContentProcessStage:
		this.afterContentHandlers = append(this.afterContentHandlers, handler)
	case LogStage:
		this.logHandlers = append(this.logHandlers, handler)
	default:
		return errors.New("stage wrong")
	}
	return nil
}

func queryParse(raw string) (param context.Param) {
	param = make(context.Param)
	if raw == "" {
		return
	}
	queryList := strings.Split(raw, "&")
	for _, q := range queryList {
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

func handlerParseParamFunc(ctx *context.Context) {
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
			Message: []byte("not found"),
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
