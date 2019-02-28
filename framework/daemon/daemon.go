package daemon

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func Daemon(nochdir, noclose int) (int, error) {
	if os.Getppid() == 1 {
		// chile process
		syscall.Umask(0)
		if nochdir == 0 {
			os.Chdir("/")
		}
		println("sub-process running")
		return 0, nil
	}

	cmd := exec.Command(os.Args[0])
	//files := make([]*os.File, 3, 6)
	if noclose == 0 {
		nullDev, err := os.OpenFile("/dev/null", 0, 0)
		if err != nil {
			return 1, err
		}
		cmd.Stdin = nullDev
		cmd.Stdout = nullDev
		cmd.Stderr = nullDev
		//files[0], files[1], files[2] = nullDev, nullDev, nullDev
	} else {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		//files[0], files[1], files[2] = os.Stdin, os.Stdout, os.Stderr
	}

	dir, _ := os.Getwd()
	sysAttrs := syscall.SysProcAttr{Setsid: true}
	cmd.SysProcAttr = &sysAttrs
	cmd.Env = os.Environ()
	cmd.Dir = dir

	//attrs := os.ProcAttr{Dir: dir, Env: os.Environ(), Files: files, Sys: &sysAttrs}

	pidFile, err := os.OpenFile("pid", os.O_RDWR|os.O_CREATE|os.O_SYNC, 0644)
	if err != nil {
		return -1, fmt.Errorf("create pid file error %s : %s", os.Args[0], err)
	}

	cmd.Start()
	//proc, err := os.StartProcess(os.Args[0], os.Args, &attrs)
	if err != nil {
		return -1, fmt.Errorf("create porcess error %s : %s", os.Args[0], err)
	}
	pid := fmt.Sprintf("%v", cmd.Process.Pid)
	pidFile.WriteString(pid)
	pidFile.Close()
	cmd.Process.Release()
	//proc.Release()
	os.Exit(0)
	return 0, nil
}
