package main

import (
	"context"
	"fmt"
	dt "github.com/JetBlink/dingtalk-notify-go-sdk"
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
	dingTalk *dt.Robot
	limiter  *pkg.LimiterServer
)

func init() {
	initLog()
}

func main() {
	pkg.ParseArg()

	dingTalk = dt.NewRobot(pkg.Opts.Robot.Token, pkg.Opts.Robot.Secret)
	limiter = pkg.NewLimiterServer(1*time.Minute, 20)

	start()
}

func start() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", requestHandle)

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

func requestHandle(w http.ResponseWriter, r *http.Request) {
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

	ats := pkg.Opts.Robot.AtMobiles
	at := r.Form.Get("at")
	if at != "" {
		ats = append(ats, strings.Split(at, ",")...)
	}

	app := r.Form.Get("app")
	if app != "" {
		content = fmt.Sprintf("%s\n%s", app, content)
	}

	if limiter.IsAvailable() {
		err = dingTalk.SendTextMessage(content, ats, pkg.Opts.Robot.IsAtAll)
		if err != nil {
			log.Error("%+v", err)
			_, _ = fmt.Fprintln(w, err)
			return
		}

		log.Info("send message <%s> success", content)
		_, _ = io.WriteString(w, "success")
		return
	} else {
		log.Error("dingTalk 1 m allow send 20 msg. msg %v discarded.",
			content)

		_, _ = io.WriteString(w, "dingTalk 1 m allow send 20 msg. msg discarded.")
		return
	}
}
