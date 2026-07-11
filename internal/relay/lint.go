package relay

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type lintState struct {
	app    *App
	errors int
}

func (state *lintState) ok(format string, args ...any) {
	fmt.Fprintf(state.app.Stdout, "OK: "+format+"\n", args...)
}

func (state *lintState) err(format string, args ...any) {
	state.errors++
	fmt.Fprintf(state.app.Stderr, "ERROR: "+format+"\n", args...)
}

func (state *lintState) requireFile(path string) {
	if info, err := os.Stat(path); err == nil && !info.IsDir() {
		state.ok("%s exists", path)
	} else {
		state.err("missing %s", path)
	}
}

func (state *lintState) requireDir(path string) {
	if info, err := os.Stat(path); err == nil && info.IsDir() {
		state.ok("%s exists", path)
	} else {
		state.err("missing %s", path)
	}
}

var agentBlockPattern = regexp.MustCompile(`(?s)<agent-relay-rules>.*?</agent-relay-rules>`)

func (state *lintState) checkAgentBlocks() {
	agentsPath := filepath.Join(state.app.ProjectDir, "AGENTS.md")
	claudePath := filepath.Join(state.app.ProjectDir, "CLAUDE.md")
	agents, agentsErr := os.ReadFile(agentsPath)
	claude, claudeErr := os.ReadFile(claudePath)
	if agentsErr != nil || claudeErr != nil {
		return
	}
	agentsBlock := agentBlockPattern.Find(agents)
	claudeBlock := agentBlockPattern.Find(claude)
	if len(agentsBlock) == 0 {
		state.err("AGENTS.md missing <agent-relay-rules> block")
	} else if len(claudeBlock) == 0 {
		state.err("CLAUDE.md missing <agent-relay-rules> block")
	} else if string(agentsBlock) != string(claudeBlock) {
		state.err("AGENTS.md and CLAUDE.md Agent Relay blocks differ")
	} else {
		state.ok("AGENTS.md and CLAUDE.md Agent Relay blocks match")
	}
}

