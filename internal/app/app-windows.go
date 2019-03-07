// +build windows

package app

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gonote/config"
	"gonote/pkg"
	_ "gonote/pkg"
	"gonote/pkg/daemon"
	"gonote/pkg/logger"
	"os"
	"strconv"
	"syscall"
)

var (
	Server = pkg.Server{}
	Config *config.Config
	Db     *sql.DB

	Args struct {
		stop       bool
		reload     bool
		daemon     bool
		configFile string
	}
)

func Main() error {
	flag.BoolVar(&Args.stop, "stop", false, "stop the daemon server gracefully")
	flag.BoolVar(&Args.reload, "reload", false, "reload the config and restart server")
	flag.BoolVar(&Args.daemon, "daemon", false, "run server in daemon mode")
	flag.StringVar(&Args.configFile, "config", "config.yml", "config file pat, example '/tmp/config.yml'")
	flag.Parse()

	if Args.stop {
		GracefullyStopDaemon()
		os.Exit(0)
	}

	if Args.reload {
		reloadDaemon()
		os.Exit(0)
	}

	err := Initialize()
	if err != nil {
		return err
	}

	return nil
}

func Initialize() (err error) {

	Config = config.ParseConfigFromFile(Args.configFile)
	logger.Init()
	logger.SetOutputFile(Config.Log.File)
	logger.SetLevel(Config.Log.Level)

	if Args.daemon {
		_, err = daemon.Daemon(0, 1, Config.Base.Pid)

		if err != nil {
			return err
		}
	}

	Server = pkg.Server{}
	Server.Initialize(Config.Net.Bind, Config.Net.Port)

	Db, err = sql.Open("mysql", Config.Mysql.Uri)
	return
}

func Stop() error {
	err := Db.Close()
	if err != nil {
		return fmt.Errorf("db connection close error: %s", err)
	}
	logger.Close()
	return nil
}

func gracefullyShutdown() {
	Server.Shutdown()
}

func reload() error {
	err := Stop()
	if err != nil {
		return err
	}
	err = Initialize()
	return err
}

func Run() {

	go Server.Run()
	sigChan := make(chan os.Signal)
Loop:
	for {
		select {
		case <-Server.GetDoneChan():
			break Loop
		case signal := <-sigChan:
			switch signal {
			case syscall.SIGTERM:
				gracefullyShutdown()

			//todo windows reload
			//case syscall.SIGUSR1:
			//err := reload()
			//if err != nil {
			//	fmt.Println(err)
			//	os.Exit(1)
			//}
			default:

			}
		}
	}
	fmt.Print("server stop")
}

func readDaemonPid() (int, error) {
	if Config.Base.Pid == "" {
		Config.Base.Pid = "pid"
	}
	pidFile, err := os.Open(Config.Base.Pid)
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

func SignalDaemon(signal syscall.Signal) error {
	//todo send signal to daemon
	//pid, err := readDaemonPid()
	//if err != nil {
	//	return err
	//}
	//err = syscall.(pid, signal)
	return errors.New("not implemented")
}

func GracefullyStopDaemon() error {
	err := SignalDaemon(syscall.SIGTERM)
	if err != nil {
		return fmt.Errorf("stop error: %s", err)
	}
	return nil
}

func reloadDaemon() error {
	//todo reload daemon
	//err := SignalDaemon(syscall.SIGUSR1)
	//if err != nil {
	//	return fmt.Errorf("stop error: %s", err)
	//}
	return errors.New("not implemented")
}
