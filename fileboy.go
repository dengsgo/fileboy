package main

import (
	"fmt"
	"gopkg.in/fsnotify/fsnotify.v1"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path"
	"strings"
	"time"
)

const (
	Version = 1

	PreError = "ERROR:"
	PreWarn  = "Warn:"
)

var (
	projectFolder = "."

	cfg *FileGirl

	watcher *fsnotify.Watcher

	taskMan *TaskMan
)

type changedFile struct {
	Name    string
	Changed int64
	Ext     string
}

func parseConfig() {
	cfg = new(FileGirl)
	fc, err := ioutil.ReadFile(projectFolder + "/filegirl.yaml")
	if err != nil {
		log.Panicln(PreError, "read filegirl.yaml file err: ", err)
	}
	err = yaml.Unmarshal(fc, cfg)
	if err != nil {
		log.Panicln(PreError, "parsed filegirl.yaml failed: ", err)
	}
	if cfg.Core.Version > Version {
		log.Panicln(PreError, "current fileboy support max version : ", Version)
	}
	// types convert map
	cfg.Monitor.TypesMap = map[string]bool{}
	for _, v := range cfg.Monitor.Types {
		cfg.Monitor.TypesMap[v] = true
	}
	log.Println(cfg)
}

func eventDispatcher(event fsnotify.Event) {
	ext := path.Ext(event.Name)
	if len(cfg.Monitor.Types) > 0 &&
		!keyInMonitorTypesMap(".*", cfg) &&
		!keyInMonitorTypesMap(ext, cfg) {
		return
	}
	switch event.Op {
	case
		fsnotify.Write,
		fsnotify.Rename:
		log.Println("EVENT", event.Op.String(), ":", event.Name)
		taskMan.Put(&changedFile{
			Name:    relativePath(projectFolder, event.Name),
			Changed: time.Now().UnixNano(),
			Ext:     ext,
		})
	case fsnotify.Remove:
	case fsnotify.Create:
	}
}

func addWatcher() {
	log.Println("collecting directory information...")
	dirs := make([]string, 0)
	for i := 0; i < len(cfg.Monitor.IncludeDirs); i++ {
		darr := dirParse2Array(cfg.Monitor.IncludeDirs[i])
		if len(darr) < 1 || len(darr) > 2 {
			log.Fatalln(PreError, "filegirl section monitor dirs is error. ", cfg.Monitor.IncludeDirs[i])
		}
		if strings.HasPrefix(darr[0], "/") {
			log.Fatalln(PreError, "dirs must be relative paths ! err path:", cfg.Monitor.IncludeDirs[i])
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
			md := projectFolder + "/" + darr[0]
			if len(darr) == 2 && darr[1] == "*" {
				dirs = arrayUniqueAdd(dirs, md)
				listFile(md, func(d string) {
					dirs = arrayUniqueAdd(dirs, d)
				})
			} else {
				dirs = arrayUniqueAdd(dirs, md)
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
	taskMan = newTaskMan(cfg.Command.DelayMillSecond, cfg.Notifier.CallUrl)
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
				log.Println(PreError, err)
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
			log.Println(PreError, "the filegirl.yaml file does not exist! ", err)
			fmt.Print(firstRunHelp)
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
				log.Println(PreError, "error create filegirl.yaml config! ", err)
				return
			}
			log.Println("create filegirl.yaml ok")
			return
		case "exec":
			parseConfig()
			newTaskMan(0, cfg.Notifier.CallUrl).run(new(changedFile))
			return
		case "version", "v", "-v", "--version":
			fmt.Println(versionDesc)
		default:
			fmt.Print(helpStr)
		}
		return
	}
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
	show()
	var err error
	projectFolder, err = os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}
	parseArgs()
}
