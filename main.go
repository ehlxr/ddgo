package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/jessevdk/go-flags"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	log "unknwon.dev/clog/v2"
)

var (
	AppName   string
	Version   string
	BuildTime string
	GitCommit string
	GoVersion string

	versionTpl = `%s
Name: %s
Version: %s
BuildTime: %s
GitCommit: %s
GoVersion: %s

`
	bannerBase64 = "DQogX19fXyAgX19fXyAgICBfX18gIF9fX19fIA0KKCAgXyBcKCAgXyBcICAvIF9fKSggIF8gICkNCiApKF8pICkpKF8pICkoIChfLS4gKShfKSggDQooX19fXy8oX19fXy8gIFxfX18vKF9fX19fKQ0K"

	opts struct {
		Addr       string `short:"a" long:"addr" default:":80" env:"ADDR" description:"Addr to listen on for HTTP server"`
		WebHookUrl string `short:"u" long:"webhook-url" env:"URL" description:"Webhook url of dingding" required:"true"`
		Version    bool   `short:"v" long:"version" description:"Show version info"`
	}
)

func init() {
	initLog()
}

func main() {
	parseArg()

	mux := http.NewServeMux()
	mux.HandleFunc("/", requestHandle)

	server := &http.Server{
		Addr:    opts.Addr,
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

	log.Info("Starting HTTP server on http://%s", opts.Addr)
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

func parseArg() {
	parser := flags.NewParser(&opts, flags.HelpFlag|flags.PassDoubleDash)
	if _, err := parser.Parse(); err != nil {
		if opts.Version {
			printVersion()
			os.Exit(0)
		}

		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			_, _ = fmt.Fprintln(os.Stdout, err)
			os.Exit(0)
		}

		_, _ = fmt.Fprintln(os.Stderr, err)

		parser.Name = AppName
		parser.WriteHelp(os.Stderr)

		os.Exit(1)
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
	info := dingToInfo(content)

	_, _ = w.Write(info)
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

// printVersion Print out version information
func printVersion() {
	banner, _ := base64.StdEncoding.DecodeString(bannerBase64)
	fmt.Printf(versionTpl, banner, AppName, Version, BuildTime, GitCommit, GoVersion)
}
