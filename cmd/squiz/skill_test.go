package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// TestEmbeddedSkillMatchesRoot asserts the literal copy of SKILL.md kept
// next to the binary for //go:embed (cmd/squiz/skill.md) is byte-identical
// to the canonical source at skills/squiz/SKILL.md.
//
// The two paths exist because //go:embed cannot reach above the package
// directory; this test is the drift detector. If it fails, sync the copy:
//
//	cp skills/squiz/SKILL.md cmd/squiz/skill.md
func TestEmbeddedSkillMatchesRoot(t *testing.T) {
	root := repoRoot(t)
	canonical := filepath.Join(root, "skills", "squiz", "SKILL.md")
	onDisk, err := os.ReadFile(canonical)
	if err != nil {
		t.Fatalf("read canonical SKILL.md at %s: %v", canonical, err)
	}
	if !bytes.Equal(embeddedSkill, onDisk) {
		t.Fatalf("embedded SKILL.md drifted from canonical source\n"+
			"  canonical: %s (%d bytes)\n"+
			"  embedded:  cmd/squiz/skill.md (%d bytes)\n"+
			"fix: cp skills/squiz/SKILL.md cmd/squiz/skill.md",
			canonical, len(onDisk), len(embeddedSkill))
	}
}

// TestMain_Skill execs `squiz skill` and asserts stdout matches the
// embedded bytes exactly. End-to-end confirmation that the dump path is
// wired and byte-clean (no surprise BOM, no trailing newline injection,
// no encoding conversion).
func TestMain_Skill(t *testing.T) {
	bin := buildBinary(t)

	cmd := exec.Command(bin, "skill")
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("squiz skill failed: %v\nstdout: %s", err, out)
	}
	if exit := cmd.ProcessState.ExitCode(); exit != 0 {
		t.Fatalf("squiz skill exit = %d, want 0", exit)
	}
	if !bytes.Equal(out, embeddedSkill) {
		t.Fatalf("squiz skill stdout (%d bytes) != embeddedSkill (%d bytes)",
			len(out), len(embeddedSkill))
	}
}
