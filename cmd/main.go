package cmd

import (
	app "gonote/internal/app"
	"gonote/web"
)

func main() {
	app.Main()
	server := &app.Server

	server.Post("/create", web.Create)
	server.Put("/fix/<uuid>", web.Update)
	server.Delete("/del/<uuid>", web.Delete)
	server.Get("/get/find", web.Get)

	app.Run()
}
