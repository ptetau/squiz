package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// schemaJSON is the embedded JSON Schema (draft 2020-12) for the squiz-plan
// index.json input. Sourced from cmd/squiz-plan/schema/squiz-plan.schema.json;
// one copy of truth for both the `squiz-plan schema` subcommand and the
// in-package drift test (TestSchemaCoversAllStructFields).
//
//go:embed schema/squiz-plan.schema.json
var schemaJSON []byte

// cmdSchema writes the embedded JSON Schema to stdout (default) or to a
// path passed via --out. Output is re-indented through json.Indent so
// the bytes on the wire are always pretty regardless of how the source
// file is formatted on disk.
//
//	squiz-plan schema                  → pretty JSON on stdout
//	squiz-plan schema --out file.json  → pretty JSON written to file.json
func cmdSchema(args []string) {
	fs := flag.NewFlagSet("schema", flag.ExitOnError)
	out := fs.String("out", "", "write schema to this path instead of stdout")
	args = reorderFlagsFirst(args, map[string]bool{})
	if err := fs.Parse(args); err != nil {
		os.Exit(2)
	}

	pretty, err := prettyJSON(schemaJSON)
	if err != nil {
		fmt.Fprintf(os.Stderr, "schema: re-indent embedded schema: %v\n", err)
		os.Exit(1)
	}

	if *out == "" {
		if _, err := io.WriteString(os.Stdout, pretty); err != nil {
			fmt.Fprintf(os.Stderr, "schema: write stdout: %v\n", err)
			os.Exit(1)
		}
		fmt.Println()
		return
	}

	absOut, err := filepath.Abs(*out)
	if err != nil {
		fmt.Fprintf(os.Stderr, "schema: resolve output path: %v\n", err)
		os.Exit(1)
	}
	if err := os.MkdirAll(filepath.Dir(absOut), 0755); err != nil {
		fmt.Fprintf(os.Stderr, "schema: make output dir: %v\n", err)
		os.Exit(1)
	}
	if err := os.WriteFile(absOut, []byte(pretty), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "schema: write %s: %v\n", absOut, err)
		os.Exit(1)
	}
	fmt.Fprintln(os.Stderr, "wrote", absOut)
}

// prettyJSON re-indents raw JSON bytes with 2-space indent. Returns the
// json.Indent error verbatim so the caller can surface a useful message
// if the embedded schema somehow shipped malformed.
func prettyJSON(raw []byte) (string, error) {
	var buf bytes.Buffer
	if err := json.Indent(&buf, raw, "", "  "); err != nil {
		return "", err
	}
	return buf.String(), nil
}
