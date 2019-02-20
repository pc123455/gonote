package main

import (
	"fmt"
	"gonote/application"
	"gonote/framework/context"
	"gonote/web"
	"reflect"
)

func main() {
	server := &application.Server
	application.Initialize("config.yml")

	server.Post("/create", web.Create)
	server.Get("/asd/<float:id>", func(ctx *context.Context) {
		num := (*ctx.Param)["id"]
		fmt.Print(reflect.TypeOf(num))
	})
	server.Put("/fix/<uuid>", web.Update)
	server.Delete("/del/<uuid>", web.Delete)
	application.Run()

	//_, error := daemon.Daemon(0, 0)
	//if error != nil {
	//}

}
