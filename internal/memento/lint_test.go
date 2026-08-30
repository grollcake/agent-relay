package memento

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestLintDoesNotExecutePaths(t *testing.T) {
	harness := newHarness(t)
	marker := filepath.Join(t.TempDir(), "command-ran")
	writeLog(t, harness.app.MementoDir,
		"2026-07-11T10:00:00 | boot | REQUEST  | Director | Bootstrap Memento AI",
		"2026-07-11T10:00:00 | boot | RUN_DONE | Director | Memento AI initialized",
		fmt.Sprintf("2026-07-11T10:01:00 | evil | REQUEST  | Director | Probe | $(touch %s)", marker),
		"2026-07-11T10:01:01 | evil | RUN_DONE | Director | Probe complete",
	)
	if err := harness.fail("lint"); err == nil {
		t.Fatal("lint accepted a malicious path")
	}
	if _, err := os.Stat(marker); !os.IsNotExist(err) {
		t.Fatal("lint executed a path from memento.log")
	}
}

func TestLintRejectsLegacyScriptsDirectory(t *testing.T) {
	harness := newHarness(t)
	if err := os.MkdirAll(harness.app.mementoPath("scripts"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := harness.fail("lint"); err == nil {
		t.Fatal("lint accepted a legacy scripts directory")
	}
}

func TestLintRejectsTamperedBinary(t *testing.T) {
	harness := newHarness(t)
	if err := os.WriteFile(harness.app.installedBinaryPath(), []byte("tampered"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := harness.fail("lint"); err == nil {
		t.Fatal("lint accepted a binary checksum mismatch")
	}
}
