package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strconv"
)

func getPidFile() string {
	return projectFolder + "/.fileboy.pid"
}

func runAsDeamon() (int, error) {
	if runtime.GOOS == "windows" {
		logAndExit(PreError, "daemons mode cannot run on windows.")
	}
	err := stopDeamon()
	if err != nil {
		logAndExit(PreError, err)
	}
	_, err = exec.LookPath("fileboy")
	if err != nil {
		logAndExit(PreError, "cannot found `fileboy` command in the PATH")
	}
	deamon := exec.Command("fileboy")
	deamon.Dir = projectFolder
	deamon.Env = os.Environ()
	deamon.Stdout = os.Stdout
	err = deamon.Start()
	if err != nil {
		logAndExit(PreError, err)
	}
	pid := deamon.Process.Pid
	if pid != 0 {
		ioutil.WriteFile(getPidFile(), []byte(strconv.Itoa(pid)), 0644)
	}
	return pid, nil
}

func stopDeamon() error {
	bs, err := ioutil.ReadFile(getPidFile())
	if err != nil {
		return nil
	}
	_ = exec.Command("kill", string(bs)).Run()
	os.Remove(getPidFile())
	return nil
}
