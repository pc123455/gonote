package main

import (
	"gonote/application"
	"gonote/web"
)

func main() {
	application.Main()
	server := &application.Server

	server.Post("/create", web.Create)
	server.Put("/fix/<uuid>", web.Update)
	server.Delete("/del/<uuid>", web.Delete)
	server.Get("/get/find", web.Get)

	application.Run()
}
