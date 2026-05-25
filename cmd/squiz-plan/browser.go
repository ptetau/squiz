package main

import (
	"os"
	"os/exec"
	"runtime"
)

// OpenInBrowser launches the OS default browser pointed at the given local
// file path. Honors SQUIZ_NO_OPEN — set to any non-empty value to make
// this a no-op (used by tests; see cmd/squiz/browser.go for full rationale).
// Duplicated from cmd/squiz/browser.go on purpose — minimising cross-
// binary coupling is cheaper than a shared helper at this size.
func OpenInBrowser(path string) error {
	if os.Getenv("SQUIZ_NO_OPEN") != "" {
		return nil
	}
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
