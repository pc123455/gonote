package framework

import (
	"gonote/framework/context"
	"net/http"
)

func handlerBadRequest(ctx *context.Context) {
	ctx.Output.SetStatus(http.StatusBadRequest)
	ctx.Output.AppendContent(ctx.Output.Error.Message)
}

func handlerUnauthorized(ctx *context.Context) {
	ctx.Output.SetStatus(http.StatusUnauthorized)
	ctx.Output.AppendContent(ctx.Output.Error.Message)
}

func handlerForbidden(ctx *context.Context) {
	ctx.Output.SetStatus(http.StatusForbidden)
	ctx.Output.AppendContent(ctx.Output.Error.Message)
}

func handlerNotFound(ctx *context.Context) {
	ctx.Output.SetStatus(http.StatusNotFound)
	ctx.Output.AppendContent(ctx.Output.Error.Message)
}

func handlerOtherError(ctx *context.Context) {
	ctx.Output.SetStatus(ctx.Output.Error.Status)
	ctx.Output.AppendContent(ctx.Output.Error.Message)
}
