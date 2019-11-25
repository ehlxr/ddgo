package pkg

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"os"
)

var (
	Opts struct {
		Addr    string `short:"a" long:"addr" default:":80" env:"ADDR" description:"Addr to listen on for HTTP server"`
		Version bool   `short:"v" long:"version" description:"Show version info"`
		Robot   Robot  `group:"DingTalk Robot Options" namespace:"robot" env-namespace:"ROBOT" `
	}
)

type Robot struct {
	Token     string   `short:"t" long:"token" env:"TOKEN" description:"DingTalk robot access token" required:"true"`
	Secret    string   `short:"s" long:"secret" env:"SECRET" description:"DingTalk robot secret"`
	AtMobiles []string `short:"m" long:"at-mobiles" env:"AT_MOBILES" env-delim:"," description:"The mobile of the person will be at"`
	IsAtAll   bool     `short:"e" long:"at-all" env:"AT_ALL" description:"Whether at everyone"`
}

func ParseArg() {
	parser := flags.NewParser(&Opts, flags.HelpFlag|flags.PassDoubleDash)
	parser.NamespaceDelimiter = "-"

	if AppName != "" {
		parser.Name = AppName
	}

	if _, err := parser.Parse(); err != nil {
		if Opts.Version {
			PrintVersion()
			os.Exit(0)
		}

		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			_, _ = fmt.Fprintln(os.Stdout, err)
			os.Exit(0)
		}

		_, _ = fmt.Fprintln(os.Stderr, err)

		parser.WriteHelp(os.Stderr)

		os.Exit(1)
	}
}
