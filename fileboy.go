// Copyright (c) 2018-2021 Author dengsgo<dengsgo@yoytang.com> [https://github.com/dengsgo/fileboy]
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"path"
	"strconv"
	"strings"
	"syscall"
	"time"

	"gopkg.in/fsnotify/fsnotify.v1"
	"gopkg.in/yaml.v2"
)

const (
	Version = 1

	InstExecWhenStart   = "exec-when-start"
	InstShouldFinish    = "should-finish"
	InstIgnoreWarn      = "ignore-warn"
	InstIgnoreInfo      = "ignore-info"
	InstIgnoreStdout    = "ignore-stdout"
	InstIgnoreExecError = "ignore-exec-error"
)

var (
	projectFolder = "."

	filegirlYamlName = "filegirl.yaml"

	cfg *FileGirl

	watcher *fsnotify.Watcher

	taskMan *TaskMan

	ioeventMapStr = map[fsnotify.Op]string{
		fsnotify.Write:  "write",
		fsnotify.Rename: "rename",
		fsnotify.Remove: "remove",
		fsnotify.Create: "create",
		fsnotify.Chmod:  "chmod",
	}
)

type changedFile struct {
	Name    string
	Changed int64
	Ext     string
	Event   string
}

func parseConfig() {
	cfg = new(FileGirl)
	fc, err := ioutil.ReadFile(getFileGirlPath())
	if err != nil {
		logError("the filegirl.yaml file in", projectFolder, "is not exist! ", err)
		fmt.Print(firstRunHelp)
		logAndExit("fileboy unable to run.")
	}
	err = yaml.Unmarshal(fc, cfg)
	if err != nil {
		logAndExit("parsed filegirl.yaml failed: ", err)
	}
	if cfg.Core.Version > Version {
		logAndExit("current fileboy support max version : ", Version)
	}
	// init map
	cfg.Monitor.TypesMap = map[string]bool{}
	cfg.Monitor.IncludeDirsMap = map[string]bool{}
	cfg.Monitor.ExceptDirsMap = map[string]bool{}
	cfg.Monitor.IncludeDirsRec = map[string]bool{}
	cfg.InstructionMap = map[string]bool{}
	// convert to map
	for _, v := range cfg.Monitor.Types {
		cfg.Monitor.TypesMap[v] = true
	}
	for _, v := range cfg.Instruction {
		cfg.InstructionMap[v] = true
	}
	log.Printf("%+v", cfg)
}

func eventDispatcher(event fsnotify.Event) {
	if event.Name == getPidFile() {
		return
	}
	ext := path.Ext(event.Name)
	if len(cfg.Monitor.Types) > 0 &&
		!keyInMonitorTypesMap(".*", cfg) &&
		!keyInMonitorTypesMap(ext, cfg) {
		return
	}

	op := ioeventMapStr[event.Op]
	if len(cfg.Monitor.Events) != 0 && !inStrArray(op, cfg.Monitor.Events) {
		return
	}
	log.Println("EVENT", event.Op.String(), ":", event.Name)
	taskMan.Put(&changedFile{
		Name:    relativePath(projectFolder, event.Name),
		Changed: time.Now().UnixNano(),
		Ext:     ext,
		Event:   op,
	})
}

func addWatcher() {
	logInfo("collecting directory information...")
	dirsMap := map[string]bool{}
	for _, dir := range cfg.Monitor.ExceptDirs {
		if dir == "." {
			logAndExit("exceptDirs must is not project root path ! err path:", dir)
		}
	}
	for _, dir := range cfg.Monitor.IncludeDirs {
		darr := dirParse2Array(dir)
		if len(darr) < 1 || len(darr) > 2 {
			logAndExit("filegirl section monitor dirs is error. ", dir)
		}
		if strings.HasPrefix(darr[0], "/") {
			logAndExit("dirs must be relative paths ! err path:", dir)
		}
		if darr[0] == "." {
			if len(darr) == 2 && darr[1] == "*" {
				// The highest priority
				dirsMap = map[string]bool{
					projectFolder: true,
				}
				listFile(projectFolder, func(d string) {
					dirsMap[d] = true
				})
				cfg.Monitor.IncludeDirsRec[projectFolder] = true
				break
			} else {
				dirsMap[projectFolder] = true
			}
		} else {
			md := projectFolder + "/" + darr[0]
			dirsMap[md] = true
			if len(darr) == 2 && darr[1] == "*" {
				listFile(md, func(d string) {
					dirsMap[d] = true
				})
				cfg.Monitor.IncludeDirsRec[md] = true
			}
		}

	}

	for dir := range dirsMap {
		logInfo("watcher add -> ", dir)
		err := watcher.Add(dir)
		if err != nil {
			logAndExit(err)
		}
	}
	logInfo("total monitored dirs: " + strconv.Itoa(len(dirsMap)))
	logInfo("fileboy is ready.")
	cfg.Monitor.DirsMap = dirsMap
}

