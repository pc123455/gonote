package main

import (
	"note/application"
)

func main() {
	application.Initialize("config.yml")
	application.Run()

	//_, error := daemon.Daemon(0, 0)
	//if error != nil {
	//}

}
