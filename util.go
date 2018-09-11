package main

import (
	"io/ioutil"
	"strings"
)

func inStringArray(value string, arr []string) bool {
	for _, v := range arr {
		if value == v {
			return true
		}
	}
	return false
}

func cmdParse2Array(s string) []string {
	a := strings.Split(s, " ")
	r := make([]string, 0)
	for i := 0; i < len(a); i++ {
		if ss := strings.Trim(a[i], " "); ss != "" {
			r = append(r, ss)
		}
	}
	return r
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

func listFile(folder string, fun func(string)) {
	files, _ := ioutil.ReadDir(folder)
	for _, file := range files {
		if file.IsDir() {
			d := folder + "/" + file.Name()
			fun(d)
			listFile(d, fun)
		}
	}
}

func arrayUniqueAdd(a []string, add string) []string {
	if inStringArray(add, a) {
		return a
	}
	return append(a, add)
}

func arrayRemoveElement(a []string, r string) []string {
	i := -1
	for k, v := range a {
		if v == r {
			i = k
			break
		}
	}
	if i == -1 {
		return a
	}
	if len(a) == 1 && i == 0 {
		return []string{}
	}
	return append(a[:i], a[i+1:]...)
}
