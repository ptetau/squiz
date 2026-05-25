package renderer

import (
	"strings"
	"testing"
)

// TestArchLibrary_AllNamesResolve walks every entry in ArchLibrary and
// asserts ArchIcon returns a non-empty <g> snippet. Catches drift if
// someone adds a name without a body.
func TestArchLibrary_AllNamesResolve(t *testing.T) {
	if len(ArchLibrary) == 0 {
		t.Fatal("ArchLibrary is empty — expected ~30 entries")
	}
	for name, body := range ArchLibrary {
		t.Run(name, func(t *testing.T) {
			got := ArchIcon(name)
			if got == "" {
				t.Errorf("ArchIcon(%q) returned empty SVG snippet", name)
			}
			if !strings.Contains(got, "<g") {
				t.Errorf("ArchIcon(%q) does not contain a <g element: %s", name, got)
			}
			if got != body {
				t.Errorf("ArchIcon(%q) differed from registry body", name)
			}
		})
	}
}

// TestArchLibrary_HasExpectedCoreNames spot-checks names the spec lists as
// foundational. Catches accidental rename/removal.
func TestArchLibrary_HasExpectedCoreNames(t *testing.T) {
	mustHave := []string{
		// compute
		"server", "container", "pod", "function", "worker", "scheduler",
		// data
		"database", "table", "blob", "storage", "cache", "stream",
		// network
		"load-balancer", "gateway", "cdn", "dns", "vpc", "subnet", "firewall",
		// services
		"api", "queue", "topic",
		// observability
		"log", "metric", "trace",
		// identity
		"user", "mobile", "browser",
		// security
		"secret", "key-icon",
	}
	for _, name := range mustHave {
		if _, ok := ArchLibrary[name]; !ok {
			t.Errorf("ArchLibrary missing expected entry %q", name)
		}
	}
}

// TestArchLibrary_ExactlyThirty pins the count to the spec target.
func TestArchLibrary_ExactlyThirty(t *testing.T) {
	if got, want := len(ArchLibrary), 30; got != want {
		t.Errorf("ArchLibrary size = %d, want %d", got, want)
	}
}

// TestArchLibrary_UnknownName — ArchIcon for an unknown name returns "".
func TestArchLibrary_UnknownName(t *testing.T) {
	if got := ArchIcon("does-not-exist"); got != "" {
		t.Errorf(`ArchIcon("does-not-exist") = %q, want ""`, got)
	}
}

// TestRenderArt_ArchPrefix exercises the dispatcher's `arch:` branch.
func TestRenderArt_ArchPrefix(t *testing.T) {
	// Valid name → wrapped SVG.
	svg, hidden := RenderArt("arch:server", 0)
	if hidden {
		t.Errorf("RenderArt(arch:server) hidden = true, want false")
	}
	if svg == "" {
		t.Fatal("RenderArt(arch:server) returned empty SVG")
	}
	if !strings.Contains(svg, "<svg") {
		t.Errorf("RenderArt(arch:server) output not wrapped in <svg>: %s", svg)
	}
	if !strings.Contains(svg, "<g") {
		t.Errorf("RenderArt(arch:server) missing <g icon body: %s", svg)
	}

	// Unknown name → error placeholder.
	bad, badHidden := RenderArt("arch:nonsense", 0)
	if badHidden {
		t.Errorf("RenderArt(arch:nonsense) hidden = true, want false")
	}
	if !strings.Contains(bad, "<svg") {
		t.Fatalf("RenderArt(arch:nonsense) not an SVG: %s", bad)
	}
	if !strings.Contains(bad, "unknown arch") {
		t.Errorf("expected 'unknown arch' marker in placeholder, got: %s", bad)
	}
	if !strings.Contains(bad, "nonsense") {
		t.Errorf("expected unknown name 'nonsense' echoed in placeholder, got: %s", bad)
	}
}
