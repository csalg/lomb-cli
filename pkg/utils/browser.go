package utils

import (
	"fmt"
	"os/exec"
	"runtime"
)

// OpenURLInBrowser opens a URL in the user's default browser.
// Ref: https://stackoverflow.com/questions/39320371/how-start-web-server-to-open-page-in-browser-in-golang
func OpenURLInBrowser(url_ string) {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url_)
	if err := exec.Command(cmd, args...).Start(); err != nil {
		fmt.Println("opening browser: " + err.Error())
	}
}
