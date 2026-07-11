package relay

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestTemplatesAreRejected(t *testing.T) {
	harness := newHarness(t)
	cases := []struct {
		event, template, path string
	}{
		{eventPlanned, "plan.md", ".agent-relay/runs/20260711-1000-probe-PLAN.md"},
		{eventExecuted, "run.md", ".agent-relay/runs/20260711-1000-probe-RUN-01.md"},
		{eventReview, "review.md", ".agent-relay/runs/20260711-1000-probe-REVIEW-01.md"},
		{eventClose, "close.md", ".agent-relay/runs/20260711-1000-probe-CLOSE.md"},
	}
	for _, testCase := range cases {
		t.Run(testCase.event, func(t *testing.T) {
			target := harness.app.projectPath(testCase.path)
			mustCopy(t, harness.app.relayPath("templates", testCase.template), target, 0o644)
			if err := harness.app.CheckArtifact(testCase.event, testCase.path, ""); err == nil {
				t.Fatalf("%s template unexpectedly passed", testCase.event)
			}
		})
	}
}

func TestTransitionMatrixAndAllowedFlows(t *testing.T) {
	harness := newHarness(t)
	events := []string{eventRequest, eventPlanned, eventExecuted, eventReview, eventFeedback, eventClose, eventRunDone}
	priors := append([]string{""}, events...)
	for _, prior := range priors {
		for _, next := range events {
			if validTransition(prior, next) {
				continue
			}
			lines := []string{
				"2026-07-11T10:00:00 | boot | REQUEST  | Director | Bootstrap Agent Relay",
				"2026-07-11T10:00:00 | boot | RUN_DONE | Director | Agent Relay initialized",
			}
			if prior != "" {
				lines = append(lines, fmt.Sprintf("2026-07-11T10:01:00 | matr | %-8s | %-8s | Matrix state", prior, expectedRoles[prior]))
			}
			writeLog(t, harness.app.RelayDir, lines...)
			err := harness.fail("append", next, "--task-id", "matr", "--role", expectedRoles[next], "--summary", "Invalid transition probe")
			if err == nil {
				t.Fatalf("invalid transition %s -> %s passed", valueOr(prior, "START"), next)
			}
		}
	}

	writeLog(t, harness.app.RelayDir,
		"2026-07-11T10:00:00 | boot | REQUEST  | Director | Bootstrap Agent Relay",
		"2026-07-11T10:00:00 | boot | RUN_DONE | Director | Agent Relay initialized",
	)
	harness.run(t, "append", eventRequest, "--task-id", "drct", "--role", "Director", "--summary", "Direct flow")
	harness.run(t, "append", eventRunDone, "--task-id", "drct", "--role", "Director", "--summary", "Direct flow complete")

	taskID, key := parseRoundOutput(t, harness.run(t, "new-round", "retry-flow", "--summary", "Retry flow"))
	planPath := ".agent-relay/runs/" + key + "-PLAN.md"
	writePlan(t, harness.app.ProjectDir, planPath, taskID)
	harness.run(t, "append", eventPlanned, "--task-id", taskID, "--role", "Planner", "--summary", "Plan complete", "--path", planPath)
	harness.run(t, "gate", "before-execute", "--task-id", taskID)

	appendRound := func(round, result string) {
		runPath := ".agent-relay/runs/" + key + "-RUN-" + round + ".md"
		reviewPath := ".agent-relay/runs/" + key + "-REVIEW-" + round + ".md"
		writeRun(t, harness.app.ProjectDir, runPath, taskID, round)
		writeReview(t, harness.app.ProjectDir, reviewPath, taskID, round, result)
		harness.run(t, "append", eventExecuted, "--task-id", taskID, "--role", "Executor", "--summary", "Run complete", "--path", runPath)
		harness.run(t, "gate", "before-review", "--task-id", taskID)
		harness.run(t, "append", eventReview, "--task-id", taskID, "--role", "Planner", "--summary", "Review complete", "--path", reviewPath)
		harness.run(t, "gate", "before-approval", "--task-id", taskID)
	}
	appendRound("01", "blockers")
	appendRound("02", "ready-for-user-decision")
	harness.run(t, "feedback", "--task-id", taskID, "--summary", "User feedback")
	appendRound("03", "ready-for-user-decision")

	status := harness.run(t, "status", "--task-id", taskID)
	if !strings.Contains(status, "last_event: REVIEW") {
		t.Fatalf("unexpected status: %s", status)
	}
	closePath := ".agent-relay/runs/" + key + "-CLOSE.md"
	writeClose(t, harness.app.ProjectDir, closePath, taskID)
	harness.run(t, "append", eventClose, "--task-id", taskID, "--role", "Director", "--summary", "Retry flow closed", "--path", closePath)
	harness.run(t, "lint")
}

