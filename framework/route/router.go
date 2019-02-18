package route

import (
	"gonote/framework"
	"regexp"
	"strings"
)

type Param map[string]interface{}

type Router interface {
	AddRoute(method string, pattern string, handler func(ctx *framework.Context)) error
	MatchRoute(method string, path string) (handler func(ctx *framework.Context), params Param)
	Initialize()
}

type baseRouter struct {
	fixRouter      Router
	variableRouter Router
}

func NewBaseRouter() *baseRouter {
	return new(baseRouter)
}

func (this *baseRouter) AddRoute(method string, pattern string, handler func(ctx *framework.Context)) error {
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

func (this *baseRouter) MatchRoute(method string, path string) (handler func(ctx *framework.Context)) {
	method = strings.ToUpper(strings.TrimSpace(method))

	// first, find handler from fix routes
	handler, _ = this.fixRouter.MatchRoute(method, path)
	if handler != nil {
		return handler
	}

	// second, find handler from variable route
	pathSequence := strings.Split(path, "/")
	currentNode := rootNode

	handler, ok := pathMap[path]
	if ok {
		return handler
	}

	return nil
}

func (this *baseRouter) Initialize() {
	this.fixRouter = new(FixRouter)
	this.variableRouter = new(VariableRouter)
}
