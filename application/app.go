package application

import (
	"note/config"
	"note/framework"
	"note/logger"
)

var (
	Server = framework.Server{}
	Config *config.Config
)

func Initialize(confFile string) {
	Config = config.ParseConfigFromFile(confFile)
	logger.Initialize(Config.Log.File, Config.Log.Level)
	Server = framework.Server{}
	Server.Initialize(Config.Net.Bind, Config.Net.Port)
}

func Run() {
	Server.Run()
}
