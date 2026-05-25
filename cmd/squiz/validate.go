package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"regexp"

	"github.com/ptetau/squiz/pkg/renderer"
)

// validateError is one machine-readable problem found during validation.
// Path uses a JSON-pointer-like dotted form (e.g. squizzes[2].options[1].id)
// to point at the offending value. Message is a single-sentence human
// summary suitable for direct CLI output.
type validateError struct {
	Path    string `json:"path"`
	Message string `json:"message"`
}

// validateReport is the JSON shape printed when --json is set. The
// Errors slice is never nil — always initialised to [] so JS consumers
// don't have to handle null.
type validateReport struct {
	Valid    bool            `json:"valid"`
	Errors   []validateError `json:"errors"`
	Warnings []validateError `json:"warnings,omitempty"`
}

// markerRE mirrors renderer.markerRE — see pkg/renderer/render.go. Duplicated
// here so we don't have to widen the renderer's public API just for this
// validator. {{squizId}} marker pattern; ID is [a-zA-Z][a-zA-Z0-9_-]*.
var markerRE = regexp.MustCompile(`\{\{([a-zA-Z][a-zA-Z0-9_-]*)\}\}`)

// cmdValidate parses an input JSON file, runs the renderer's shape check
// (ParseDocument), then layers on a handful of business rules the parser
// doesn't currently enforce:
//
//   - Squiz `id` uniqueness across the document
//   - Option `id` uniqueness within each squiz; non-empty ids
//   - Theme name (if set) must be one of the 8 known names
//   - Density (if set) must be `compact` or `comfortable`
//   - `{{markers}}` in spec.paragraphs that reference unknown squizzes
//     are WARNINGS (not errors) — the renderer treats them as non-fatal
//     and emits them as literal text so authors notice typos. They never
//     change exit code.
//
// Exit code: 0 on success (warnings allowed), 1 on any error.
//
//	squiz validate <input.json>          → text mode
//	squiz validate <input.json> --json   → machine-readable report
func cmdValidate(args []string) {
	fs := flag.NewFlagSet("validate", flag.ExitOnError)
	jsonOut := fs.Bool("json", false, "emit a JSON report instead of text")
	args = reorderFlagsFirst(args, map[string]bool{"json": true})
	if err := fs.Parse(args); err != nil {
		os.Exit(2)
	}

	if fs.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "validate: missing input JSON path")
		os.Exit(2)
	}
	inputPath := fs.Arg(0)
	data, err := os.ReadFile(inputPath)
	if err != nil {
		emitValidateFatal(*jsonOut, "", fmt.Sprintf("read input: %v", err))
		os.Exit(1)
	}

	doc, parseErr := renderer.ParseDocument(data)
	if parseErr != nil {
		emitValidateFatal(*jsonOut, "", fmt.Sprintf("parse: %v", parseErr))
		os.Exit(1)
	}

	errors, warnings, squizCount, optionCount := validateDocument(doc)

	// Layer the v0.8.0 taste lints on top of the hard-error pass. These
	// are warnings only — they never affect the exit code. They use the
	// `warn:` prefix in text mode (distinct from the existing `warning:`
	// prefix for marker warnings) so authors can grep them out separately.
	lintWarnings := runLints(doc)

	emitValidate(*jsonOut, errors, warnings, lintWarnings, validateCounts{
		squizzes: squizCount,
		options:  optionCount,
	})
	if len(errors) > 0 {
		os.Exit(1)
	}
}

type validateCounts struct {
	squizzes int
	options  int
}

