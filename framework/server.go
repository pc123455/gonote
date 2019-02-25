package framework

import (
	"errors"
	"fmt"
	"go/types"
	"gonote/framework/context"
	"gonote/framework/logger"
	"gonote/framework/route"
	"net"
	"net/http"
	"path"
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

type Handler struct {
	handle func(writer http.ResponseWriter, request *http.Request)
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

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "*" {
		if r.ProtoAtLeast(1, 1) {
			w.Header().Set("Connection", "close")
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// CONNECT requests are not canonicalized.
	if r.Method == "CONNECT" {
		// If r.URL.Path is /tree and its handler is not registered,
		// the /tree -> /tree/ redirect applies to CONNECT requests
		// but the path canonicalization does not.
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

	// All other requests have any port stripped and path cleaned
	// before passing to mux.handler.

	//host := stripHostPort(r.Host)
	path := cleanPath(r.URL.Path)

	if path != r.URL.Path {
		//_, pattern = mux.handler(host, path)
		url := *r.URL
		url.Path = path
		//return http.RedirectHandler(url.String(), StatusMovedPermanently), pattern
	}

	var handler HandlerFunc = nil
	logger.Infof("%q %q %q ", r.Proto, r.Method, r.RequestURI)
	ctx := context.NewContext(w, r)

	defer func() {
		err := recover()
		switch err.(type) {
		case types.Nil:
		case context.HttpError:
			httpError := err.(context.HttpError)
			errHandler := h.errHandlerMap[httpError.Status]
			if errHandler == nil {
				errHandler = h.defaultErrorHandlerFunc
			}
			errHandler(ctx)
		default:
			logger.Errorf(string(debug.Stack()))
			ctx.Output.SetStatus(http.StatusInternalServerError)
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
			currentHandlers = h.preAccessHandlers
		case AccessStage:
			currentHandlers = h.accessHandlers
		case BeforeRouteStage:
			currentHandlers = h.beforeRouteHandlers
		case RouteStage:
			handler = h.handlerRouteFunc(ctx)
			currentStage++
			continue
		case BeforeContentProcessStage:
			currentHandlers = h.beforeContentHandlers
		case ContentProcessStage:
			handler(ctx)
			currentStage++
			continue
		case AfterContentProcessStage:
			currentHandlers = h.afterContentHandlers
			ctx.Output.Write()
		case LogStage:
			currentHandlers = h.logHandlers
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
}

func (h *Handler) Initialize() {

	//initialize error handler
	h.errHandlerMap = make(ErrorHandlerMap)
	h.defaultErrorHandlerFunc = handlerOtherError

	//initialize handlers
	h.preAccessHandlers = make([]HandlerFunc, 0)
	h.accessHandlers = make([]HandlerFunc, 0)
	h.beforeRouteHandlers = make([]HandlerFunc, 0)
	h.beforeContentHandlers = []HandlerFunc{handlerParseParamFunc}
	h.afterContentHandlers = make([]HandlerFunc, 0)
	h.logHandlers = make([]HandlerFunc, 0)
}

func (this *Handler) AppendFilterHandler(stage int, handler HandlerFunc) error {
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

func (this *Handler) handlerRouteFunc(ctx *context.Context) (handler HandlerFunc) {
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

func (this *Handler) AddHandleFunc(method, pattern string, handler HandlerFunc) {
	this.router.AddRoute(method, pattern, handler)
}

func stripHostPort(h string) string {
	// If no port on host, return unchanged
	if strings.IndexByte(h, ':') == -1 {
		return h
	}
	host, _, err := net.SplitHostPort(h)
	if err != nil {
		return h // on error, return unchanged
	}
	return host
}

// Return the canonical path for p, eliminating . and .. elements.
func cleanPath(p string) string {
	if p == "" {
		return "/"
	}
	if p[0] != '/' {
		p = "/" + p
	}
	np := path.Clean(p)
	// path.Clean removes trailing slash except for root;
	// put the trailing slash back if necessary.
	if p[len(p)-1] == '/' && np != "/" {
		np += "/"
	}
	return np
}

type ErrorHandlerMap map[int]HandlerFunc

type Server struct {
	BindIP  string
	Port    int
	server  *http.Server
	handler *Handler
}

func (this *Server) Get(pattern string, handler HandlerFunc) {
	this.handler.AddHandleFunc("GET", pattern, handler)
}

func (this *Server) Post(pattern string, handler HandlerFunc) {
	this.handler.AddHandleFunc("POST", pattern, handler)
}

func (this *Server) Put(pattern string, handler HandlerFunc) {
	this.handler.AddHandleFunc("Put", pattern, handler)
}

func (this *Server) Delete(pattern string, handler HandlerFunc) {
	this.handler.AddHandleFunc("Delete", pattern, handler)
}

func (this *Server) Initialize(ip string, port int) {
	//this.router = route.NewBaseRouter()

	this.handler = &Handler{
		router: route.NewBaseRouter(),
	}

	this.server = &http.Server{
		Addr:           fmt.Sprintf(":%v", port),
		Handler:        this.handler,
		MaxHeaderBytes: 1 << 30,
	}
}

func (this *Server) Run() {
	this.server.ListenAndServe()
}
