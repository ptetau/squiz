package main

import (
	"fmt"
	"strings"
)

// stringBuilder is a tiny strings.Builder wrapper used by wf.go helpers
// that need format-and-append. Keeps the registry's helper functions
// terse and consistent.
type stringBuilder struct {
	b strings.Builder
}

func (s *stringBuilder) Appendf(format string, args ...any) {
	fmt.Fprintf(&s.b, format, args...)
}

func (s *stringBuilder) String() string { return s.b.String() }