func initWatcher() {
	var err error
	if watcher != nil {
		_ = watcher.Close()
	}
	watcher, err = fsnotify.NewWatcher()
	if err != nil {
		logAndExit(err)
	}
	taskMan = newTaskMan(cfg.Command.DelayMillSecond, cfg.Notifier.CallUrl)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				// directory structure changes, dynamically add, delete and monitor according to rules
				// TODO // this method cannot be triggered when the parent folder of the change folder is not monitored
				go watchChangeHandler(event)
				eventDispatcher(event)
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				logError(err)
			}
		}
	}()
	addWatcher()
}

func watchChangeHandler(event fsnotify.Event) {
	// stop the fileboy daemon process when the .fileboy.pid file is changed
	if event.Name == getPidFile() &&
		(event.Op == fsnotify.Remove ||
			event.Op == fsnotify.Write ||
			event.Op == fsnotify.Rename) {
		logUInfo("exit daemon process")
		stopSelf()
		return
	}
	if event.Op != fsnotify.Create && event.Op != fsnotify.Rename {
		return
	}
	_, err := ioutil.ReadDir(event.Name)
	if err != nil {
		return
	}
	do := false
	for rec := range cfg.Monitor.IncludeDirsRec {
		if !strings.HasPrefix(event.Name, rec) {
			continue
		}
		// check exceptDirs
		if hitDirs(event.Name, &cfg.Monitor.ExceptDirs) {
			continue
		}

		_ = watcher.Remove(event.Name)
		err := watcher.Add(event.Name)
		if err == nil {
			do = true
			logInfo("watcher add -> ", event.Name)
		} else {
			logWarn("watcher add faild:", event.Name, err)
		}
	}

	if do {
		return
	}

	// check map
	if _, ok := cfg.Monitor.DirsMap[event.Name]; ok {
		_ = watcher.Remove(event.Name)
		err := watcher.Add(event.Name)
		if err == nil {
			logInfo("watcher add -> ", event.Name)
		} else {
			logWarn("watcher add faild:", event.Name, err)
		}
	}
}

func parseArgs() {
	switch {
	case len(os.Args) == 1:
		show()
		parseConfig()
		done := make(chan bool)
		initWatcher()
		defer watcher.Close()
		if keyInInstruction(InstExecWhenStart) {
			taskMan.run(new(changedFile))
		}
		<-done
		return
	case len(os.Args) > 1:
		c := os.Args[1]
		switch c {
		case "deamon", "daemon":
			pid, err := runAsDaemon()
			if err != nil {
				logAndExit(err)
			}
			logUInfo("PID:", pid)
			logUInfo("fileboy is ready. the main process will run as a daemons")
			return
		case "stop":
			err := stopDaemon()
			if err != nil {
				logAndExit(err)
			}
			logUInfo("fileboy daemon is stoped.")
			return
		case "init":
			_, err := ioutil.ReadFile(getFileGirlPath())
			if err == nil {
				logError("profile filegirl.yaml already exists.")
				logAndExit("delete it first when you want to regenerate filegirl.yaml")
			}
			err = ioutil.WriteFile(getFileGirlPath(), []byte(exampleFileGirl), 0644)
			if err != nil {
				logError("profile filegirl.yaml create failed! ", err)
				return
			}
			logUInfo("profile filegirl.yaml created ok")
			return
		case "exec":
			parseConfig()
			newTaskMan(0, cfg.Notifier.CallUrl).run(new(changedFile))
			return
		case "version", "v", "-v", "--version":
			fmt.Println(versionDesc)
		case "help", "--help", "--h", "-h":
			fmt.Print(helpStr)
		default:
			logAndExit("unknown parameter, use 'fileboy help' to view available commands")
		}
		return
	default:
		logAndExit("unknown parameters, use `fileboy help` show help info.")
	}
}

func signalHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		if taskMan != nil && taskMan.cmd != nil && taskMan.cmd.Process != nil {
			if err := taskMan.cmd.Process.Kill(); err != nil {
				logWarn("stopping the process failed: PID:", taskMan.cmd.ProcessState.Pid(), ":", err)
			}
		}
		os.Exit(0)
	}()
}

func getFileGirlPath() string {
	return projectFolder + "/" + filegirlYamlName
}

func show() {
	fmt.Print(logo)
	rand.Seed(time.Now().UnixNano())
	fmt.Println(englishSay[rand.Intn(len(englishSay))])
	fmt.Println("")
	fmt.Println(statement)
}

func main() {
	log.SetPrefix("[FileBoy]: ")
	log.SetFlags(2)
	log.SetOutput(os.Stdout)
	// show()
	var err error
	projectFolder, err = os.Getwd()
	if err != nil {
		logAndExit(err)
	}
	signalHandler()
	parseArgs()
}
