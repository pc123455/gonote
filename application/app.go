package application

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"gonote/config"
	"gonote/framework"
	"gonote/framework/logger"
	"os"
	"syscall"
	"time"
)

var (
	Server = framework.Server{}
	Config *config.Config
	Db     *sql.DB
)

func Initialize(confFile string) (err error) {
	Config = config.ParseConfigFromFile(confFile)
	logger.Initialize(Config.Log.File, Config.Log.Level)
	Server = framework.Server{}
	Server.Initialize(Config.Net.Bind, Config.Net.Port)

	Db, err = sql.Open("mysql", Config.Mysql.Uri)
	return
}

func delayShutdown() {
	time.Sleep(time.Second * 10)
	Server.Shutdown()
}

func Run() {
	go Server.Run()
	sigChan := make(chan os.Signal)
	for signal := range sigChan {
		switch signal {
		case syscall.SIGTERM:
			Server.Shutdown()
		default:

		}
	}
	go delayShutdown()
	Server.Wait()
	print("server stop")
}