func (state *lintState) checkLog() {
	logPath := state.app.relayPath("relay.log")
	file, err := os.Open(logPath)
	if err != nil {
		return
	}
	defer file.Close()

	lastEvents := map[string]string{}
	taskKeys := map[string]string{}
	pendingRounds := map[string]string{}
	pendingKeys := map[string]string{}
	legacyLines := 0

	scanner := bufio.NewScanner(file)
	for lineNumber := 1; scanner.Scan(); lineNumber++ {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		record, parseErr := parseRecord(line, lineNumber)
		if parseErr != nil {
			if isLegacyRecordLine(line) {
				legacyLines++
				continue
			}
			state.err("%s", parseErr)
			continue
		}
		if record.TaskID == "" {
			state.err("missing task-id on line %d", lineNumber)
		}
		expectedRole, validEvent := expectedRoles[record.Event]
		if !validEvent {
			state.err("invalid event on line %d: %s", lineNumber, record.Event)
		} else if record.Role != expectedRole {
			state.err("invalid role on line %d: %s requires %s, got %s", lineNumber, record.Event, expectedRole, record.Role)
		}
		if record.TaskID != "" && validEvent {
			prior := lastEvents[record.TaskID]
			if !validTransition(prior, record.Event) {
				shownPrior := valueOr(prior, "START")
				state.err("invalid transition on line %d for task-id %s: %s -> %s", lineNumber, record.TaskID, shownPrior, record.Event)
			}
			lastEvents[record.TaskID] = record.Event
		}

		requiresArtifact := record.Event == eventPlanned || record.Event == eventExecuted || record.Event == eventReview || record.Event == eventClose
		if requiresArtifact && record.Path == "" {
			state.err("missing path on line %d: %s", lineNumber, record.Event)
		}
		if record.Path != "" {
			if requiresArtifact {
				if artifactErr := state.app.CheckArtifact(record.Event, record.Path, record.TaskID); artifactErr != nil {
					state.err("relay log line %d: %s", lineNumber, artifactErr)
				}
			} else if info, statErr := os.Stat(state.app.projectPath(record.Path)); statErr != nil || info.IsDir() {
				state.err("artifact path not found on line %d: %s", lineNumber, record.Path)
			}
		}

		switch record.Event {
		case eventPlanned:
			if record.Path != "" {
				taskKeys[record.TaskID] = artifactKey(record.Path, record.Event)
			}
		case eventExecuted:
			if record.Path == "" {
				continue
			}
			round := artifactRound(record.Path, record.Event)
			key := artifactKey(record.Path, record.Event)
			if key != taskKeys[record.TaskID] {
				state.err("RUN artifact key does not match PLAN on line %d for task-id %s", lineNumber, record.TaskID)
			}
			if round == "" {
				state.err("invalid RUN artifact name on line %d: %s", lineNumber, record.Path)
			} else {
				pendingRounds[record.TaskID] = round
				pendingKeys[record.TaskID] = key
			}
		case eventReview:
			if record.Path == "" {
				continue
			}
			round := artifactRound(record.Path, record.Event)
			key := artifactKey(record.Path, record.Event)
			pendingRound, hasRun := pendingRounds[record.TaskID]
			if round == "" {
				state.err("invalid REVIEW artifact name on line %d: %s", lineNumber, record.Path)
			} else if !hasRun {
				state.err("REVIEW has no corresponding EXECUTED on line %d for task-id %s", lineNumber, record.TaskID)
			} else {
				if round != pendingRound {
					state.err("round mismatch on line %d for task-id %s: RUN-%s -> REVIEW-%s", lineNumber, record.TaskID, pendingRound, round)
				}
				if key != pendingKeys[record.TaskID] {
					state.err("artifact key mismatch on line %d for task-id %s", lineNumber, record.TaskID)
				}
				if key != taskKeys[record.TaskID] {
					state.err("REVIEW artifact key does not match PLAN on line %d for task-id %s", lineNumber, record.TaskID)
				}
				delete(pendingRounds, record.TaskID)
				delete(pendingKeys, record.TaskID)
			}
		case eventClose:
			if record.Path != "" && artifactKey(record.Path, record.Event) != taskKeys[record.TaskID] {
				state.err("CLOSE artifact key does not match PLAN on line %d for task-id %s", lineNumber, record.TaskID)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		state.err("cannot read relay.log: %s", err)
	}
	if legacyLines > 0 {
		state.ok("legacy relay log lines preserved: %d", legacyLines)
	}

	closedCount := 0
	for _, event := range lastEvents {
		if event == eventClose || event == eventRunDone {
			closedCount++
		}
	}
	state.ok("relay tasks: %d total, %d closed, %d open", len(lastEvents), closedCount, len(lastEvents)-closedCount)
}

func (a *App) Lint() error {
	state := &lintState{app: a}
	for _, name := range []string{"PROTOCOL.md", "DIRECTOR.md", "PLANNER.md", "EXECUTOR.md", "HOW-TO-UPDATE.md", "VERSION", "GUIDANCE.md", "LESSON-LEARNED.md", "relay.log"} {
		state.requireFile(a.relayPath(name))
	}
	state.requireFile(a.installedBinaryPath())
	state.requireFile(a.relayPath("bin", "SHA256SUMS"))
	for _, name := range []string{"templates", "runs", "lesson-learned", "bin"} {
		state.requireDir(a.relayPath(name))
	}
	if a.GOOS != "windows" {
		if info, err := os.Stat(a.installedBinaryPath()); err == nil && info.Mode()&0o111 != 0 {
			state.ok("agent-relay binary is executable")
		} else {
			state.err("agent-relay binary is not executable")
		}
	}
	if err := verifyChecksum(a.installedBinaryPath(), a.relayPath("bin", "SHA256SUMS"), filepath.Join(a.GOOS+"-"+a.GOARCH, binaryName(a.GOOS))); err != nil {
		state.err("agent-relay binary checksum failed: %s", err)
	} else {
		state.ok("agent-relay binary checksum matches")
	}
	if _, err := os.Stat(a.relayPath("scripts")); os.IsNotExist(err) {
		state.ok("legacy scripts directory absent")
	} else {
		state.err("legacy scripts directory must be removed")
	}
	if _, err := os.Stat(a.relayPath("protocol-guard")); os.IsNotExist(err) {
		state.ok("legacy protocol-guard absent")
	} else {
		state.err("legacy protocol-guard must be removed")
	}
	state.checkAgentBlocks()
	state.checkLog()
	if state.errors > 0 {
		return fmt.Errorf("relay-lint failed: %d error(s)", state.errors)
	}
	fmt.Fprintln(a.Stdout, "relay-lint passed")
	return nil
}
