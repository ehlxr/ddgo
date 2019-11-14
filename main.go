package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jessevdk/go-flags"
	"io"
	"io/ioutil"
	"net/http"
	log "unknwon.dev/clog/v2"
)

func init() {
	err := log.NewConsole()
	if err != nil {
		panic("unable to create new logger: " + err.Error())
	}

	err = log.NewFile(log.FileConfig{
		Level:    log.LevelInfo,
		Filename: "ddgo.log",
		FileRotationConfig: log.FileRotationConfig{
			Rotate: true,
			Daily:  true,
		},
	})
	if err != nil {
		panic("unable to create new logger: " + err.Error())
	}
}

var opts struct {
	Addr       string `short:"a" long:"addr" default:"0.0.0.0:10141" env:"ADDR" description:"Addr to listen on for http requests"`
	WebHookUrl string `short:"u" long:"webhook-url" env:"URL" description:"Webhook url of dingding" required:"true"`
}

func dingToInfo(msg string) []byte {
	content, data := make(map[string]string), make(map[string]interface{})

	content["content"] = msg
	data["msgtype"] = "text"
	data["text"] = content
	b, _ := json.Marshal(data)

	log.Info("send to %s data <%s>",
		opts.WebHookUrl,
		b)

	resp, err := http.Post(opts.WebHookUrl,
		"application/json",
		bytes.NewBuffer(b))
	if err != nil {
		log.Error("send request to %s %+v",
			opts.Addr,
			err)

	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	log.Info("send to %s data <%s> result is %s",
		opts.WebHookUrl,
		b,
		body)
	return body
}

func send(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Error("parse request form %+v",
			err)
		_, _ = io.WriteString(w,
			fmt.Sprintf("parse request form %+v", err))
		return
	}

	content := r.Form.Get("content")
	if content == "" {
		log.Error("read content from request form nil")
		_, _ = io.WriteString(w, "read content from request form nil")
		return
	}
	info := dingToInfo(content)

	_, _ = w.Write(info)
}

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal("parse arg %+v",
			err)
	}

	http.HandleFunc("/", send)

	log.Info("server on http://%s", opts.Addr)
	if err := http.ListenAndServe(opts.Addr, nil); err != nil {
		log.Fatal("ListenAndServe %+v",
			err)
	}
}
