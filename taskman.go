package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"sync"
	"time"
)

type TaskMan struct {
	lastTaskId int64
	delay      int
	cmd        *exec.Cmd
	notifier   *NetNotifier
	putLock    sync.Mutex
	runLock    sync.Mutex
}

func newTaskMan(delay int, callUrl string) *TaskMan {
	return &TaskMan{
		delay:    delay,
		notifier: newNetNotifier(callUrl),
	}
}

func (t *TaskMan) Put(cf *changedFile) {
	if t.delay < 1 {
		t.preRun(cf)
		return
	}
	t.putLock.Lock()
	defer t.putLock.Unlock()
	t.lastTaskId = cf.Changed
	go func() {
		<-time.After(time.Millisecond * time.Duration(t.delay))
		if t.lastTaskId > cf.Changed {
			return
		}
		t.preRun(cf)
	}()
}

func (t *TaskMan) preRun(cf *changedFile) {
	if t.cmd != nil && t.cmd.Process != nil && t.cmd.ProcessState != nil && !t.cmd.ProcessState.Exited() {
		log.Println("stop old process ")
		if err := t.cmd.Process.Kill(); err != nil {
			log.Println(PreWarn, "stopped err, reason:", err)
		}
	}
	go t.run(cf)
}

func (t *TaskMan) run(cf *changedFile) {
	go t.notifier.Put(cf)
	t.runLock.Lock()
	defer t.runLock.Unlock()
	for i := 0; i < len(cfg.Command.Exec); i++ {
		carr := cmdParse2Array(cfg.Command.Exec[i], cf)
		log.Println("EXEC", carr)
		t.cmd = exec.Command(carr[0], carr[1:]...)
		//cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: syscall.CREATE_UNICODE_ENVIRONMENT}
		t.cmd.Stdin = os.Stdin
		//cmd.Stdout = os.Stdout
		t.cmd.Stderr = os.Stderr
		t.cmd.Dir = projectFolder
		t.cmd.Env = os.Environ()
		stdout, err := t.cmd.StdoutPipe()
		if err != nil {
			log.Println(PreError, err.Error())
			return
		}
		err = t.cmd.Start()
		if err != nil {
			log.Println(PreError, "run command", carr, "error. ", err)
		}
		reader := bufio.NewReader(stdout)
		for {
			line, err2 := reader.ReadString('\n')
			if err2 != nil || io.EOF == err2 {
				break
			}
			fmt.Print(line)
		}
		err = t.cmd.Wait()
		if err != nil {
			log.Println(PreWarn, "cmd wait err ", err)
			break
		}
		if t.cmd.Process != nil && t.cmd.ProcessState != nil && !t.cmd.ProcessState.Exited() {
			if err = t.cmd.Process.Kill(); err != nil {
				log.Println(PreWarn, "cmd cannot kill, reason:", err)
			}
		}
	}

	log.Println("end ")
}
