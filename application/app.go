package application

import (
	"database/sql"
	"flag"
	_ "github.com/go-sql-driver/mysql"
	"gonote/config"
	"gonote/framework"
	"gonote/framework/logger"
	"os"
	"strconv"
	"syscall"
)

var (
	Server = framework.Server{}
	Config *config.Config
	Db     *sql.DB

	Args struct {
		stop   bool
		reload bool
		daemon bool
	}
)

func Initialize(confFile string) (err error) {
	flag.BoolVar(&Args.stop, "stop", false, "stop the daemon server gracefully")
	flag.BoolVar(&Args.reload, "reload", false, "reload the config and restart server")
	flag.BoolVar(&Args.daemon, "daemon", false, "run server in daemon mode")
	flag.Parse()

	Config = config.ParseConfigFromFile(confFile)
	logger.Initialize(Config.Log.File, Config.Log.Level)
	Server = framework.Server{}
	Server.Initialize(Config.Net.Bind, Config.Net.Port)

	Db, err = sql.Open("mysql", Config.Mysql.Uri)
	return
}

func gracefullyShutdown() {
	Server.Shutdown()
}

func Run() {
	go Server.Run()
	sigChan := make(chan os.Signal)
	select {
	case <-Server.GetDoneChan():
	case signal := <-sigChan:
		switch signal {
		case syscall.SIGTERM:
			gracefullyShutdown()
		default:

		}
	}
	print("server stop")
}

func readDaemonPid() (int, error) {
	pidFile, err := os.Open("pid")
	if err != nil {
		return 0, err
	}
	var buff []byte
	_, err = pidFile.Read(buff)
	if err != nil {
		return 0, err
	}
	pid, err := strconv.Atoi(string(buff))
	if err != nil {
		return 0, err
	}
	return pid, nil
}

func SignalNotify(signal syscall.Signal) error {
	_, err := readDaemonPid()
	if err != nil {
		return err
	}
	//syscall.Kill(pid, signal)
	return nil
}
