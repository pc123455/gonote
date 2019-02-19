package application

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"gonote/config"
	"gonote/framework"
	"gonote/framework/logger"
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

func Run() {
	Server.Run()
}
