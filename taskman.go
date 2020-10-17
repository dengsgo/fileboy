package main

import (
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
				if len(t.waitQueue) < 1 {
					return
				}
				cf := t.waitQueue[len(t.waitQueue)-1]
				if len(t.waitQueue) > 1 {
					logInfo("redundant tasks dropped:", len(t.waitQueue)-1)
				}
				t.waitQueue = []*changedFile{}
				go t.preRun(cf)
			}
		}()
	}

	return t
}

func (t *TaskMan) Put(cf *changedFile) {
	if t.delay < 1 {
		t.dispatcher(cf)
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
		t.dispatcher(cf)
	}()
}

func (t *TaskMan) dispatcher(cf *changedFile) {
	if keyInInstruction(InstShouldFinish) {
		t.waitQueue = append(t.waitQueue, cf)
		if t.cmd == nil {
			t.waitChan <- true
			return
		}
		logInfo("waitting for the last task to finish")
		logInfo("waiting tasks:", len(t.waitQueue))
	} else {
		t.preRun(cf)
	}
}

func (t *TaskMan) preRun(cf *changedFile) {
	if t.cmd != nil && t.cmd.Process != nil {
		if err := t.cmd.Process.Kill(); err != nil {
			logInfo("stop old process ")
			logWarn("stopped err, reason:", err)
		}
	}
	go t.run(cf)
	go t.notifier.Put(cf)
}

func (t *TaskMan) run(cf *changedFile) {
	t.runLock.Lock()
	defer t.runLock.Unlock()
	for i := 0; i < len(cfg.Command.Exec); i++ {
		carr := cmdParse2Array(cfg.Command.Exec[i], cf)
		logInfo("EXEC", carr)
		t.cmd = exec.Command(carr[0], carr[1:]...)
		//cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: syscall.CREATE_UNICODE_ENVIRONMENT}
		t.cmd.Stdin = os.Stdin
		t.cmd.Stdout = os.Stdout
		if keyInInstruction(InstIgnoreStdout) {
			t.cmd.Stdout = nil
		}
		t.cmd.Stderr = os.Stderr
		t.cmd.Dir = projectFolder
		t.cmd.Env = os.Environ()
		err := t.cmd.Start()
		if err != nil {
			logError("run command", carr, "error. ", err)
			if keyInInstruction(InstIgnoreExecError) {
				continue
			}
			break
		}
		err = t.cmd.Wait()
		if err != nil {
			logError("command exec failed:", carr, err)
			if keyInInstruction(InstIgnoreExecError) {
				continue
			}
			break
		}
		if t.cmd.Process != nil {
			err := t.cmd.Process.Kill()
			logInfo(t.cmd.ProcessState)
			if t.cmd.ProcessState != nil && !t.cmd.ProcessState.Exited() {
				logError("command cannot stop!", carr, err)
			}
		}
	}
	if keyInInstruction(InstShouldFinish) {
		t.cmd = nil
		t.waitChan <- true
	}
	logInfo("EXEC end")
}
