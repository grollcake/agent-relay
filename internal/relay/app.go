package relay

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type App struct {
	ProjectDir string
	RelayDir   string
	Stdout     io.Writer
	Stderr     io.Writer
	Now        func() time.Time
	Random     io.Reader
	GOOS       string
	GOARCH     string
}

func New(projectDir, relayDir string) *App {
	return &App{
		ProjectDir: projectDir,
		RelayDir:   relayDir,
		Stdout:     os.Stdout,
		Stderr:     os.Stderr,
		Now:        time.Now,
		Random:     rand.Reader,
		GOOS:       runtime.GOOS,
		GOARCH:     runtime.GOARCH,
	}
}

func Discover() (*App, error) {
	if configured := os.Getenv("AGENT_RELAY_DIR"); configured != "" {
		relayDir, err := filepath.Abs(configured)
		if err != nil {
			return nil, err
		}
		return New(filepath.Dir(relayDir), relayDir), nil
	}

	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	for dir := cwd; ; dir = filepath.Dir(dir) {
		relayDir := filepath.Join(dir, ".agent-relay")
		if info, statErr := os.Stat(relayDir); statErr == nil && info.IsDir() {
			return New(dir, relayDir), nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
	}

	executable, err := os.Executable()
	if err == nil {
		binDir := filepath.Dir(executable)
		if filepath.Base(binDir) == "bin" {
			relayDir := filepath.Dir(binDir)
			return New(filepath.Dir(relayDir), relayDir), nil
		}
	}
	return nil, errors.New("cannot locate .agent-relay; run from a project or set AGENT_RELAY_DIR")
}

func (a *App) Run(args []string) error {
	if len(args) == 0 {
		a.usage()
		return errors.New("missing command")
	}

	command := args[0]
	args = args[1:]
	switch command {
	case "append":
		return a.runAppend(args)
	case "gate":
		return a.runGate(args)
	case "new-round":
		return a.runNewRound(args)
	case "feedback":
		return a.runFeedback(args)
	case "status":
		return a.runStatus(args)
	case "subagent-prompt", "prompt":
		return a.runPrompt(args)
	case "check-artifact":
		return a.runCheckArtifact(args)
	case "lint":
		return a.Lint()
	case "merge-agent-block":
		return a.runMergeAgentBlock(args)
	case "update":
		return a.runUpdate(args)
	case "version":
		return a.runVersion(args)
	case "help", "-h", "--help":
		a.usage()
		return nil
	default:
		a.usage()
		return fmt.Errorf("unknown command: %s", command)
	}
}

func (a *App) usage() {
	fmt.Fprint(a.Stdout, `Usage:
  agent-relay append <EVENT> --task-id <id> --role <role> --summary <text> [--path <path>]
  agent-relay gate <before-execute|before-review|before-approval> --task-id <id>
  agent-relay new-round <slug> --summary <text> [--branch <name>]
  agent-relay feedback --task-id <id> --summary <text>
  agent-relay status [--task-id <id>]
  agent-relay prompt <plan|review|exec> --task-id <id> --key <key> [--run-number <NN>]
  agent-relay check-artifact <PLANNED|EXECUTED|REVIEW|CLOSE> <path> [task-id]
  agent-relay lint
  agent-relay merge-agent-block <target-file> <source-file>
  agent-relay update --upstream <agent-relay-repo> [--apply]
  agent-relay version
`)
}

func (a *App) relayPath(parts ...string) string {
	return filepath.Join(append([]string{a.RelayDir}, parts...)...)
}

func (a *App) projectPath(path string) string {
	return filepath.Join(a.ProjectDir, filepath.FromSlash(path))
}

func (a *App) runGit(args ...string) (string, error) {
	command := exec.Command("git", args...)
	command.Dir = a.ProjectDir
	output, err := command.CombinedOutput()
	if err != nil {
		message := strings.TrimSpace(string(output))
		if message == "" {
			message = err.Error()
		}
		return "", errors.New(message)
	}
	return strings.TrimSpace(string(output)), nil
}

func binaryName(goos string) string {
	if goos == "windows" {
		return "agent-relay.exe"
	}
	return "agent-relay"
}

func (a *App) installedBinaryPath() string {
	return a.relayPath("bin", binaryName(a.GOOS))
}
