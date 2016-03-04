package execstreamer

import (
	"errors"
	"os/exec"
	"strings"
)

func NewExecutorFromName(name string) (*executor, error) {
	cleanName := strings.ToLower(strings.TrimSpace(name))
	switch cleanName {
	case "winshell":
		return &executor{"cmd", []string{"/c"}}, nil
	case "bash":
		return &executor{"bash", []string{"-c"}}, nil
	default:
		return nil, errors.New("Executor not supported, name: " + name)
	}
}

type executor struct {
	exe  string
	args []string
}

func (e *executor) getFinalExeAndArgs(cmdExe string, cmdArgs ...string) (finalExe string, finalArgs []string) {
	if strings.TrimSpace(e.exe) == "" {
		return cmdExe, cmdArgs
	}

	finalExe = e.exe
	finalArgs = e.args
	finalArgs = append(finalArgs, cmdExe)
	finalArgs = append(finalArgs, cmdArgs...)
	return
}

func (e *executor) GetCommand(cmdExe string, cmdArgs ...string) *exec.Cmd {
	finalExe, finalArgs := e.getFinalExeAndArgs(cmdExe, cmdArgs...)

	return exec.Command(finalExe, finalArgs...)
}