// validateDocument runs the post-ParseDocument business-rule checks and
// returns the categorised findings + counts for the summary line.
func validateDocument(doc *renderer.Document) (errs, warns []validateError, squizCount, optCount int) {
	// Theme check. Empty is fine (renderer auto-rotates).
	if doc.Theme != "" && !knownTheme(doc.Theme) {
		errs = append(errs, validateError{
			Path:    "theme",
			Message: fmt.Sprintf("unknown theme %q (must be one of: %s)", doc.Theme, themeListForMsg()),
		})
	}

	// Density check. ParseDocument defaults blank → "compact", so by the
	// time we see it here it should be one of the two valid values; any
	// other value is something the author put in the file deliberately.
	if doc.Density != "" && doc.Density != "compact" && doc.Density != "comfortable" {
		errs = append(errs, validateError{
			Path:    "density",
			Message: fmt.Sprintf("unknown density %q (must be \"compact\" or \"comfortable\")", doc.Density),
		})
	}

	// Squiz ID uniqueness + per-squiz option checks.
	squizIDs := make(map[string]int, len(doc.Squizzes))
	for i, s := range doc.Squizzes {
		base := fmt.Sprintf("squizzes[%d]", i)
		if s.ID == "" {
			errs = append(errs, validateError{
				Path:    base + ".id",
				Message: "empty (must be a stable slug)",
			})
		} else if prev, dup := squizIDs[s.ID]; dup {
			errs = append(errs, validateError{
				Path:    base + ".id",
				Message: fmt.Sprintf("duplicate id %q (first declared at squizzes[%d])", s.ID, prev),
			})
		} else {
			squizIDs[s.ID] = i
		}

		// Option IDs unique within this squiz.
		optIDs := make(map[string]int, len(s.Options))
		for j, o := range s.Options {
			optBase := fmt.Sprintf("%s.options[%d]", base, j)
			if o.ID == "" {
				errs = append(errs, validateError{
					Path:    optBase + ".id",
					Message: "empty (must be a stable slug)",
				})
				continue
			}
			if prev, dup := optIDs[o.ID]; dup {
				errs = append(errs, validateError{
					Path:    optBase + ".id",
					Message: fmt.Sprintf("duplicate id %q within this squiz (first at options[%d])", o.ID, prev),
				})
				continue
			}
			optIDs[o.ID] = j
		}

		optCount += len(s.Options)
	}
	squizCount = len(doc.Squizzes)

	// Marker scan in spec paragraphs — warnings only.
	for i, p := range doc.Spec.Paragraphs {
		for _, m := range markerRE.FindAllStringSubmatch(p.Text, -1) {
			id := m[1]
			if _, ok := squizIDs[id]; !ok {
				warns = append(warns, validateError{
					Path:    fmt.Sprintf("spec.paragraphs[%d].text", i),
					Message: fmt.Sprintf("marker {{%s}} references unknown squiz (will render as literal text)", id),
				})
			}
		}
	}

	return errs, warns, squizCount, optCount
}

// knownTheme reports whether name is one of the 8 ship-with themes.
// Delegates to renderer.Themes so adding a theme there propagates here.
func knownTheme(name string) bool {
	for _, t := range renderer.Themes {
		if t == name {
			return true
		}
	}
	return false
}

func themeListForMsg() string {
	out := ""
	for i, t := range renderer.Themes {
		if i > 0 {
			out += ", "
		}
		out += t
	}
	return out
}

// emitValidateFatal renders a single error (typically: file unreadable or
// stdlib parse failure) in whichever mode the user asked for. Kept
// separate so the success path doesn't have to special-case the
// "no document loaded" state.
func emitValidateFatal(asJSON bool, path, msg string) {
	if asJSON {
		report := validateReport{
			Valid:  false,
			Errors: []validateError{{Path: path, Message: msg}},
		}
		buf, _ := json.MarshalIndent(report, "", "  ")
		fmt.Println(string(buf))
		return
	}
	if path != "" {
		fmt.Fprintf(os.Stderr, "%s: %s\n", path, msg)
	} else {
		fmt.Fprintln(os.Stderr, msg)
	}
	fmt.Fprintln(os.Stderr, "invalid (1 error)")
}

// emitValidate writes the validation result in either text or JSON.
// Text mode: one finding per line as `path: message`. Marker warnings get
// the `warning: ` prefix (pre-v0.8.0 behaviour). Lint warnings get the
// `warn: ` prefix (v0.8.0 taste lints). Both are folded into the JSON
// report's single `warnings` array. Ends with a 1-line summary on
// stdout/stderr depending on success/failure. Exits cleanly — the caller
// decides exit code from len(errs).
func emitValidate(asJSON bool, errs, warns, lints []validateError, counts validateCounts) {
	if asJSON {
		combined := make([]validateError, 0, len(warns)+len(lints))
		combined = append(combined, warns...)
		combined = append(combined, lints...)
		report := validateReport{
			Valid:    len(errs) == 0,
			Errors:   errs,
			Warnings: combined,
		}
		if report.Errors == nil {
			report.Errors = []validateError{}
		}
		buf, _ := json.MarshalIndent(report, "", "  ")
		fmt.Println(string(buf))
		return
	}

	w := os.Stderr
	if len(errs) == 0 {
		w = os.Stdout
	}
	for _, e := range errs {
		fmt.Fprintf(os.Stderr, "%s: %s\n", e.Path, e.Message)
	}
	for _, e := range warns {
		fmt.Fprintf(os.Stderr, "warning: %s: %s\n", e.Path, e.Message)
	}
	for _, e := range lints {
		fmt.Fprintf(os.Stderr, "warn: %s: %s\n", e.Path, e.Message)
	}
	if len(errs) == 0 {
		fmt.Fprintf(w, "valid (%d %s, %d option%s)\n",
			counts.squizzes, squizPlural(counts.squizzes),
			counts.options, plural(counts.options))
	} else {
		fmt.Fprintf(w, "invalid (%d error%s)\n", len(errs), plural(len(errs)))
	}
}

func plural(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}

// squizPlural handles the irregular pluralisation "squiz" → "squizzes".
// Mirrors how the renderer talks about decisions in the topbar.
func squizPlural(n int) string {
	if n == 1 {
		return "squiz"
	}
	return "squizzes"
}
