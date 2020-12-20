// Copyright (c) 2018-2020 Author dengsgo<dengsgo@yoytang.com> [https://github.com/dengsgo/fileboy]
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

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

func runAsDaemon() (int, error) {
	if runtime.GOOS == "windows" {
		logAndExit("daemons mode cannot run on windows.")
	}
	err := stopDaemon()
	if err != nil {
		logAndExit(err)
	}
	_, err = exec.LookPath("fileboy")
	if err != nil {
		logAndExit("cannot found `fileboy` command in the PATH")
	}
	daemon := exec.Command("fileboy")
	daemon.Dir = projectFolder
	daemon.Env = os.Environ()
	daemon.Stdout = os.Stdout
	err = daemon.Start()
	if err != nil {
		logAndExit(err)
	}
	pid := daemon.Process.Pid
	if pid != 0 {
		ioutil.WriteFile(getPidFile(), []byte(strconv.Itoa(pid)), 0644)
	}
	return pid, nil
}

func stopDaemon() error {
	bs, err := ioutil.ReadFile(getPidFile())
	if err != nil {
		return nil
	}
	_ = exec.Command("kill", string(bs)).Run()
	os.Remove(getPidFile())
	return nil
}

func stopSelf() {
	pid := os.Getpid()
	os.Remove(getPidFile())
	_ = exec.Command("kill", strconv.Itoa(pid)).Run()
}
