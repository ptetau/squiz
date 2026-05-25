package main

import (
	"reflect"
	"strings"
	"testing"
)

// boolFlags mirrors what cmdRender passes — only --stdout and --open are
// valueless. Everything else (--out, --theme) consumes the next arg.
var testBoolFlags = map[string]bool{"stdout": true, "open": true}

func TestReorderFlagsFirst(t *testing.T) {
	cases := []struct {
		name string
		in   []string
		want []string
	}{
		{
			name: "no flags",
			in:   []string{"in.json"},
			want: []string{"in.json"},
		},
		{
			name: "flags already first",
			in:   []string{"--theme", "paper", "--out", "x.html", "in.json"},
			want: []string{"--theme", "paper", "--out", "x.html", "in.json"},
		},
		{
			name: "value-flag after positional gets reordered with its value",
			in:   []string{"in.json", "--out", "x.html"},
			want: []string{"--out", "x.html", "in.json"},
		},
		{
			name: "bool flag after positional",
			in:   []string{"in.json", "--open"},
			want: []string{"--open", "in.json"},
		},
		{
			name: "interleaved",
			in:   []string{"--theme", "paper", "in.json", "--out", "x.html", "--open"},
			want: []string{"--theme", "paper", "--out", "x.html", "--open", "in.json"},
		},
		{
			name: "--name=value form stays together",
			in:   []string{"in.json", "--theme=phosphor"},
			want: []string{"--theme=phosphor", "in.json"},
		},
		{
			name: "single-dash short form",
			in:   []string{"in.json", "-open"},
			want: []string{"-open", "in.json"},
		},
		{
			name: "-- terminator stops flag detection",
			in:   []string{"--theme", "paper", "--", "--not-a-flag.json"},
			want: []string{"--theme", "paper", "--", "--not-a-flag.json"},
		},
		{
			name: "bare - is positional, not flag",
			in:   []string{"-", "--theme", "paper"},
			want: []string{"--theme", "paper", "-"},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := reorderFlagsFirst(tc.in, testBoolFlags)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("reorderFlagsFirst(%v)\n  got:  %v\n  want: %v", tc.in, got, tc.want)
			}
		})
	}
}

func TestHasFlag(t *testing.T) {
	cases := []struct {
		args []string
		name string
		want bool
	}{
		{[]string{"--open"}, "open", true},
		{[]string{"-open"}, "open", true},
		{[]string{"--open=true"}, "open", true},
		{[]string{"--openother"}, "open", false},
		{[]string{}, "open", false},
		{[]string{"in.json", "--theme", "paper"}, "theme", true},
		{[]string{"in.json", "--theme", "paper"}, "open", false},
	}
	for _, tc := range cases {
		t.Run(strings.Join(tc.args, " ")+"_"+tc.name, func(t *testing.T) {
			if got := hasFlag(tc.args, tc.name); got != tc.want {
				t.Errorf("hasFlag(%v, %q) = %v, want %v", tc.args, tc.name, got, tc.want)
			}
		})
	}
}
