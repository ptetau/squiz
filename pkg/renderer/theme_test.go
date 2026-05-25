package renderer

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// isolateHome redirects user-home lookups to a per-test temp dir so the
// real ~/.squiz/themes.json is never touched.
func isolateHome(t *testing.T) string {
	t.Helper()
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)
	t.Setenv("USERPROFILE", tmp)
	return tmp
}

func TestResolveTheme_OverrideWins(t *testing.T) {
	home := isolateHome(t)
	workDir := t.TempDir()

	got := ResolveTheme(workDir, "amber")
	if got != "amber" {
		t.Errorf("ResolveTheme(_, \"amber\") = %q, want %q", got, "amber")
	}

	// Override must skip the cache path entirely — no themes.json file written.
	cachePath := filepath.Join(home, ".squiz", "themes.json")
	if _, err := os.Stat(cachePath); !os.IsNotExist(err) {
		t.Errorf("themes.json exists after override-only call (err=%v); override should skip cache", err)
	}
}

func TestResolveTheme_InvalidOverrideFallsThrough(t *testing.T) {
	isolateHome(t)
	workDir := t.TempDir()

	got := ResolveTheme(workDir, "nonsense")
	if !validTheme(got) {
		t.Errorf("ResolveTheme(_, \"nonsense\") = %q, not in Themes; should fall through to rotation", got)
	}
}

func TestResolveTheme_PersistsAndStable(t *testing.T) {
	home := isolateHome(t)
	workDir := t.TempDir()

	first := ResolveTheme(workDir, "")
	if !validTheme(first) {
		t.Fatalf("first call returned invalid theme %q", first)
	}

	second := ResolveTheme(workDir, "")
	if second != first {
		t.Errorf("second call returned %q, want stable %q", second, first)
	}

	cachePath := filepath.Join(home, ".squiz", "themes.json")
	data, err := os.ReadFile(cachePath)
	if err != nil {
		t.Fatalf("themes.json not written: %v", err)
	}

	var c themeCache
	if err := json.Unmarshal(data, &c); err != nil {
		t.Fatalf("themes.json not valid JSON: %v", err)
	}
	key := repoKey(workDir)
	gotTheme, ok := c.Repos[key]
	if !ok {
		t.Fatalf("themes.json has no entry for key %q; repos=%v", key, c.Repos)
	}
	if gotTheme != first {
		t.Errorf("themes.json[%q] = %q, want %q", key, gotTheme, first)
	}
}

func TestResolveTheme_DifferentWorkdirsRotate(t *testing.T) {
	isolateHome(t)

	results := make([]string, 3)
	for i := 0; i < 3; i++ {
		wd := t.TempDir()
		results[i] = ResolveTheme(wd, "")
		if !validTheme(results[i]) {
			t.Fatalf("rotation slot %d returned invalid theme %q", i, results[i])
		}
	}

	// Sequential workDirs in a clean cache should advance the rotation index
	// and yield three distinct themes (Themes has 8 entries, well over 3).
	if results[0] == results[1] || results[1] == results[2] || results[0] == results[2] {
		t.Errorf("expected three distinct themes from rotation, got %v", results)
	}

	// Sanity: in a fresh cache the first slot should be Themes[0].
	if results[0] != Themes[0] {
		t.Errorf("first rotation slot = %q, want %q", results[0], Themes[0])
	}
}

func TestValidTheme(t *testing.T) {
	if len(Themes) != 8 {
		t.Errorf("Themes has %d entries, expected 8 (test needs updating if list changed)", len(Themes))
	}
	for _, name := range Themes {
		if !validTheme(name) {
			t.Errorf("validTheme(%q) = false, want true (it's in Themes)", name)
		}
	}
	if validTheme("definitely-not-a-theme") {
		t.Error("validTheme(\"definitely-not-a-theme\") = true, want false")
	}
	if validTheme("") {
		t.Error("validTheme(\"\") = true, want false")
	}
}
