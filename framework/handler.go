package framework

import (
	"gonote/framework/context"
	"net/http"
)

func handlerBadRequest(ctx *context.Context) {
	ctx.Output.SetStatus(http.StatusBadRequest)
	ctx.Output.AppendContent([]byte("400 bad request"))
}

func handlerUnauthorized(ctx *context.Context) {
	ctx.Output.SetStatus(http.StatusUnauthorized)
	ctx.Output.AppendContent([]byte("401 unauthorized"))
}

func handlerForbidden(ctx *context.Context) {
	ctx.Output.SetStatus(http.StatusForbidden)
	ctx.Output.AppendContent([]byte("403 forbidden"))
}

func handlerNotFound(ctx *context.Context) {
	ctx.Output.SetStatus(http.StatusNotFound)
	ctx.Output.AppendContent([]byte("404 not found"))
}

func handlerOtherError(ctx *context.Context) {
	ctx.Output.SetStatus(ctx.Output.Error.Status)
	ctx.Output.AppendContent([]byte(ctx.Output.Error.Message))
}
