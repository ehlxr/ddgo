package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	dtn "github.com/JetBlink/dingtalk-notify-go-sdk"
	"github.com/ehlxr/ddgo/pkg"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"
	log "unknwon.dev/clog/v2"
)

var (
	dingTalk *dtn.Robot
	limiter  *pkg.LimiterServer
)

func init() {
	initLog()
}

func main() {
	pkg.ParseArg()

	dingTalk = dtn.NewRobot(pkg.Opts.Robot.Token, pkg.Opts.Robot.Secret)
	limiter = pkg.NewLimiterServer(1*time.Minute, 20)

	start()
}

func start() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Header.Get("Content-Type") {
		case "application/json":
			handlePostJson(w, r)
		case "application/x-www-form-urlencoded":
			handlePostForm(w, r)
		}
	})

	server := &http.Server{
		Addr:    pkg.Opts.Addr,
		Handler: mux,
	}

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	go func() {
		<-quit

		if err := server.Shutdown(context.Background()); err != nil {
			log.Fatal("Shutdown server:", err)
		}
	}()

	log.Info("Starting HTTP server on http://%s", pkg.Opts.Addr)
	if err := server.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			log.Info("Server closed under request")
		} else {
			log.Fatal("Server closed unexpected")
		}
	}
}

func initLog() {
	err := log.NewConsole()
	if err != nil {
		panic("unable to create new logger: " + err.Error())
	}

	err = log.NewFile(log.FileConfig{
		Level:    log.LevelInfo,
		Filename: "./logs/ddgo.log",
		FileRotationConfig: log.FileRotationConfig{
			Rotate: true,
			Daily:  true,
		},
	})
	if err != nil {
		panic("unable to create new logger: " + err.Error())
	}
}

func handlePostForm(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Error("parse request form %+v",
			err)
		_, _ = io.WriteString(w,
			fmt.Sprintf("parse request form %+v", err))
		return
	}

	if err = sendMessage(&Message{
		r.Form.Get("content"),
		r.Form.Get("at"),
		r.Form.Get("app"),
		r.Form.Get("dt"),
		r.Form.Get("ds"),
	}); err != nil {
		http.Error(w, err.Error(), 400)
		log.Error(err.Error())

		return
	}

	_, _ = io.WriteString(w, "success")
	return
}

func handlePostJson(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		http.Error(w, "Please sendMessage a request body", 400)
		log.Error("Please sendMessage a request body")
		return
	}

	m := &Message{}
	err := json.NewDecoder(r.Body).Decode(m)
	if err != nil {
		http.Error(w, err.Error(), 400)
		log.Error(err.Error())
		return
	}

	if err = sendMessage(m); err != nil {
		http.Error(w, err.Error(), 400)
		log.Error(err.Error())
		return
	}

	_, _ = io.WriteString(w, "success")
	return
}

type Message struct {
	Content string `json:"content" desc:"消息内容"`
	At      string `json:"at" desc:"被@人的手机号"`
	App     string `json:"app" desc:"发送消息应用名称（添加到消息之前）"`
	Dt      string `json:"dt" desc:"钉钉机器人 token"`
	Ds      string `json:"ds" desc:"钉钉机器人签名 secret"`
}

func sendMessage(m *Message) error {
	if m.Content == "" {
		return errors.New("content must not be null")
	}

	ats := pkg.Opts.Robot.AtMobiles
	if m.At != "" {
		ats = append(ats, strings.Split(m.At, ",")...)
	}

	if m.App != "" {
		m.Content = fmt.Sprintf("%s\n%s", m.App, m.Content)
	}

	dtRobot := dingTalk

	if m.Ds != "" && m.Dt != "" {
		dtRobot = dtn.NewRobot(m.Dt, m.Ds)
	}

	if limiter.IsAvailable() {
		err := dtRobot.SendTextMessage(m.Content, ats, pkg.Opts.Robot.IsAtAll)
		if err != nil {
			return err
		}

		log.Info("sendMessage message <%s> success", m.Content)
	} else {
		return errors.New("dingTalk 1 m allow sendMessage 20 msg. msg discarded")
	}
	return nil
}
