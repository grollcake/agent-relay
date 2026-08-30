package memento

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func newTestUpstream(t *testing.T, app *App) string {
	t.Helper()
	root := repositoryRoot(t)
	upstream := t.TempDir()
	mustCopy(t, filepath.Join(root, "VERSION"), filepath.Join(upstream, "VERSION"), 0o644)
	mustCopy(t, filepath.Join(root, "bootstrap", "AGENTS.md"), filepath.Join(upstream, "bootstrap", "AGENTS.md"), 0o644)
	mustCopy(t, filepath.Join(root, "bootstrap", "CLAUDE.md"), filepath.Join(upstream, "bootstrap", "CLAUDE.md"), 0o644)
	upstreamMemento := filepath.Join(upstream, "bootstrap", ".memento")
	for _, name := range managedMementoFiles {
		mustCopy(t, filepath.Join(root, "bootstrap", ".memento", name), filepath.Join(upstreamMemento, name), 0o644)
	}
	if err := copyTree(filepath.Join(root, "bootstrap", ".memento", "templates"), filepath.Join(upstreamMemento, "templates"), nil); err != nil {
		t.Fatal(err)
	}
	binary := filepath.Join(upstreamMemento, "bin", app.GOOS+"-"+app.GOARCH, binaryName(app.GOOS))
	if err := os.MkdirAll(filepath.Dir(binary), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(binary, []byte("new test binary"), 0o755); err != nil {
		t.Fatal(err)
	}
	digest, err := fileSHA256(binary)
	if err != nil {
		t.Fatal(err)
	}
	checksum := digest + "  " + filepath.ToSlash(app.upstreamBinaryRelative()) + "\n"
	if err := os.WriteFile(filepath.Join(upstreamMemento, "bin", "SHA256SUMS"), []byte(checksum), 0o644); err != nil {
		t.Fatal(err)
	}
	return upstream
}

func TestUpdatePreservesProjectState(t *testing.T) {
	harness := newHarness(t)
	upstream := newTestUpstream(t, harness.app)
	if err := os.Remove(harness.app.installedBinaryPath()); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(harness.app.mementoPath("bin", "SHA256SUMS")); err != nil {
		t.Fatal(err)
	}
	state := map[string][]byte{
		"GUIDANCE.md":              []byte("custom guidance\n"),
		"LESSON-LEARNED.md":        []byte("custom lesson index\n"),
		"lesson-learned/custom.md": []byte("custom lesson\n"),
		"runs/preserved.md":        []byte("custom run\n"),
	}
	for path, content := range state {
		fullPath := harness.app.mementoPath(strings.Split(path, "/")...)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(fullPath, content, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(harness.app.mementoPath("VERSION"), []byte("0.9.0\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(harness.app.mementoPath("PROTOCOL.md"), []byte("stale managed protocol\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(harness.app.mementoPath("scripts"), 0o755); err != nil {
		t.Fatal(err)
	}
	legacy := []byte("TASK_BEGIN agent=Director task=legacy\n")
	currentLog := readFile(t, harness.app.mementoPath("memento.log"))
	if err := os.WriteFile(harness.app.mementoPath("memento.log"), append(legacy, currentLog...), 0o644); err != nil {
		t.Fatal(err)
	}
	beforeLog := readFile(t, harness.app.mementoPath("memento.log"))
	beforeLines := bytes.Count(beforeLog, []byte("\n"))

	harness.run(t, "update", "--upstream", upstream, "--apply")
	for path, expected := range state {
		actual := readFile(t, harness.app.mementoPath(strings.Split(path, "/")...))
		if !bytes.Equal(actual, expected) {
			t.Fatalf("%s changed during update", path)
		}
	}
	afterLog := readFile(t, harness.app.mementoPath("memento.log"))
	if !bytes.HasPrefix(afterLog, beforeLog) {
		t.Fatal("existing memento.log lines changed")
	}
	if lines := bytes.Count(afterLog, []byte("\n")); lines != beforeLines+2 {
		t.Fatalf("update appended %d lines, want 2", lines-beforeLines)
	}
	expectedProtocol := readFile(t, filepath.Join(upstream, "bootstrap", ".memento", "PROTOCOL.md"))
	if actual := readFile(t, harness.app.mementoPath("PROTOCOL.md")); !bytes.Equal(actual, expectedProtocol) {
		t.Fatal("managed protocol was not updated")
	}
	if _, err := os.Stat(harness.app.mementoPath("scripts")); !os.IsNotExist(err) {
		t.Fatal("legacy scripts directory was not removed")
	}
	if actual := string(readFile(t, harness.app.installedBinaryPath())); actual != "new test binary" {
		t.Fatalf("binary was not updated: %q", actual)
	}
	harness.run(t, "lint")
}

func TestUpdateDryRunDoesNotMutate(t *testing.T) {
	harness := newHarness(t)
	upstream := newTestUpstream(t, harness.app)
	before := readFile(t, harness.app.mementoPath("memento.log"))
	output := harness.run(t, "update", "--upstream", upstream)
	if !strings.Contains(output, "Dry run") {
		t.Fatalf("unexpected dry-run output: %s", output)
	}
	if after := readFile(t, harness.app.mementoPath("memento.log")); !bytes.Equal(before, after) {
		t.Fatal("dry-run changed memento.log")
	}
}

func TestUpdateRejectsTamperedUpstreamBinary(t *testing.T) {
	harness := newHarness(t)
	upstream := newTestUpstream(t, harness.app)
	if err := os.WriteFile(harness.app.upstreamBinary(upstream), []byte("tampered"), 0o755); err != nil {
		t.Fatal(err)
	}
	before := readFile(t, harness.app.mementoPath("VERSION"))
	if err := harness.fail("update", "--upstream", upstream, "--apply"); err == nil {
		t.Fatal("update accepted a binary checksum mismatch")
	}
	if after := readFile(t, harness.app.mementoPath("VERSION")); !bytes.Equal(before, after) {
		t.Fatal("failed update changed VERSION")
	}
}
