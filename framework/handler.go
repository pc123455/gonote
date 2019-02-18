package framework

func handler404(ctx *Context) {
	ctx.WriteHeader(404)
	ctx.ResponseWriter.Write([]byte("404 not found"))
}
