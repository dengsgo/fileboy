package main

import (
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

func keyInMonitorTypesMap(k string, maps map[string]bool) bool {
	_, ok := maps[k]
	return ok
}

func keyInInstruction(k string) bool {
	_, ok := cfg.InstructionMap[k]
	return ok
}

func cmdParse2Array(s string, cf *changedFile) []string {
	a := strings.Split(s, " ")
	r := make([]string, 0)
	for i := 0; i < len(a); i++ {
		if ss := strings.Trim(a[i], " "); ss != "" {
			r = append(r, strParseRealStr(ss, cf))
		}
	}
	return r
}

func strParseRealStr(s string, cf *changedFile) string {
	return strings.ReplaceAll(
		strings.ReplaceAll(
			strings.ReplaceAll(
				strings.ReplaceAll(s, "{{file}}", cf.Name),
				"{{ext}}", cf.Ext,
			),
			"{{changed}}", strconv.FormatInt(cf.Changed, 10),
		),
		"{{event}}", cf.Event,
	)
}

func dirParse2Array(s string) []string {
	a := strings.Split(s, ",")
	r := make([]string, 0)
	for i := 0; i < len(a); i++ {
		if ss := strings.Trim(a[i], " "); ss != "" {
			r = append(r, ss)
		}
	}
	return r
}

func hitDirs(d string, dirs *[]string) bool {
	d += "/"
	for _, v := range *dirs {
		if strings.HasPrefix(d, projectFolder+"/"+v+"/") {
			return true
		}
	}
	return false
}

func listFile(folder string, fun func(string)) {
	files, _ := ioutil.ReadDir(folder)
	for _, file := range files {
		if file.IsDir() {
			d := folder + "/" + file.Name()
			if hitDirs(d, &cfg.Monitor.ExceptDirs) {
				continue
			}
			fun(d)
			listFile(d, fun)
		}
	}
}

func relativePath(folder, p string) string {
	s := strings.ReplaceAll(strings.TrimPrefix(p, folder), "\\", "/")
	if strings.HasPrefix(s, "/") && len(s) > 1 {
		s = s[1:]
	}
	return s
}

func inStrArray(s string, arr []string) bool {
	for _, v := range arr {
		if s == v {
			return true
		}
	}
	return false
}

func logInfo(v ...interface{}) {
	if keyInInstruction(InstIgnoreInfo) {
		return
	}
	logUInfo(v...)
}

func logUInfo(v ...interface{}) {
	v = append([]interface{}{"I:"}, v...)
	log.Println(v...)
}

func logWarn(v ...interface{}) {
	if keyInInstruction(InstIgnoreWarn) {
		return
	}
	v = append([]interface{}{"W:"}, v...)
	log.Println(v...)
}

func logError(v ...interface{}) {
	v = append([]interface{}{"E:"}, v...)
	log.Println(v...)
}

func logAndExit(v ...interface{}) {
	v = append([]interface{}{"O:"}, v...)
	log.Println(v...)
	os.Exit(15)
}
