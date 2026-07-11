package relay

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var managedRelayFiles = []string{
	"HOW-TO-UPDATE.md",
	"PROTOCOL.md",
	"DIRECTOR.md",
	"PLANNER.md",
	"EXECUTOR.md",
}

func (a *App) upstreamBinary(upstream string) string {
	return filepath.Join(upstream, "bootstrap", ".agent-relay", "bin", a.GOOS+"-"+a.GOARCH, binaryName(a.GOOS))
}

func (a *App) upstreamBinaryRelative() string {
	return filepath.Join(a.GOOS+"-"+a.GOARCH, binaryName(a.GOOS))
}

func samePath(left, right string) bool {
	leftAbs, leftErr := filepath.Abs(left)
	rightAbs, rightErr := filepath.Abs(right)
	if leftErr != nil || rightErr != nil {
		return false
	}
	return strings.EqualFold(filepath.Clean(leftAbs), filepath.Clean(rightAbs))
}

func (a *App) preflightUpdate(upstream string) error {
	upstreamRelay := filepath.Join(upstream, "bootstrap", ".agent-relay")
	checksumFile := filepath.Join(upstreamRelay, "bin", "SHA256SUMS")
	required := []string{filepath.Join(upstream, "VERSION"), filepath.Join(upstream, "bootstrap", "AGENTS.md"), a.upstreamBinary(upstream), checksumFile}
	for _, name := range managedRelayFiles {
		required = append(required, filepath.Join(upstreamRelay, name))
	}
	for _, path := range required {
		if info, err := os.Stat(path); err != nil || info.IsDir() {
			return fmt.Errorf("missing upstream file: %s", path)
		}
	}
	if err := verifyChecksum(a.upstreamBinary(upstream), checksumFile, a.upstreamBinaryRelative()); err != nil {
		return fmt.Errorf("invalid upstream binary: %w", err)
	}
	if runtime.GOOS == "windows" {
		if executable, err := os.Executable(); err == nil && samePath(executable, a.installedBinaryPath()) {
			return errors.New("Windows update must be run with the new upstream agent-relay.exe from a temporary path")
		}
	}
	return nil
}

func (a *App) runUpdate(args []string) error {
	parsed, err := parseArguments(args, map[string]bool{"--upstream": true}, map[string]bool{"--apply": true})
	if err != nil {
		return err
	}
	if len(parsed.pos) != 0 {
		return errors.New("update accepts flags only")
	}
	upstream, err := requireValue(parsed, "--upstream")
	if err != nil {
		return err
	}
	upstream, err = filepath.Abs(upstream)
	if err != nil {
		return err
	}
	if info, statErr := os.Stat(filepath.Join(upstream, "bootstrap", ".agent-relay")); statErr != nil || !info.IsDir() {
		return fmt.Errorf("invalid upstream: %s", upstream)
	}
	currentVersionBytes, err := os.ReadFile(a.relayPath("VERSION"))
	if err != nil {
		return fmt.Errorf("missing %s", a.relayPath("VERSION"))
	}
	nextVersionBytes, err := os.ReadFile(filepath.Join(upstream, "VERSION"))
	if err != nil {
		return errors.New("missing upstream VERSION")
	}
	currentVersion := strings.TrimSpace(string(currentVersionBytes))
	nextVersion := strings.TrimSpace(string(nextVersionBytes))
	fmt.Fprintf(a.Stdout, "Agent Relay update: %s -> %s\n", currentVersion, nextVersion)

	if !parsed.flags["--apply"] {
		fmt.Fprint(a.Stdout, `Dry run. Re-run with --apply to update:
- WARNING: --apply replaces the managed files below; review local customizations first.
- AGENTS.md Agent Relay block
- CLAUDE.md Agent Relay block when present
- .agent-relay/HOW-TO-UPDATE.md
- .agent-relay/PROTOCOL.md
- .agent-relay/DIRECTOR.md
- .agent-relay/PLANNER.md
- .agent-relay/EXECUTOR.md
- .agent-relay/bin/agent-relay
- .agent-relay/templates/
- .agent-relay/VERSION

Preserved:
- .agent-relay/GUIDANCE.md
- .agent-relay/LESSON-LEARNED.md
- .agent-relay/lesson-learned/
- .agent-relay/relay.log
- .agent-relay/runs/
`)
		return nil
	}
	if err := a.preflightUpdate(upstream); err != nil {
		return err
	}

	upstreamRelay := filepath.Join(upstream, "bootstrap", ".agent-relay")
	agentsPath := filepath.Join(a.ProjectDir, "AGENTS.md")
	if _, err := os.Stat(agentsPath); os.IsNotExist(err) {
		if err := copyFile(filepath.Join(upstream, "bootstrap", "AGENTS.md"), agentsPath, 0o644); err != nil {
			return err
		}
	} else if err := MergeAgentBlock(agentsPath, filepath.Join(upstream, "bootstrap", "AGENTS.md")); err != nil {
		return err
	}
	claudePath := filepath.Join(a.ProjectDir, "CLAUDE.md")
	if _, err := os.Stat(claudePath); err == nil {
		if err := MergeAgentBlock(claudePath, filepath.Join(upstream, "bootstrap", "CLAUDE.md")); err != nil {
			return err
		}
	}
	for _, name := range managedRelayFiles {
		if err := copyFile(filepath.Join(upstreamRelay, name), a.relayPath(name), 0o644); err != nil {
			return err
		}
	}
	if err := os.MkdirAll(a.relayPath("templates"), 0o755); err != nil {
		return err
	}
	templates, err := filepath.Glob(filepath.Join(upstreamRelay, "templates", "*.md"))
	if err != nil {
		return err
	}
	for _, source := range templates {
		if err := copyFile(source, a.relayPath("templates", filepath.Base(source)), 0o644); err != nil {
			return err
		}
	}
	if err := copyFile(a.upstreamBinary(upstream), a.installedBinaryPath(), 0o755); err != nil {
		return err
	}
	if err := copyFile(filepath.Join(upstreamRelay, "bin", "SHA256SUMS"), a.relayPath("bin", "SHA256SUMS"), 0o644); err != nil {
		return err
	}
	if err := copyFile(filepath.Join(upstream, "VERSION"), a.relayPath("VERSION"), 0o644); err != nil {
		return err
	}
	if err := os.RemoveAll(a.relayPath("scripts")); err != nil {
		return err
	}

	records, err := a.readRecords()
	if err != nil {
		return err
	}
	taskID, err := a.generateTaskID(records)
	if err != nil {
		return err
	}
	now := a.Now().Format("2006-01-02T15:04:05")
	if _, err := a.appendRecord(Record{Timestamp: now, TaskID: taskID, Event: eventRequest, Role: "Director", Summary: fmt.Sprintf("Update Agent Relay %s -> %s", currentVersion, nextVersion)}); err != nil {
		return err
	}
	if _, err := a.appendRecord(Record{Timestamp: now, TaskID: taskID, Event: eventRunDone, Role: "Director", Summary: fmt.Sprintf("Agent Relay updated to %s", nextVersion)}); err != nil {
		return err
	}
	fmt.Fprintf(a.Stdout, "Agent Relay updated: %s -> %s\n", currentVersion, nextVersion)
	return nil
}

func (a *App) runVersion(args []string) error {
	if len(args) != 0 {
		return errors.New("version accepts no arguments")
	}
	version, err := os.ReadFile(a.relayPath("VERSION"))
	if err != nil {
		return err
	}
	fmt.Fprintln(a.Stdout, strings.TrimSpace(string(version)))
	return nil
}
