package main

import (
	"os/exec"
	"runtime"
)

// OpenInBrowser launches the OS default browser pointed at the given local file path.
// Duplicated from cmd/squiz/browser.go on purpose — the brief calls for
// minimising cross-binary coupling; 12 lines is cheaper than a shared helper.
func OpenInBrowser(path string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		// cmd /c start "" "path" — the empty "" is the window title; needed because
		// start treats the first quoted arg as the title.
		cmd = exec.Command("cmd", "/c", "start", "", path)
	case "darwin":
		cmd = exec.Command("open", path)
	default:
		cmd = exec.Command("xdg-open", path)
	}
	return cmd.Start()
}
