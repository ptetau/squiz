package main

import (
	"os/exec"
	"strings"
	"testing"
)

// TestMain_Help_ListsTopics execs `squiz help` (no topic) and asserts
// stdout includes the headline topic names. The topic-list surface goes to
// stdout (not stderr) so it can be piped/grep'd.
func TestMain_Help_ListsTopics(t *testing.T) {
	bin := buildBinary(t)

	cmd := exec.Command(bin, "help")
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("squiz help failed: %v\nstdout: %s", err, out)
	}
	if exit := cmd.ProcessState.ExitCode(); exit != 0 {
		t.Fatalf("squiz help exit = %d, want 0", exit)
	}
	s := string(out)
	for _, want := range []string{"art", "themes"} {
		if !strings.Contains(s, want) {
			t.Errorf("stdout missing topic %q\nfull stdout:\n%s", want, s)
		}
	}
}

// TestMain_Help_KnownTopic execs `squiz help art` and asserts the topic
// body is substantial. Topic refs are documented as stand-alone — anything
// shorter than ~500 bytes is suspect.
func TestMain_Help_KnownTopic(t *testing.T) {
	bin := buildBinary(t)

	cmd := exec.Command(bin, "help", "art")
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("squiz help art failed: %v\nstdout: %s", err, out)
	}
	if exit := cmd.ProcessState.ExitCode(); exit != 0 {
		t.Fatalf("squiz help art exit = %d, want 0", exit)
	}
	if len(out) < 500 {
		t.Errorf("squiz help art body too short: %d bytes (want > 500)\nbody: %s",
			len(out), out)
	}
}

// TestMain_Help_UnknownTopic execs `squiz help nosuch` and asserts:
//   - exit code 1
//   - stderr lists the available topics (so the user can recover)
func TestMain_Help_UnknownTopic(t *testing.T) {
	bin := buildBinary(t)

	cmd := exec.Command(bin, "help", "nosuch")
	// We expect non-zero exit, so Output() will return an *exec.ExitError;
	// capture stderr separately rather than failing on err.
	var stderrBuf strings.Builder
	cmd.Stderr = &stderrBuf
	stdout, _ := cmd.Output()

	if exit := cmd.ProcessState.ExitCode(); exit != 1 {
		t.Fatalf("squiz help nosuch exit = %d, want 1\nstdout: %s\nstderr: %s",
			exit, stdout, stderrBuf.String())
	}
	stderr := stderrBuf.String()
	if !strings.Contains(stderr, "art") || !strings.Contains(stderr, "themes") {
		t.Errorf("stderr should list available topics; got:\n%s", stderr)
	}
}
