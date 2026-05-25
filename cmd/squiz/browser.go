package main

import (
	"os"
	"os/exec"
	"runtime"
)

// OpenInBrowser launches the OS default browser pointed at the given local
// file path. Honors SQUIZ_NO_OPEN — set to any non-empty value to make
// this a no-op. Used by tests that exercise the --open path: without the
// opt-out, the OS browser launches asynchronously and may try to read the
// file AFTER the test's t.TempDir cleanup has deleted it, producing a
// "file not found" popup on the user's desktop.
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
