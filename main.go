package main

import (
	"gonote/application"
	"gonote/web"
)

func main() {
	server := &application.Server
	application.Initialize("config.yml")

	server.Post("/create", web.Create)
	application.Run()

	//_, error := daemon.Daemon(0, 0)
	//if error != nil {
	//}

}
