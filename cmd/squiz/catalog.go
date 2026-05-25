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

// cmdCatalog implements `squiz catalog [name] [--json] [--previews [--out f]]`.
//
//	squiz catalog                  # list catalog names (wf, arch, dsl, themes)
//	squiz catalog wf               # text listing
//	squiz catalog wf --json        # JSON array
//	squiz catalog wf --previews    # write gallery HTML (wf-gallery.html in cwd)
//	squiz catalog arch             # same shape
//	squiz catalog dsl              # primitives + grammar
//	squiz catalog themes           # 8 themes + vibe lines
func cmdCatalog(args []string) {
	fs := flag.NewFlagSet("catalog", flag.ExitOnError)
	asJSON := fs.Bool("json", false, "emit machine-readable JSON")
	previews := fs.Bool("previews", false, "write a self-contained gallery HTML (wf/arch only)")
	out := fs.String("out", "", "output path for --previews (default: <name>-gallery.html in cwd)")
	theme := fs.String("theme", "paper", "theme for the gallery page")
	// Flag-anywhere parity with cmdRender: bool flags don't consume the next arg.
	args = reorderFlagsFirst(args, map[string]bool{"json": true, "previews": true})
	if err := fs.Parse(args); err != nil {
		os.Exit(2)
	}

	if fs.NArg() == 0 {
		// No catalog name: list the available catalogs.
		if *asJSON {
			io.WriteString(os.Stdout, `["wf","arch","dsl","themes"]`+"\n")
			return
		}
		io.WriteString(os.Stdout, renderer.FormatNamesText())
		return
	}

	name := fs.Arg(0)
	switch name {
	case "wf":
		emitCatalogEntries("wf", renderer.WFCatalog(), renderer.WFRender, *asJSON, *previews, *out, *theme)
	case "arch":
		emitCatalogEntries("arch", renderer.ArchCatalog(), renderer.ArchRender, *asJSON, *previews, *out, *theme)
	case "dsl":
		if *previews {
			fmt.Fprintln(os.Stderr, "catalog dsl: --previews not supported (DSL primitives are dynamic)")
			os.Exit(2)
		}
		emitDSL(*asJSON)
	case "themes":
		if *previews {
			fmt.Fprintln(os.Stderr, "catalog themes: --previews not supported")
			os.Exit(2)
		}
		emitThemes(*asJSON)
	default:
		fmt.Fprintf(os.Stderr, "catalog: unknown name %q (want: wf, arch, dsl, themes)\n", name)
		os.Exit(2)
	}
}

func emitCatalogEntries(name string, entries []renderer.CatalogEntry, render func(string) string, asJSON, previews bool, out, theme string) {
	if previews {
		path := out
		if path == "" {
			path = name + "-gallery.html"
		}
		abs, err := filepath.Abs(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "resolve --out: %v\n", err)
			os.Exit(1)
		}
		if err := os.MkdirAll(filepath.Dir(abs), 0o755); err != nil {
			fmt.Fprintf(os.Stderr, "mkdir output dir: %v\n", err)
			os.Exit(1)
		}
		title := strings.ToUpper(name) + " · catalog"
		html := renderer.RenderGalleryHTML(title, theme, entries, render)
		if err := os.WriteFile(abs, []byte(html), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "write gallery: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintln(os.Stderr, "wrote", abs)
		return
	}
	if asJSON {
		s, err := renderer.FormatCatalogJSON(entries)
		if err != nil {
			fmt.Fprintf(os.Stderr, "json: %v\n", err)
			os.Exit(1)
		}
		io.WriteString(os.Stdout, s)
		return
	}
	io.WriteString(os.Stdout, renderer.FormatCatalogText(entries))
}

func emitDSL(asJSON bool) {
	if asJSON {
		s, err := renderer.FormatDSLJSON()
		if err != nil {
			fmt.Fprintf(os.Stderr, "json: %v\n", err)
			os.Exit(1)
		}
		io.WriteString(os.Stdout, s)
		return
	}
	io.WriteString(os.Stdout, renderer.FormatDSLText())
}

func emitThemes(asJSON bool) {
	if asJSON {
		s, err := renderer.FormatThemesJSON()
		if err != nil {
			fmt.Fprintf(os.Stderr, "json: %v\n", err)
			os.Exit(1)
		}
		io.WriteString(os.Stdout, s)
		return
	}
	io.WriteString(os.Stdout, renderer.FormatThemesText())
}
