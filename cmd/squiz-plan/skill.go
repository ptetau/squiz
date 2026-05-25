package main

import (
	_ "embed"
	"flag"
	"fmt"
	"os"
)

// embeddedSkill is the canonical SKILL.md content baked into the binary at
// compile time. The source of truth lives at skills/squiz-plan/SKILL.md; we
// keep a literal copy at cmd/squiz-plan/skill.md because //go:embed cannot
// reach above the package directory. A sync-check test (skill_test.go)
// asserts the two files stay byte-identical; if it fails, run:
//
//	cp skills/squiz-plan/SKILL.md cmd/squiz-plan/skill.md
//
//go:embed skill.md
var embeddedSkill []byte

// cmdSkill dumps the embedded SKILL.md to stdout (or --out path).
//
// Usage:
//
//	squiz-plan skill              # write SKILL.md to stdout
//	squiz-plan skill --out PATH   # write SKILL.md to PATH
//
// No other flags. This is a content-dump command, not a renderer.
func cmdSkill(args []string) {
	fs := flag.NewFlagSet("skill", flag.ExitOnError)
	out := fs.String("out", "", "write SKILL.md to this path instead of stdout")
	if err := fs.Parse(args); err != nil {
		os.Exit(2)
	}

	if *out != "" {
		if err := os.WriteFile(*out, embeddedSkill, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "write skill: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintln(os.Stderr, "wrote", *out)
		return
	}

	if _, err := os.Stdout.Write(embeddedSkill); err != nil {
		fmt.Fprintf(os.Stderr, "write skill: %v\n", err)
		os.Exit(1)
	}
}
