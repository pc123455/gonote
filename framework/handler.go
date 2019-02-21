package framework

import "gonote/framework/context"

func handler404(ctx *context.Context) {
	ctx.Output.SetStatus(404)
	ctx.Output.Write([]byte("404 not found"))
}
