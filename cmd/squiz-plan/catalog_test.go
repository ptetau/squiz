package main

import (
	"encoding/json"
	"os/exec"
	"strings"
	"testing"
)

// TestMain_CatalogWF runs `squiz-plan catalog wf` and asserts known names appear.
func TestMain_CatalogWF(t *testing.T) {
	bin := buildBinary(t)

	cmd := exec.Command(bin, "catalog", "wf")
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("squiz-plan catalog wf failed: %v", err)
	}
	s := string(out)
	for _, want := range []string{"calendar-grid", "spark-rising"} {
		if !strings.Contains(s, want) {
			t.Errorf("catalog wf output missing %q\n---\n%s", want, s)
		}
	}
}

// TestMain_CatalogWF_JSON asserts --json parses as an array ≥ 50.
func TestMain_CatalogWF_JSON(t *testing.T) {
	bin := buildBinary(t)

	cmd := exec.Command(bin, "catalog", "wf", "--json")
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("squiz-plan catalog wf --json failed: %v", err)
	}
	var entries []map[string]any
	if err := json.Unmarshal(out, &entries); err != nil {
		t.Fatalf("catalog wf --json not parseable: %v\n---\n%s", err, out)
	}
	if len(entries) < 50 {
		t.Errorf("catalog wf --json returned %d entries, want ≥ 50", len(entries))
	}
}

// TestMain_CatalogArch sanity-checks the arch catalog.
func TestMain_CatalogArch(t *testing.T) {
	bin := buildBinary(t)

	cmd := exec.Command(bin, "catalog", "arch")
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("squiz-plan catalog arch failed: %v", err)
	}
	s := string(out)
	for _, want := range []string{"server", "database", "queue"} {
		if !strings.Contains(s, want) {
			t.Errorf("catalog arch output missing %q\n---\n%s", want, s)
		}
	}
}

// TestMain_CatalogThemes asserts all 8 ship themes appear.
func TestMain_CatalogThemes(t *testing.T) {
	bin := buildBinary(t)

	cmd := exec.Command(bin, "catalog", "themes")
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("squiz-plan catalog themes failed: %v", err)
	}
	s := string(out)
	for _, want := range []string{"paper", "phosphor", "amber", "beige", "rose", "ocean", "forest", "slate"} {
		if !strings.Contains(s, want) {
			t.Errorf("catalog themes output missing %q\n---\n%s", want, s)
		}
	}
}

// TestMain_CatalogDSL asserts every primitive is listed.
func TestMain_CatalogDSL(t *testing.T) {
	bin := buildBinary(t)

	cmd := exec.Command(bin, "catalog", "dsl")
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("squiz-plan catalog dsl failed: %v", err)
	}
	s := string(out)
	for _, want := range []string{"grid:", "spark:", "bars:", "swatches:", "pills:", "sample:", "circle-pack:", "text:", "flow:", "box:", "arrow:"} {
		if !strings.Contains(s, want) {
			t.Errorf("catalog dsl output missing %q\n---\n%s", want, s)
		}
	}
}
