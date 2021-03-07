// Copyright (c) 2018-2021 Author dengsgo<dengsgo@yoytang.com> [https://github.com/dengsgo/fileboy]
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

type postParams struct {
	ProjectFolder string `json:"project_folder"`
	File          string `json:"file"`
	Changed       int64  `json:"changed"`
	Ext           string `json:"ext"`
	Event         string `json:"event"`
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
		logWarn("notifier call url ignore. ", n.CallUrl)
		return
	}
	n.dispatch(&postParams{
		ProjectFolder: projectFolder,
		File:          cf.Name,
		Changed:       cf.Changed,
		Ext:           cf.Ext,
		Event:         cf.Event,
	})
}

func (n *NetNotifier) dispatch(params *postParams) {
	b, err := json.Marshal(params)
	if err != nil {
		logError("json.Marshal n.params. ", err)
		return
	}
	client := http.DefaultClient
	client.Timeout = time.Second * 15
	req, err := http.NewRequest("POST", n.CallUrl, bytes.NewBuffer(b))
	if err != nil {
		logError("http.NewRequest. ", err)
		return
	}
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("User-Agent", "FileBoy Net Notifier v1.16")
	resp, err := client.Do(req)
	if err != nil {
		logError("notifier call failed. err:", err)
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
	logInfo("notifier done .")
}
