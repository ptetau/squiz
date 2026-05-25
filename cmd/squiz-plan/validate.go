package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/ptetau/squiz/internal/planview"
)

// validateError is one machine-readable problem found during validation.
// Path uses a file path or `index`/`<section>.json` form to point at the
// source location; Message is a single-sentence human summary.
type validateError struct {
	Path    string `json:"path"`
	Message string `json:"message"`
}

// validateReport is the JSON shape printed when --json is set.
type validateReport struct {
	Valid  bool            `json:"valid"`
	Errors []validateError `json:"errors"`
}

// cmdValidate loads a plan via planview.LoadPlan (which already runs the
// full battery: index parse, section parse, prefix check, ID uniqueness,
// ref existence, option-id uniqueness) and reports any error returned.
//
// We delegate end-to-end rather than re-implementing the rules: the
// parser is the source of truth, this command just dresses its first
// error in a friendlier shell. (LoadPlan returns on first error, so
// callers see one finding at a time — there's no batched mode to wrap.)
//
// Exit code: 0 on success, 1 on any error.
//
//	squiz-plan validate <plan/index.json>        → text mode
//	squiz-plan validate <plan/index.json> --json → machine-readable report
func cmdValidate(args []string) {
	fs := flag.NewFlagSet("validate", flag.ExitOnError)
	jsonOut := fs.Bool("json", false, "emit a JSON report instead of text")
	args = reorderFlagsFirst(args, map[string]bool{"json": true})
	if err := fs.Parse(args); err != nil {
		os.Exit(2)
	}

	if fs.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "validate: missing input index.json path")
		os.Exit(2)
	}
	indexPath := fs.Arg(0)

	plan, loadErr := planview.LoadPlan(indexPath)
	if loadErr != nil {
		path, msg := splitLoadError(loadErr.Error())
		emitValidate(*jsonOut, []validateError{{Path: path, Message: msg}}, planCounts{})
		os.Exit(1)
	}

	emitValidate(*jsonOut, nil, planCounts{
		sections: len(plan.Sections),
		items:    countItems(plan),
	})
}

type planCounts struct {
	sections int
	items    int
}

func countItems(p *planview.Plan) int {
	n := 0
	for _, s := range p.Sections {
		n += len(s.Items)
	}
	return n
}

// splitLoadError takes a LoadPlan error string and best-effort separates
// the "path-ish" prefix from the human message. LoadPlan formats its
// errors as `<file-or-section>: <message>` (see parser.go:48–183), so
// splitting on the first `: ` gives us a clean (path, message) pair.
// Errors without that shape come back as ("", whole-string).
func splitLoadError(s string) (path, msg string) {
	if i := strings.Index(s, ": "); i > 0 {
		return s[:i], s[i+2:]
	}
	return "", s
}

// emitValidate writes the validation result in either text or JSON.
// Text mode: one finding per line as `path: message`, ending with a
// 1-line summary on a separate channel (stdout for success, stderr for
// failure). JSON mode: the report struct on stdout.
func emitValidate(asJSON bool, errs []validateError, counts planCounts) {
	if asJSON {
		report := validateReport{
			Valid:  len(errs) == 0,
			Errors: errs,
		}
		if report.Errors == nil {
			report.Errors = []validateError{}
		}
		buf, _ := json.MarshalIndent(report, "", "  ")
		fmt.Println(string(buf))
		return
	}

	for _, e := range errs {
		if e.Path != "" {
			fmt.Fprintf(os.Stderr, "%s: %s\n", e.Path, e.Message)
		} else {
			fmt.Fprintln(os.Stderr, e.Message)
		}
	}
	if len(errs) == 0 {
		fmt.Fprintf(os.Stdout, "valid (%d section%s, %d item%s)\n",
			counts.sections, plural(counts.sections),
			counts.items, plural(counts.items))
	} else {
		fmt.Fprintf(os.Stderr, "invalid (%d error%s)\n", len(errs), plural(len(errs)))
	}
}

func plural(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}
