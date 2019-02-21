package main

import (
	"gonote/application"
	"gonote/web"
)

func main() {
	server := &application.Server
	application.Initialize("config.yml")

	server.Post("/create", web.Create)
	server.Put("/fix/<uuid>", web.Update)
	server.Delete("/del/<uuid>", web.Delete)
	server.Get("/get/find", web.Get)
	application.Run()

	//_, error := daemon.Daemon(0, 0)
	//if error != nil {
	//}

}
