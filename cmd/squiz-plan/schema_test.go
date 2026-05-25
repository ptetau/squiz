package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/ptetau/squiz/internal/planview"
)

// TestMain_Schema invokes `<bin> schema`, parses stdout as JSON, and
// asserts the embedded schema has the expected top-level shape (a
// JSON-Schema object with `properties.sections` declared).
func TestMain_Schema(t *testing.T) {
	bin := buildBinary(t)
	cmd := exec.Command(bin, "schema")
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("squiz-plan schema failed: %v", err)
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
	for _, want := range []string{"title", "lede", "theme", "density", "sections"} {
		if _, ok := props[want]; !ok {
			t.Errorf("schema.properties missing key %q", want)
		}
	}
}

// TestMain_SchemaToFile verifies --out writes to a path and the file is
// well-formed JSON.
func TestMain_SchemaToFile(t *testing.T) {
	bin := buildBinary(t)
	out := filepath.Join(t.TempDir(), "squiz-plan.schema.json")
	cmd := exec.Command(bin, "schema", "--out", out)
	if combined, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("squiz-plan schema --out failed: %v\n%s", err, combined)
	}
	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("read schema out file: %v", err)
	}
	var got map[string]any
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("written file not valid JSON: %v\n%s", err, data)
	}
}

// TestSchemaCoversAllStructFields walks every json-tagged field on
// planview's input types and asserts the embedded JSON Schema has a
// matching property at the expected location. The schema covers
// index.json directly; item/option live under $defs since they're only
// described for reference (the parser owns their validation).
//
// Locations checked:
//
//	Index       → root.properties
//	SectionFile → $defs.sectionFile.properties
//	Item        → $defs.item.properties
//	Option      → $defs.option.properties
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
		{"Index", reflect.TypeOf(planview.Index{}), rootProps},
		{"SectionFile", reflect.TypeOf(planview.SectionFile{}), defProps(t, defs, "sectionFile")},
		{"Item", reflect.TypeOf(planview.Item{}), defProps(t, defs, "item")},
		{"Option", reflect.TypeOf(planview.Option{}), defProps(t, defs, "option")},
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
