package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
)

type postParams struct {
	ProjectFolder string `json:"project_folder"`
	File          string `json:"file"`
	Changed       int64  `json:"changed"`
	Ext           string `json:"ext"`
}

type NetNotifier struct {
	CallUrl string
	CanPost bool
}

func newNetNotifier(callUrl string) *NetNotifier {
	callPost := true
	if strings.TrimSpace(callUrl) == "" {
		callPost = false
	}
	return &NetNotifier{
		CallUrl: callUrl,
		CanPost: callPost,
	}
}

func (n *NetNotifier) Put(cf *changedFile) {
	if !n.CanPost {
		log.Println(PreWarn, "notifier call url ignore. ", n.CallUrl)
		return
	}
	n.dispatch(&postParams{
		ProjectFolder: projectFolder,
		File:          cf.Name,
		Changed:       cf.Changed,
		Ext:           cf.Ext,
	})
}

func (n *NetNotifier) dispatch(params *postParams) {
	b, err := json.Marshal(params)
	if err != nil {
		log.Println(PreError, "json.Marshal n.params. ", err)
		return
	}
	client := &http.Client{
		Timeout: time.Second * 15,
	}
	req, err := http.NewRequest("POST", n.CallUrl, bytes.NewBuffer(b))
	if err != nil {
		log.Println(PreError, "http.NewRequest. ", err)
		return
	}
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("User-Agent", "FileBoy Net Notifier v1.8")
	resp, err := client.Do(req)
	if err != nil {
		log.Println(PreError, "notifier call failed. err:", err)
		return
	}
	defer func() {
		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}
	}()
	if resp.StatusCode >= 300 {
		// todo retry???
	}
	log.Println("notifier done .")
}
