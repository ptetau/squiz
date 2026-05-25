package main

import (
	_ "embed"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

// embeddedExample is the canonical squiz fixture, shipped inside the
// binary so users can scaffold a starter without remembering the JSON
// schema or hunting for testdata/ paths after a release-archive install.
//
// MUST be kept byte-identical to testdata/smoke.json — there's a sync
// check test (TestEmbeddedExampleMatchesRoot) that fails if they drift.
//
//go:embed example/smoke.json
var embeddedExample []byte

// cmdExample writes the canonical sample. Three modes:
//
//	squiz example                  → write squiz-example.json to cwd
//	squiz example --out path.json  → write to a specific path
//	squiz example --stdout         → write to stdout
//
// After writing to a file, prints a one-line "try: squiz <path> --open"
// hint so the user knows what to do next.
func cmdExample(args []string) {
	fs := flag.NewFlagSet("example", flag.ExitOnError)
	out := fs.String("out", "", "where to write the example JSON (default: squiz-example.json in cwd)")
	stdout := fs.Bool("stdout", false, "write to stdout instead of a file")
	args = reorderFlagsFirst(args, map[string]bool{"stdout": true})
	if err := fs.Parse(args); err != nil {
		os.Exit(2)
	}

	if *stdout {
		if _, err := os.Stdout.Write(embeddedExample); err != nil {
			fmt.Fprintf(os.Stderr, "write stdout: %v\n", err)
			os.Exit(1)
		}
		return
	}

	outPath := *out
	if outPath == "" {
		outPath = "squiz-example.json"
	}
	absOut, err := filepath.Abs(outPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "resolve output path: %v\n", err)
		os.Exit(1)
	}
	if err := os.MkdirAll(filepath.Dir(absOut), 0755); err != nil {
		fmt.Fprintf(os.Stderr, "make output dir: %v\n", err)
		os.Exit(1)
	}
	if err := os.WriteFile(absOut, embeddedExample, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "write example: %v\n", err)
		os.Exit(1)
	}
	fmt.Fprintln(os.Stderr, "wrote", absOut)
	fmt.Fprintln(os.Stderr, "try:  squiz", filepath.Base(absOut), "--open")
}
