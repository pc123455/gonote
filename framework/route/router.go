package route

import (
	"gonote/framework/context"
	"regexp"
	"strings"
)

type Router interface {
	AddRoute(method string, pattern string, handler func(ctx *context.Context)) error
	MatchRoute(method string, path string) (handler func(ctx *context.Context), params context.Param)
	Initialize()
}

type baseRouter struct {
	fixRouter      Router
	variableRouter Router
}

func NewBaseRouter() (router *baseRouter) {
	router = new(baseRouter)
	router.Initialize()
	return
}

func (this *baseRouter) AddRoute(method string, pattern string, handler func(ctx *context.Context)) error {
	//regex := "(/[a-zA-Z0-9._~-])*(/<[a-zA-Z]+>)(/[a-zA-Z0-9._~-])*"
	pattern = strings.TrimSpace(pattern)
	method = strings.ToUpper(strings.TrimSpace(method))
	pathReg := regexp.MustCompile("/<[\\w:]+>/?")
	if !pathReg.Match([]byte(pattern)) {
		//"pattern" is a fix route
		this.fixRouter.AddRoute(method, pattern, handler)
	} else {
		//"pattern" is a variable route
		this.variableRouter.AddRoute(method, pattern, handler)
	}

	return nil
}

func (this *baseRouter) MatchRoute(method string, path string) (handler func(ctx *context.Context), param context.Param) {
	method = strings.ToUpper(strings.TrimSpace(method))
	// first, find handler from fix routes
	handler, param = this.fixRouter.MatchRoute(method, path)
	if handler != nil {
		return
	}

	// second, find handler from variable route
	handler, param = this.variableRouter.MatchRoute(method, path)

	return
}

func (this *baseRouter) Initialize() {
	this.fixRouter = new(FixRouter)
	this.fixRouter.Initialize()
	this.variableRouter = new(VariableRouter)
	this.variableRouter.Initialize()
}
