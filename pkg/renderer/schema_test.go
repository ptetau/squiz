package renderer

import "testing"

func TestParseDocument_Defaults(t *testing.T) {
	data := []byte(`{"squizzes": []}`)
	d, err := ParseDocument(data)
	if err != nil {
		t.Fatalf("ParseDocument returned error: %v", err)
	}
	if d == nil {
		t.Fatal("ParseDocument returned nil document")
	}
	if d.Density != "compact" {
		t.Errorf("Density = %q, want %q", d.Density, "compact")
	}
	if d.Cursor == nil {
		t.Fatal("Cursor = nil, want non-nil pointer to true")
	}
	if *d.Cursor != true {
		t.Errorf("*Cursor = %v, want true", *d.Cursor)
	}
	// Theme stays empty intentionally so auto-rotation can take over.
	if d.Theme != "" {
		t.Errorf("Theme = %q, want empty string (auto-rotation depends on this)", d.Theme)
	}
}

func TestParseDocument_PreservesExplicit(t *testing.T) {
	data := []byte(`{
		"density": "comfortable",
		"cursor": false,
		"theme": "phosphor",
		"squizzes": []
	}`)
	d, err := ParseDocument(data)
	if err != nil {
		t.Fatalf("ParseDocument returned error: %v", err)
	}
	if d.Density != "comfortable" {
		t.Errorf("Density = %q, want %q", d.Density, "comfortable")
	}
	if d.Cursor == nil {
		t.Fatal("Cursor = nil, want non-nil pointer to false")
	}
	if *d.Cursor != false {
		t.Errorf("*Cursor = %v, want false", *d.Cursor)
	}
	if d.Theme != "phosphor" {
		t.Errorf("Theme = %q, want %q", d.Theme, "phosphor")
	}
}

func TestParseDocument_Malformed(t *testing.T) {
	data := []byte(`{invalid`)
	d, err := ParseDocument(data)
	if err == nil {
		t.Fatalf("ParseDocument returned no error for malformed JSON; doc=%+v", d)
	}
	if d != nil {
		t.Errorf("ParseDocument returned non-nil document for malformed JSON: %+v", d)
	}
}

func TestOption_ResolvedArt(t *testing.T) {
	cases := []struct {
		name string
		opt  Option
		want string
	}{
		{
			name: "Art set, ArtSVG empty",
			opt:  Option{Art: "wf:calendar-grid"},
			want: "wf:calendar-grid",
		},
		{
			name: "ArtSVG set, Art empty",
			opt:  Option{ArtSVG: "<svg></svg>"},
			want: "<svg></svg>",
		},
		{
			name: "both set, Art preferred",
			opt:  Option{Art: "grid:7x7@0.55", ArtSVG: "<svg></svg>"},
			want: "grid:7x7@0.55",
		},
		{
			name: "both empty",
			opt:  Option{},
			want: "",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.opt.ResolvedArt()
			if got != tc.want {
				t.Errorf("ResolvedArt() = %q, want %q", got, tc.want)
			}
		})
	}
}