func TestArtifactConsistency(t *testing.T) {
	harness := newHarness(t)
	taskID, key := parseRoundOutput(t, harness.run(t, "new-round", "consistent-flow", "--summary", "Consistency flow"))
	planPath := ".agent-relay/runs/" + key + "-PLAN.md"

	writePlan(t, harness.app.ProjectDir, planPath, taskID)
	content := strings.ReplaceAll(string(readFile(t, harness.app.projectPath(planPath))), "Date: 2026-07-11\n", "")
	writeArtifact(t, harness.app.ProjectDir, planPath, content)
	if err := harness.fail("append", eventPlanned, "--task-id", taskID, "--role", "Planner", "--summary", "Missing date", "--path", planPath); err == nil {
		t.Fatal("PLAN without a date passed")
	}

	writePlan(t, harness.app.ProjectDir, planPath, "wrong")
	if err := harness.fail("append", eventPlanned, "--task-id", taskID, "--role", "Planner", "--summary", "Wrong task", "--path", planPath); err == nil {
		t.Fatal("PLAN task-id mismatch passed")
	}
	writePlan(t, harness.app.ProjectDir, planPath, taskID)
	harness.run(t, "append", eventPlanned, "--task-id", taskID, "--role", "Planner", "--summary", "Plan complete", "--path", planPath)

	wrongRun := ".agent-relay/runs/20260711-1000-wrong-key-RUN-01.md"
	writeRun(t, harness.app.ProjectDir, wrongRun, taskID, "01")
	if err := harness.fail("append", eventExecuted, "--task-id", taskID, "--role", "Executor", "--summary", "Wrong key", "--path", wrongRun); err == nil {
		t.Fatal("RUN key mismatch passed")
	}
	runPath := ".agent-relay/runs/" + key + "-RUN-01.md"
	writeRun(t, harness.app.ProjectDir, runPath, taskID, "01")
	harness.run(t, "append", eventExecuted, "--task-id", taskID, "--role", "Executor", "--summary", "Run complete", "--path", runPath)

	wrongReview := ".agent-relay/runs/" + key + "-REVIEW-02.md"
	writeReview(t, harness.app.ProjectDir, wrongReview, taskID, "02", "ready-for-user-decision")
	if err := harness.fail("append", eventReview, "--task-id", taskID, "--role", "Planner", "--summary", "Wrong round", "--path", wrongReview); err == nil {
		t.Fatal("REVIEW round mismatch passed")
	}
	reviewPath := ".agent-relay/runs/" + key + "-REVIEW-01.md"
	writeReview(t, harness.app.ProjectDir, reviewPath, taskID, "01", "ready-for-user-decision")
	harness.run(t, "append", eventReview, "--task-id", taskID, "--role", "Planner", "--summary", "Review complete", "--path", reviewPath)

	closePath := ".agent-relay/runs/" + key + "-CLOSE.md"
	writeClose(t, harness.app.ProjectDir, closePath, "wrong")
	if err := harness.fail("append", eventClose, "--task-id", taskID, "--role", "Director", "--summary", "Wrong task", "--path", closePath); err == nil {
		t.Fatal("CLOSE task-id mismatch passed")
	}
	writeClose(t, harness.app.ProjectDir, closePath, taskID)
	harness.run(t, "append", eventClose, "--task-id", taskID, "--role", "Director", "--summary", "Consistency flow closed", "--path", closePath)

	cleanLog := readFile(t, harness.app.relayPath("relay.log"))
	log := strings.ReplaceAll(string(cleanLog), runPath, wrongRun)
	if err := os.WriteFile(harness.app.relayPath("relay.log"), []byte(log), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := harness.fail("lint"); err == nil {
		t.Fatal("lint accepted an artifact key mismatch")
	}
	if err := os.WriteFile(harness.app.relayPath("relay.log"), cleanLog, 0o644); err != nil {
		t.Fatal(err)
	}
	review := strings.ReplaceAll(string(readFile(t, harness.app.projectPath(reviewPath))), "Task ID: "+taskID, "Task ID: tampered")
	writeArtifact(t, harness.app.ProjectDir, reviewPath, review)
	if err := harness.fail("lint"); err == nil {
		t.Fatal("lint accepted a task-id mismatch")
	}
}

func TestCommandSmoke(t *testing.T) {
	harness := newHarness(t)
	if err := harness.fail("new-round", "bad", "--summary", "invalid | summary"); err == nil {
		t.Fatal("summary delimiter passed")
	}
	prompt := harness.run(t, "prompt", "plan", "--task-id", "abcd", "--key", "20260711-1000-test")
	if !strings.Contains(prompt, "Planner for round `abcd`") {
		t.Fatalf("unexpected prompt: %s", prompt)
	}
	version := strings.TrimSpace(harness.run(t, "version"))
	if version == "" {
		t.Fatal("version output is empty")
	}
	if err := harness.fail("gate", "before-review", "--task-id", "missing"); err == nil {
		t.Fatal("gate accepted a missing task")
	}
}

func TestMergeAgentBlock(t *testing.T) {
	harness := newHarness(t)
	target := filepath.Join(harness.app.ProjectDir, "TARGET.md")
	source := filepath.Join(harness.app.ProjectDir, "SOURCE.md")
	if err := os.WriteFile(source, []byte("before\n<agent-relay-rules>\nnew\n</agent-relay-rules>\nafter\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(target, []byte("project rules\n<agent-relay-rules>\nold\n</agent-relay-rules>\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	harness.run(t, "merge-agent-block", target, source)
	merged := string(readFile(t, target))
	if strings.Contains(merged, "old") || !strings.Contains(merged, "new") || !strings.Contains(merged, "project rules") {
		t.Fatalf("unexpected merged file: %s", merged)
	}
	appendTarget := filepath.Join(harness.app.ProjectDir, "APPEND.md")
	if err := os.WriteFile(appendTarget, []byte("project rules\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	harness.run(t, "merge-agent-block", appendTarget, source)
	if !strings.Contains(string(readFile(t, appendTarget)), "<agent-relay-rules>") {
		t.Fatal("merge did not append a missing block")
	}
}
