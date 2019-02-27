package route

import (
	"gonote/framework/context"
	"gonote/framework/logger"
)

type fixRoute map[string]func(ctx *context.Context)

type FixRouter struct {
	fixRouteMethodMap map[string]fixRoute
}

func (this *FixRouter) AddRoute(method string, pattern string, handler func(ctx *context.Context)) error {
	var route fixRoute
	route, ok := this.fixRouteMethodMap[method]
	if !ok {
		route = make(fixRoute, 10)
		this.fixRouteMethodMap[method] = route
	}

	if h := route[pattern]; h == nil {
		route[pattern] = handler
	} else {
		logger.Warnf("method: %s and path: %s already exist", method, pattern)
	}
	return nil
}

func (this *FixRouter) MatchRoute(method string, path string) (handler func(ctx *context.Context), param context.Param) {
	fixRoute, ok := this.fixRouteMethodMap[method]
	param = context.Param{}
	if ok {
		handler = fixRoute[path]
		if handler != nil {
			return
		}
	}
	return
}

func (this *FixRouter) Initialize() {
	this.fixRouteMethodMap = make(map[string]fixRoute, 6)
}
