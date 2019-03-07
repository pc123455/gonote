// +build linux

package daemon

import (
	"errors"
	"fmt"
	"gonote/pkg/frameworkwork/logger"
	"os"
	"os/exec"
	"strconv"
	"syscall"
)

const (
	stageVar = "__DAEMON_STAGE"
)

const (
	stageParent = iota
	stageChild
)

func Daemon(nochdir, noclose int, pidFilename string) (int, error) {
	fmt.Println("sub-process running")
	id := os.Getpid()
	fmt.Printf("pid: %v\n", id)
	fmt.Printf("gid: %v\n", os.Getgid())
	fmt.Printf("ppid: %v\n", os.Getppid())
	stage, err := getStage()
	if err != nil {
		err = resetStage()
		if err != nil {
			return -1, err
		}
		stage = stageParent
	}
	if stage == stageChild {
		// chile process
		syscall.Umask(0)
		if nochdir == 0 {
			os.Chdir("/")
		}
		resetStage()
		fmt.Println("child precess start")
		//os.Exit(0)
		return 0, nil
	}

	err = os.Setenv(stageVar, strconv.Itoa(stageChild))
	if err != nil {
		fmt.Println("set environment variable error")
		return -1, err
	}
	cmd := exec.Command(os.Args[0])
	//files := make([]*os.File, 3, 6)
	nullDev, err := os.OpenFile("/dev/null", 0, 0)
	if err != nil {
		return 1, err
	}
	if noclose == 0 {
		cmd.Stdin = nullDev
		cmd.Stdout = nullDev
		cmd.Stderr = nullDev
		//files[0], files[1], files[2] = nullDev, nullDev, nullDev
	} else {
		cmd.Stdin = nullDev
		cmd.Stdout = logger.LogFile
		cmd.Stderr = logger.LogFile
		//files[0], files[1], files[2] = os.Stdin, os.Stdout, os.Stderr
	}

	dir, _ := os.Getwd()
	sysAttrs := syscall.SysProcAttr{Setsid: true}
	cmd.SysProcAttr = &sysAttrs
	err = os.Setenv(stageVar, strconv.Itoa(stageChild))
	if err != nil {
		return -1, fmt.Errorf("set enviornment error: %s", err)
	}
	cmd.Env = os.Environ()
	cmd.Dir = dir

	//attrs := os.ProcAttr{Dir: dir, Env: os.Environ(), Files: files, Sys: &sysAttrs}

	pidFile, err := os.OpenFile(pidFilename, os.O_RDWR|os.O_CREATE|os.O_SYNC, 0644)
	if err != nil {
		return -1, fmt.Errorf("create pid file error %s : %s", os.Args[0], err)
	}

	err = cmd.Start()
	//proc, err := os.StartProcess(os.Args[0], os.Args, &attrs)
	if err != nil {
		resetStage()
		return -1, fmt.Errorf("create porcess error %s : %s", os.Args[0], err)
	}
	pid := strconv.Itoa(cmd.Process.Pid)
	pidFile.WriteString(pid)
	pidFile.Close()
	cmd.Process.Release()
	//proc.Release()
	//time.Sleep(2 * time.Second)
	os.Exit(0)
	return 0, nil
}

func getStage() (int, error) {
	stageStr := os.Getenv(stageVar)
	if stageStr == "" {
		return -1, errors.New("stage is not set")
	}
	stage, err := strconv.Atoi(stageStr)
	if err != nil {
		return -1, err
	}
	if stage > stageChild {
		return -1, errors.New("stage is invalid")
	}
	return stage, nil
}

func resetStage() error {
	return os.Setenv(stageVar, strconv.Itoa(stageParent))
}
