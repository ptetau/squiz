package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/ptetau/squiz/pkg/renderer"
)

// cmdPreview implements `squiz-plan preview <spec> [--theme T] [--out F] [--stdout]`.
//
// Same shape as `squiz preview` — squiz-plan agents author plans that
// reference the same wf/arch/dsl art forms, so this lets them sanity-check
// a single form without booting the whole plan renderer.
//
//	squiz-plan preview wf:calendar-grid                    # → wf-calendar-grid.html
//	squiz-plan preview arch:server --theme phosphor
//	squiz-plan preview "flow:[client,api,db]" --out p.html
//	squiz-plan preview wf:calendar-grid --stdout
func cmdPreview(args []string) {
	fs := flag.NewFlagSet("preview", flag.ExitOnError)
	out := fs.String("out", "", "output HTML path (default: derived from spec)")
	stdout := fs.Bool("stdout", false, "write HTML to stdout instead of a file")
	theme := fs.String("theme", "paper", "theme (paper|phosphor|amber|beige|rose|ocean|forest|slate)")
	args = reorderFlagsFirst(args, map[string]bool{"stdout": true})
	if err := fs.Parse(args); err != nil {
		os.Exit(2)
	}

	if fs.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "preview: missing spec (e.g. wf:calendar-grid)")
		os.Exit(2)
	}

	spec := fs.Arg(0)
	svg, hidden := renderer.RenderArt(spec, 0)
	if hidden {
		fmt.Fprintln(os.Stderr, "preview: spec resolves to 'no art' (nothing to render)")
		os.Exit(1)
	}

	html := renderer.RenderPreviewHTML(spec, svg, *theme)

	if *stdout {
		io.WriteString(os.Stdout, html)
		return
	}

	outPath := *out
	if outPath == "" {
		outPath = previewDefaultFilename(spec)
	}
	abs, err := filepath.Abs(outPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "resolve --out: %v\n", err)
		os.Exit(1)
	}
	if err := os.MkdirAll(filepath.Dir(abs), 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "mkdir output dir: %v\n", err)
		os.Exit(1)
	}
	if err := os.WriteFile(abs, []byte(html), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "write output: %v\n", err)
		os.Exit(1)
	}
	fmt.Fprintln(os.Stderr, "wrote", abs)
}

// previewDefaultFilename turns a spec into a stable .html filename in cwd.
// Mirrors cmd/squiz/preview.go's helper — kept local per the codebase's
// "binaries are self-contained" convention (see reorderFlagsFirst note).
func previewDefaultFilename(spec string) string {
	s := strings.ToLower(spec)
	safe := strings.Builder{}
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			safe.WriteRune(r)
		default:
			safe.WriteRune('-')
		}
	}
	name := strings.Trim(safe.String(), "-")
	for strings.Contains(name, "--") {
		name = strings.ReplaceAll(name, "--", "-")
	}
	if name == "" {
		name = "preview"
	}
	return name + ".html"
}
