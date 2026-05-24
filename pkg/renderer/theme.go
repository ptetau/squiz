package renderer

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

// Themes ship with the binary. New themes appended to the end of this list
// will start showing up as repos are seen for the first time.
var Themes = []string{
	"paper",
	"phosphor",
	"amber",
	"beige",
	"rose",
	"ocean",
	"forest",
	"slate",
}

// themeCache is the on-disk persisted mapping `repo-key → theme`. The
// rotation index advances each time a new repo is encountered so themes
// cycle deterministically (paper → phosphor → … → slate → paper).
type themeCache struct {
	RotationIndex int               `json:"rotation_index"`
	Repos         map[string]string `json:"repos"`
}

var cacheMu sync.Mutex

func cachePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".squiz", "themes.json"), nil
}

func loadCache() *themeCache {
	c := &themeCache{Repos: map[string]string{}}
	path, err := cachePath()
	if err != nil {
		return c
	}
	data, err := os.ReadFile(path)
	if err != nil {
		// Missing or unreadable cache = fresh start. Not an error worth surfacing.
		return c
	}
	if err := json.Unmarshal(data, c); err != nil {
		fmt.Fprintf(os.Stderr, "warn: theme cache corrupt (%v), starting fresh\n", err)
		return &themeCache{Repos: map[string]string{}}
	}
	if c.Repos == nil {
		c.Repos = map[string]string{}
	}
	return c
}

func saveCache(c *themeCache) {
	path, err := cachePath()
	if err != nil {
		return
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		fmt.Fprintf(os.Stderr, "warn: theme cache mkdir: %v\n", err)
		return
	}
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return
	}
	// Atomic write: temp + rename so a crashed save doesn't leave a half-file.
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o600); err != nil {
		fmt.Fprintf(os.Stderr, "warn: theme cache write: %v\n", err)
		return
	}
	if err := os.Rename(tmp, path); err != nil {
		fmt.Fprintf(os.Stderr, "warn: theme cache rename: %v\n", err)
	}
}

// repoKey is the cache key for `workDir`. Prefers `git remote get-url origin`
// (survives cloning, consistent across machines); falls back to the canonical
// absolute path when not a git repo or git is unavailable.
func repoKey(workDir string) string {
	abs, err := filepath.Abs(workDir)
	if err != nil {
		abs = workDir
	}
	if remote := gitRemote(abs); remote != "" {
		return "git:" + remote
	}
	return "path:" + filepath.ToSlash(abs)
}

func gitRemote(dir string) string {
	cmd := exec.Command("git", "-C", dir, "remote", "get-url", "origin")
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

// ResolveTheme returns the theme to use for a render, applying precedence:
//  1. explicit override (CLI flag or JSON `theme` field) wins
//  2. cached mapping for this repo
//  3. next theme in rotation (assigned, persisted, and returned)
//
// Falls back to "paper" if anything goes wrong.
func ResolveTheme(workDir, override string) string {
	if override != "" {
		if !validTheme(override) {
			fmt.Fprintf(os.Stderr, "warn: unknown theme %q, falling back to auto\n", override)
		} else {
			return override
		}
	}

	cacheMu.Lock()
	defer cacheMu.Unlock()

	cache := loadCache()
	key := repoKey(workDir)

	if t, ok := cache.Repos[key]; ok && validTheme(t) {
		return t
	}

	// New repo: assign next theme in rotation and persist.
	idx := cache.RotationIndex % len(Themes)
	chosen := Themes[idx]
	cache.Repos[key] = chosen
	cache.RotationIndex++
	saveCache(cache)
	return chosen
}

func validTheme(name string) bool {
	for _, t := range Themes {
		if t == name {
			return true
		}
	}
	return false
}

// ThemeForRepoOnly returns the theme that WOULD be assigned to this workDir
// without persisting anything. Used by introspection commands.
func ThemeForRepoOnly(workDir string) (string, string, error) {
	cache := loadCache()
	key := repoKey(workDir)
	if t, ok := cache.Repos[key]; ok {
		return t, key, nil
	}
	if len(Themes) == 0 {
		return "", key, errors.New("no themes registered")
	}
	return Themes[cache.RotationIndex%len(Themes)], key, nil
}
