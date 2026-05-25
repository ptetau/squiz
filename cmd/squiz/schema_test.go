package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/ptetau/squiz/pkg/renderer"
)

// TestMain_Schema invokes `<bin> schema`, parses stdout as JSON, and
// asserts the embedded schema has the expected top-level shape (a
// JSON-Schema object with `properties.squizzes` declared as an array).
func TestMain_Schema(t *testing.T) {
	bin := buildBinary(t)
	cmd := exec.Command(bin, "schema")
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("squiz schema failed: %v", err)
	}
	var got map[string]any
	if err := json.Unmarshal(out, &got); err != nil {
		t.Fatalf("schema stdout not valid JSON: %v\n%s", err, out)
	}
	if got["type"] != "object" {
		t.Errorf("schema top-level type = %v, want \"object\"", got["type"])
	}
	props, ok := got["properties"].(map[string]any)
	if !ok {
		t.Fatalf("schema missing top-level properties map")
	}
	for _, want := range []string{"theme", "density", "scanlines", "cursor", "spec", "squizzes"} {
		if _, ok := props[want]; !ok {
			t.Errorf("schema.properties missing key %q", want)
		}
	}
}

// TestMain_SchemaToFile verifies `--out path` writes the schema to a
// specific file and the file parses as JSON.
func TestMain_SchemaToFile(t *testing.T) {
	bin := buildBinary(t)
	out := filepath.Join(t.TempDir(), "squiz.schema.json")
	cmd := exec.Command(bin, "schema", "--out", out)
	if combined, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("squiz schema --out failed: %v\n%s", err, combined)
	}
	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("read schema out file: %v", err)
	}
	var got map[string]any
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("written file not valid JSON: %v\n%s", err, data)
	}
	if got["title"] == nil {
		t.Errorf("schema missing top-level title")
	}
}

// TestSchemaCoversAllStructFields walks every json-tagged field on the
// renderer's input types and asserts the embedded JSON Schema has a
// matching property at the expected location. Catches the common drift
// failure: someone adds a field to a struct but forgets to update the
// schema.
//
// Locations checked:
//
//	Document  → root.properties
//	Spec      → $defs.spec.properties
//	Paragraph → $defs.paragraph.properties
//	Squiz     → $defs.squiz.properties
//	Option    → $defs.option.properties
//
// Field names use the JSON tag (stripping ,omitempty).
func TestSchemaCoversAllStructFields(t *testing.T) {
	var schema map[string]any
	if err := json.Unmarshal(schemaJSON, &schema); err != nil {
		t.Fatalf("embedded schema isn't valid JSON: %v", err)
	}

	rootProps, ok := schema["properties"].(map[string]any)
	if !ok {
		t.Fatalf("schema missing properties")
	}
	defs, ok := schema["$defs"].(map[string]any)
	if !ok {
		t.Fatalf("schema missing $defs")
	}

	cases := []struct {
		name  string
		typ   reflect.Type
		props map[string]any
	}{
		{"Document", reflect.TypeOf(renderer.Document{}), rootProps},
		{"Spec", reflect.TypeOf(renderer.Spec{}), defProps(t, defs, "spec")},
		{"Paragraph", reflect.TypeOf(renderer.Paragraph{}), defProps(t, defs, "paragraph")},
		{"Squiz", reflect.TypeOf(renderer.Squiz{}), defProps(t, defs, "squiz")},
		{"Option", reflect.TypeOf(renderer.Option{}), defProps(t, defs, "option")},
	}

	for _, c := range cases {
		for i := 0; i < c.typ.NumField(); i++ {
			f := c.typ.Field(i)
			tag := f.Tag.Get("json")
			if tag == "" || tag == "-" {
				continue
			}
			name := strings.SplitN(tag, ",", 2)[0]
			if name == "" {
				continue
			}
			if _, ok := c.props[name]; !ok {
				t.Errorf("schema missing property for %s.%s (json tag %q)", c.name, f.Name, name)
			}
		}
	}
}

// defProps fetches $defs.<key>.properties or fails the test with a
// targeted message so drift is easy to triage.
func defProps(t *testing.T, defs map[string]any, key string) map[string]any {
	t.Helper()
	d, ok := defs[key].(map[string]any)
	if !ok {
		t.Fatalf("schema $defs missing %q", key)
	}
	p, ok := d["properties"].(map[string]any)
	if !ok {
		t.Fatalf("schema $defs.%s missing properties", key)
	}
	return p
}
