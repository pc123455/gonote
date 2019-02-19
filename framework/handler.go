package framework

import "gonote/framework/context"

func handler404(ctx *context.Context) {
	ctx.WriteHeader(404)
	ctx.ResponseWriter.Write([]byte("404 not found"))
}
