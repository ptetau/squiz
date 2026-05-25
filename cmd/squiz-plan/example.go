package main

import (
	"embed"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

// embeddedPlan is the canonical squiz-plan fixture tree, shipped inside
// the binary so users can scaffold a starter without remembering the
// multi-file layout or hunting for testdata/ paths after a release-archive
// install.
//
// The 7 files MUST stay byte-identical to testdata/plan-example/* —
// there's a sync check test (TestEmbeddedPlanMatchesRoot) that fails if
// they drift.
//
//go:embed example/*.json
var embeddedPlan embed.FS

// cmdExample writes the canonical plan tree (7 files) into a directory.
// Two modes:
//
//	squiz-plan example                → write to ./squiz-plan-example/
//	squiz-plan example --out dir/     → write to a specific directory
//
// Auto-creates the directory if missing. After writing, prints a one-line
// "try: squiz-plan <dir>/index.json --open" hint.
func cmdExample(args []string) {
	fs := flag.NewFlagSet("example", flag.ExitOnError)
	out := fs.String("out", "", "directory to write the example plan into (default: squiz-plan-example/)")
	args = reorderFlagsFirst(args, map[string]bool{})
	if err := fs.Parse(args); err != nil {
		os.Exit(2)
	}

	outDir := *out
	if outDir == "" {
		outDir = "squiz-plan-example"
	}
	absOut, err := filepath.Abs(outDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "resolve output path: %v\n", err)
		os.Exit(1)
	}
	if err := os.MkdirAll(absOut, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "make output dir: %v\n", err)
		os.Exit(1)
	}

	entries, err := embeddedPlan.ReadDir("example")
	if err != nil {
		fmt.Fprintf(os.Stderr, "read embedded plan: %v\n", err)
		os.Exit(1)
	}
	written := 0
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		data, err := embeddedPlan.ReadFile("example/" + e.Name())
		if err != nil {
			fmt.Fprintf(os.Stderr, "read embedded %s: %v\n", e.Name(), err)
			os.Exit(1)
		}
		dst := filepath.Join(absOut, e.Name())
		if err := os.WriteFile(dst, data, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "write %s: %v\n", dst, err)
			os.Exit(1)
		}
		written++
	}

	fmt.Fprintf(os.Stderr, "wrote %d files into %s\n", written, absOut)
	fmt.Fprintf(os.Stderr, "try:  squiz-plan %s --open\n",
		filepath.Join(filepath.Base(absOut), "index.json"))
}
