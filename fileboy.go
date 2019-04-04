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
	"strconv"
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
		log.Println(PreError, "the filegirl.yaml file in", projectFolder, "is not exist! ", err)
		fmt.Print(firstRunHelp)
		log.Fatalln("fileboy unable to run.")
	}
	err = yaml.Unmarshal(fc, cfg)
	if err != nil {
		log.Panicln(PreError, "parsed filegirl.yaml failed: ", err)
	}
	if cfg.Core.Version > Version {
		log.Panicln(PreError, "current fileboy support max version : ", Version)
	}
	// init map
	cfg.Monitor.TypesMap = map[string]bool{}
	cfg.Monitor.IncludeDirsMap = map[string]bool{}
	cfg.Monitor.ExceptDirsMap = map[string]bool{}
	// convert to map
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
	dirsMap := map[string]bool{
		projectFolder: true,
	}
	for _, dir := range cfg.Monitor.IncludeDirs {
		darr := dirParse2Array(dir)
		if len(darr) < 1 || len(darr) > 2 {
			log.Fatalln(PreError, "filegirl section monitor dirs is error. ", dir)
		}
		if strings.HasPrefix(darr[0], "/") {
			log.Fatalln(PreError, "dirs must be relative paths ! err path:", dir)
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
			}
		}

	}
	for _, dir := range cfg.Monitor.ExceptDirs {
		if dir == "." {
			log.Fatalln(PreError, "exceptDirs must is not project root path ! err path:", dir)
		}
		p := projectFolder + "/" + dir
		delete(dirsMap, p)
		listFile(p, func(d string) {
			delete(dirsMap, d)
		})
	}
	for dir := range dirsMap {
		log.Println("watcher add -> ", dir)
		err := watcher.Add(dir)
		if err != nil {
			log.Fatalln(err)
		}
	}
	log.Println("total monitored dirs: " + strconv.Itoa(len(dirsMap)))
	log.Println("fileboy is ready.")
	cfg.Monitor.DirsMap = dirsMap
}

func initWatcher() {
	var err error
	if watcher != nil {
		_ = watcher.Close()
	}
	watcher, err = fsnotify.NewWatcher()
	if err != nil {
		log.Fatalln(err)
	}
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
}

func parseArgs() {
	switch len(os.Args) {
	case 1:
		parseConfig()
		done := make(chan bool)
		initWatcher()
		defer watcher.Close()
		<-done
		return
	case 2:
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
	default:
		log.Fatalln("Unknown parameters, use `fileboy help` show help info.")
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
