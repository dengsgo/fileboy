package main

import (
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

	waitChan  chan bool
	waitQueue []*changedFile
}

func newTaskMan(delay int, callUrl string) *TaskMan {
	t := &TaskMan{
		delay:     delay,
		notifier:  newNetNotifier(callUrl),
		waitChan:  make(chan bool, 1),
		waitQueue: []*changedFile{},
	}
	if keyInInstruction(InstShouldFinish) {
		go func() {
			for {
				<-t.waitChan
				if len(t.waitQueue) > 0 {
					cf := t.waitQueue[len(t.waitQueue)-1]
					if len(t.waitQueue) > 1 {
						log.Println("Number of redundant tasks dropped:", len(t.waitQueue)-1)
					}
					t.waitQueue = []*changedFile{}
					go t.preRun(cf)
				}
			}
		}()
	}

	return t
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
		if keyInInstruction(InstShouldFinish) {
			t.waitQueue = append(t.waitQueue, cf)
			if t.cmd == nil {
				t.waitChan <- true
				return
			}
			log.Println("Waitting for the last task to finish")
			log.Println("Number of waiting tasks:", len(t.waitQueue))
		} else {
			t.preRun(cf)
		}
	}()
}

func (t *TaskMan) preRun(cf *changedFile) {
	if t.cmd != nil && t.cmd.Process != nil {
		if err := t.cmd.Process.Kill(); err != nil {
			log.Println("stop old process ")
			log.Println(PreWarn, "stopped err, reason:", err)
		}
	}
	go t.run(cf)
	go t.notifier.Put(cf)
}

func (t *TaskMan) waitFinish() {
	log.Println("prostate", t.cmd.Process.Pid)
	if t.cmd.ProcessState != nil && !t.cmd.ProcessState.Exited() {

	}
}

func (t *TaskMan) run(cf *changedFile) {
	t.runLock.Lock()
	defer t.runLock.Unlock()
	for i := 0; i < len(cfg.Command.Exec); i++ {
		carr := cmdParse2Array(cfg.Command.Exec[i], cf)
		log.Println("EXEC", carr)
		t.cmd = exec.Command(carr[0], carr[1:]...)
		//cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: syscall.CREATE_UNICODE_ENVIRONMENT}
		t.cmd.Stdin = os.Stdin
		t.cmd.Stdout = os.Stdout
		t.cmd.Stderr = os.Stderr
		t.cmd.Dir = projectFolder
		t.cmd.Env = os.Environ()
		err := t.cmd.Start()
		if err != nil {
			log.Println(PreError, "run command", carr, "error. ", err)
			break
		}
		err = t.cmd.Wait()
		if err != nil {
			log.Println(PreWarn, "command exec failed:", carr, err)
			break
		}
		if t.cmd.Process != nil {
			err := t.cmd.Process.Kill()
			log.Println(t.cmd.ProcessState)
			if t.cmd.ProcessState != nil && !t.cmd.ProcessState.Exited() {
				log.Println(PreError, "command cannot stop!", carr, err)
			}
		}
	}
	if keyInInstruction(InstShouldFinish) {
		t.cmd = nil
		t.waitChan <- true
	}
	log.Println("EXEC end")
}
