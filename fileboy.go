package main

import (
	"bufio"
	"fmt"
	"gopkg.in/fsnotify/fsnotify.v1"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"sync"
	"time"
)

const (
	Version = 1
)

var (
	projectFolder string = "."

	cfg *FileGirl

	watcher *fsnotify.Watcher

	cmd *exec.Cmd

	runLock sync.Mutex
)

type wDirState struct {
	Name      string
	Recursive bool
}

func parseConfig() {
	cfg = new(FileGirl)
	fc, err := ioutil.ReadFile(projectFolder + "/filegirl.yaml")
	if err != nil {
		log.Panicln("read filegirl.yaml file err: ", err)
	}
	err = yaml.Unmarshal(fc, cfg)
	if err != nil {
		log.Panicln("parsed filegirl.yaml failed: ", err)
	}
	if cfg.Core.Version > Version {
		log.Panicln("current fileboy support max version : ", Version)
	}
	log.Println(cfg)
}

func eventDispatcher(event fsnotify.Event) {
	ext := path.Ext(event.Name)
	if len(cfg.Monitor.Types) > 0 &&
		cfg.Monitor.Types[0] != ".*" &&
		!inStringArray(ext, cfg.Monitor.Types) {
		//log.Println(ext, cfg.Monitor.Types, inStringArray(ext, cfg.Monitor.Types))
		return
	}
	switch event.Op {
	case fsnotify.Create:
	case fsnotify.Write:
		log.Println("event write : ", event.Name)
		if cmd != nil {
			err := cmd.Process.Kill()
			if err != nil {
				log.Println("err: ", err)
			}
			log.Println("stop old process ")
		}
		go run()
	case fsnotify.Remove:
	case fsnotify.Rename:
	}
}

func run() {
	runLock.Lock()
	defer runLock.Unlock()
	for i := 0; i < len(cfg.Command.Exec); i++ {
		carr := cmdParse2Array(cfg.Command.Exec[i])
		cmd = exec.Command(carr[0], carr[1:]...)
		//cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: syscall.CREATE_UNICODE_ENVIRONMENT}
		cmd.Stdin = os.Stdin
		//cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			log.Println("error=>", err.Error())
			return
		}
		cmd.Start()
		reader := bufio.NewReader(stdout)
		for {
			line, err2 := reader.ReadString('\n')
			if err2 != nil || io.EOF == err2 {
				break
			}
			log.Print(line)
		}
		err = cmd.Wait()
		if err != nil {
			log.Println("cmd wait err ", err)
			break
		}
		err = cmd.Process.Kill()
		if err != nil {
			log.Println("cmd cannot kill ", err)
		}
	}

	log.Println("end ")
}

func addWatcher() {
	dirs := make([]string, 0)
	for i := 0; i < len(cfg.Monitor.IncludeDirs); i++ {
		darr := dirParse2Array(cfg.Monitor.IncludeDirs[i])
		if len(darr) < 1 || len(darr) > 2 {
			log.Fatalln("filegirl section monitor dirs is error. ", cfg.Monitor.IncludeDirs[i])
		}
		if darr[0] == "." {
			if len(darr) == 2 && darr[1] == "*" {
				dirs = make([]string, 0)
				dirs = append(dirs, ".")
				listFile(projectFolder, func(d string) {
					dirs = arrayUniqueAdd(dirs, d)
				})
			} else {
				dirs = arrayUniqueAdd(dirs, projectFolder)
			}
			break
		} else {
			if len(darr) == 2 && darr[1] == "*" {
				listFile(projectFolder+"/"+darr[0], func(d string) {
					dirs = arrayUniqueAdd(dirs, d)
				})
			} else {
				dirs = arrayUniqueAdd(dirs, projectFolder+"/"+darr[0])
			}
		}

	}
	for i := 0; i < len(cfg.Monitor.ExceptDirs); i++ {
		p := projectFolder + "/" + cfg.Monitor.ExceptDirs[i]
		dirs = arrayRemoveElement(dirs, p)
		listFile(p, func(d string) {
			dirs = arrayRemoveElement(dirs, d)
		})
	}
	for _, dir := range dirs {
		log.Println("watcher add -> ", dir)
		err := watcher.Add(dir)
		if err != nil {
			log.Fatalln(err)
		}
	}
	log.Println("fileboy is ready.")
}

func initWatcher() {
	parseConfig()
	var err error
	watcher, err = fsnotify.NewWatcher()
	if err != nil {
		log.Fatalln(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				eventDispatcher(event)
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()
	addWatcher()
	<-done
}

func parseArgs() {
	l := len(os.Args)
	if l == 1 {
		_, err := ioutil.ReadFile(projectFolder + "/filegirl.yaml")
		if err != nil {
			log.Println("the filegirl.yaml file is not exist! ", err)
			fmt.Print(firstRunHelp)
			fmt.Print(helpStr)
			return
		}
		initWatcher()
		return
	}
	if l == 2 {
		c := os.Args[1]
		switch c {
		case "init":
			err := ioutil.WriteFile(projectFolder+"/filegirl.yaml", []byte(exampleFileGirl), 0644)
			if err != nil {
				log.Println("error create filegirl.yaml config! ", err)
				return
			}
			log.Println("create filegirl.yaml ok")
			return
		case "exec":
			parseConfig()
			run()
			return
		default:
			fmt.Print(helpStr)
		}
		return
	}
}

func show() {
	fmt.Print(logo)
	rand.Seed(time.Now().UnixNano())
	fmt.Println(englishSay[rand.Intn(len(englishSay))], "\r\n")
	fmt.Println("Version: ", Version, "   Author: deng@yoytang.com")
}

func main() {
	log.SetPrefix("[FileBoy]: ")
	log.SetFlags(2)
	show()
	var err error
	projectFolder, err = os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}
	parseArgs()
}
