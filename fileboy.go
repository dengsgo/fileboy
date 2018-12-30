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
	"time"
)

const (
	Version = 1
)

var (
	projectFolder string = "."

	cfg *FileGirl

	watcher *fsnotify.Watcher

	taskMan *TaskMan
)

type changeFile struct {
	Name    string
	Changed int64
	Ext     string
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
		taskMan.Put(&changeFile{
			Name:    relativePath(projectFolder, event.Name),
			Changed: time.Now().UnixNano(),
			Ext:     ext,
		})
	case fsnotify.Remove:
	case fsnotify.Rename:
	}
}

func addWatcher() {
	log.Println("collecting directory information...")
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
	taskMan = newTaskMan(2000)
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
			log.Println("the filegirl.yaml file does not exist! ", err)
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
			newTaskMan(0).run(new(changeFile))
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
