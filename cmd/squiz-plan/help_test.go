package main

import (
	"os/exec"
	"strings"
	"testing"
)

// TestMain_Help_ListsTopics execs `squiz-plan help` (no topic) and asserts
// stdout includes the headline topic names.
func TestMain_Help_ListsTopics(t *testing.T) {
	bin := buildBinary(t)

	cmd := exec.Command(bin, "help")
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("squiz-plan help failed: %v\nstdout: %s", err, out)
	}
	if exit := cmd.ProcessState.ExitCode(); exit != 0 {
		t.Fatalf("squiz-plan help exit = %d, want 0", exit)
	}
	s := string(out)
	for _, want := range []string{"art", "themes"} {
		if !strings.Contains(s, want) {
			t.Errorf("stdout missing topic %q\nfull stdout:\n%s", want, s)
		}
	}
}

// TestMain_Help_KnownTopic execs `squiz-plan help art` and asserts the
// topic body is substantial.
func TestMain_Help_KnownTopic(t *testing.T) {
	bin := buildBinary(t)

	cmd := exec.Command(bin, "help", "art")
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("squiz-plan help art failed: %v\nstdout: %s", err, out)
	}
	if exit := cmd.ProcessState.ExitCode(); exit != 0 {
		t.Fatalf("squiz-plan help art exit = %d, want 0", exit)
	}
	if len(out) < 500 {
		t.Errorf("squiz-plan help art body too short: %d bytes (want > 500)\nbody: %s",
			len(out), out)
	}
}

// TestMain_Help_UnknownTopic execs `squiz-plan help nosuch` and asserts:
//   - exit code 1
//   - stderr lists the available topics
func TestMain_Help_UnknownTopic(t *testing.T) {
	bin := buildBinary(t)

	cmd := exec.Command(bin, "help", "nosuch")
	var stderrBuf strings.Builder
	cmd.Stderr = &stderrBuf
	stdout, _ := cmd.Output()

	if exit := cmd.ProcessState.ExitCode(); exit != 1 {
		t.Fatalf("squiz-plan help nosuch exit = %d, want 1\nstdout: %s\nstderr: %s",
			exit, stdout, stderrBuf.String())
	}
	stderr := stderrBuf.String()
	if !strings.Contains(stderr, "art") || !strings.Contains(stderr, "themes") {
		t.Errorf("stderr should list available topics; got:\n%s", stderr)
	}
}
