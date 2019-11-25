package pkg

import (
	"encoding/base64"
	"fmt"
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
)

// PrintVersion Print out version information
func PrintVersion() {
	banner, _ := base64.StdEncoding.DecodeString(bannerBase64)
	fmt.Printf(versionTpl, banner, AppName, Version, BuildTime, GitCommit, GoVersion)
}
