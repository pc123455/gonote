package main

import (
	"gonote/application"
	"gonote/framework/context"
)

func main() {
	server := &application.Server
	application.Initialize("config.yml")

	server.Post("/asd", func(ctx *context.Context) {
		ctx.ResponseWriter.Write([]byte("asd"))
	})
	server.Get("/asd/qwe", func(ctx *context.Context) {
		ctx.ResponseWriter.Write([]byte("asdqwe"))
	})
	server.Get("/<id>/qwe", func(ctx *context.Context) {
		ctx.ResponseWriter.Write([]byte("id que"))
	})
	server.Get("/asd/<id>", func(ctx *context.Context) {
		ctx.ResponseWriter.Write([]byte("asd id"))
	})
	server.Get("/<id>/<name>", func(ctx *context.Context) {
		ctx.ResponseWriter.Write([]byte("id name"))
	})
	application.Run()

	//_, error := daemon.Daemon(0, 0)
	//if error != nil {
	//}

}
