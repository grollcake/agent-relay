package memento

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
	MementoDir string
	Stdout     io.Writer
	Stderr     io.Writer
	Now        func() time.Time
	Random     io.Reader
	GOOS       string
	GOARCH     string
}

func New(projectDir, mementoDir string) *App {
	return &App{
		ProjectDir: projectDir,
		MementoDir: mementoDir,
		Stdout:     os.Stdout,
		Stderr:     os.Stderr,
		Now:        time.Now,
		Random:     rand.Reader,
		GOOS:       runtime.GOOS,
		GOARCH:     runtime.GOARCH,
	}
}

func Discover() (*App, error) {
	if configured := os.Getenv("MEMENTO_DIR"); configured != "" {
		mementoDir, err := filepath.Abs(configured)
		if err != nil {
			return nil, err
		}
		return New(filepath.Dir(mementoDir), mementoDir), nil
	}

	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	for dir := cwd; ; dir = filepath.Dir(dir) {
		mementoDir := filepath.Join(dir, ".memento")
		if info, statErr := os.Stat(mementoDir); statErr == nil && info.IsDir() {
			return New(dir, mementoDir), nil
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
			mementoDir := filepath.Dir(binDir)
			return New(filepath.Dir(mementoDir), mementoDir), nil
		}
	}
	return nil, errors.New("cannot locate .memento; run from a project or set MEMENTO_DIR")
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
  memento append <EVENT> --task-id <id> --role <role> --summary <text> [--path <path>]
  memento gate <before-execute|before-review|before-approval> --task-id <id>
  memento new-round <slug> --summary <text> [--branch <name>]
  memento feedback --task-id <id> --summary <text>
  memento status [--task-id <id>]
  memento prompt <plan|review|exec> --task-id <id> --key <key> [--run-number <NN>]
  memento check-artifact <PLANNED|EXECUTED|REVIEW|CLOSE> <path> [task-id]
  memento lint
  memento merge-agent-block <target-file> <source-file>
  memento update --upstream <memento-repo> [--apply]
  memento version
`)
}

func (a *App) mementoPath(parts ...string) string {
	return filepath.Join(append([]string{a.MementoDir}, parts...)...)
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
		return "memento.exe"
	}
	return "memento"
}

func (a *App) installedBinaryPath() string {
	return a.mementoPath("bin", binaryName(a.GOOS))
}
