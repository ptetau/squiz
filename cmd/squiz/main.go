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

// version is overridden at release time via -ldflags "-X main.version=…".
var version = "dev"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(2)
	}

	switch os.Args[1] {
	case "render":
		cmdRender(os.Args[2:])
	case "version", "--version", "-v":
		fmt.Println("squiz " + version)
	case "help", "--help", "-h":
		printUsage()
	default:
		// Shorthand: `squiz foo.json [flags…]` == `squiz render foo.json [flags…] --open`.
		// Any flags the user already passed are forwarded verbatim; --open is
		// only auto-appended when the user didn't ask for it explicitly.
		if strings.HasSuffix(strings.ToLower(os.Args[1]), ".json") {
			args := append([]string{}, os.Args[1:]...)
			if !hasFlag(args, "open") {
				args = append(args, "--open")
			}
			cmdRender(args)
			return
		}
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", os.Args[1])
		printUsage()
		os.Exit(2)
	}
}

func cmdRender(args []string) {
	fs := flag.NewFlagSet("render", flag.ExitOnError)
	out := fs.String("out", "", "output HTML path (default: <input>.html next to input)")
	stdout := fs.Bool("stdout", false, "write HTML to stdout instead of a file")
	open := fs.Bool("open", false, "open the rendered HTML in the default browser")
	theme := fs.String("theme", "", "force theme (paper|phosphor|amber|beige|rose|ocean|forest|slate)")
	// Reorder so flags precede positionals — Go's stdlib flag.Parse stops at
	// the first non-flag arg, so without this `squiz render foo.json --out x`
	// would silently drop --out. boolFlags lists flags that don't consume
	// the following arg as a value.
	args = reorderFlagsFirst(args, map[string]bool{"stdout": true, "open": true})
	if err := fs.Parse(args); err != nil {
		os.Exit(2)
	}

	if fs.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "render: missing input JSON path")
		os.Exit(2)
	}

	inputPath := fs.Arg(0)
	absInput, err := filepath.Abs(inputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "resolve input path: %v\n", err)
		os.Exit(1)
	}
	data, err := os.ReadFile(absInput)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read input: %v\n", err)
		os.Exit(1)
	}

	doc, err := renderer.ParseDocument(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse input: %v\n", err)
		os.Exit(1)
	}

	// Compute the output path BEFORE render so the renderer can embed
	// the absolute path into the page for self-referential anchors.
	// Deterministic: always <basename>.html next to the input unless
	// the user passes --out.
	outPath := *out
	if outPath == "" {
		base := strings.TrimSuffix(filepath.Base(absInput), filepath.Ext(absInput))
		outPath = filepath.Join(filepath.Dir(absInput), base+".html")
	}
	absOut, err := filepath.Abs(outPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "resolve output path: %v\n", err)
		os.Exit(1)
	}

	html, err := renderer.Render(doc, renderer.RenderOpts{
		OutputPath:    absOut,
		ThemeOverride: *theme,
		WorkDir:       filepath.Dir(absInput),
		Version:       version,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "render: %v\n", err)
		os.Exit(1)
	}

	if *stdout {
		io.WriteString(os.Stdout, html)
		return
	}

	if err := os.WriteFile(absOut, []byte(html), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "write output: %v\n", err)
		os.Exit(1)
	}
	fmt.Fprintln(os.Stderr, "wrote", absOut, "·", doc.Theme)

	if *open {
		if err := OpenInBrowser(absOut); err != nil {
			fmt.Fprintf(os.Stderr, "open browser: %v\n", err)
		}
	}
}

// reorderFlagsFirst walks args and returns a new slice where all flags come
// before all positionals, preserving original order within each group.
// Treats `-x`, `--x`, `--x=v` as flags; for the `--x v` form, the value
// comes along unless `x` is in boolFlags (which take no value). The literal
// `--` terminator passes through and stops flag-detection (POSIX convention).
func reorderFlagsFirst(args []string, boolFlags map[string]bool) []string {
	flags := make([]string, 0, len(args))
	positionals := make([]string, 0, len(args))
	seenTerminator := false

	for i := 0; i < len(args); i++ {
		a := args[i]
		if seenTerminator || !strings.HasPrefix(a, "-") || a == "-" {
			positionals = append(positionals, a)
			continue
		}
		if a == "--" {
			seenTerminator = true
			flags = append(flags, a)
			continue
		}
		flags = append(flags, a)
		// --name=value already carries its value; nothing to consume.
		if strings.Contains(a, "=") {
			continue
		}
		name := strings.TrimLeft(a, "-")
		if boolFlags[name] {
			continue
		}
		// String flag without `=` consumes the next arg as its value.
		if i+1 < len(args) {
			flags = append(flags, args[i+1])
			i++
		}
	}
	return append(flags, positionals...)
}

// hasFlag reports whether args contains `-name`, `--name`, `-name=…`, or
// `--name=…`. Cheap presence check used to avoid double-appending flags
// the user already specified.
func hasFlag(args []string, name string) bool {
	for _, a := range args {
		if a == "-"+name || a == "--"+name {
			return true
		}
		if strings.HasPrefix(a, "-"+name+"=") || strings.HasPrefix(a, "--"+name+"=") {
			return true
		}
	}
	return false
}

func printUsage() {
	fmt.Fprintln(os.Stderr, `squiz `+version+` — render a Squiz spec from JSON to interactive HTML

Usage:
  squiz render <input.json> [--out path] [--stdout] [--open] [--theme name]
  squiz <input.json> [flags…]           (shorthand: render + open; flags forward)
  squiz version

Flags may appear before or after the input path.

Examples:
  squiz habits.json                     render + open in browser
  squiz habits.json --theme phosphor    shorthand with extra flags
  squiz render spec.json --out doc.html
  squiz render spec.json --stdout > doc.html`)
}
