package execstreamer

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"runtime/debug"
	"strings"
	"sync"
)

//ExecStreamer is the streamer interface (built by the ExecStreamerBuilder)
type ExecStreamer interface {
	StartExec() (*exec.Cmd, error)
	ExecAndWait() error
}

type execStreamer struct {
	stdOutAndErrWaitGroup *sync.WaitGroup

	ExecutorName string
	Exe          string
	Args         []string
	Dir          string
	Env          []string

	StdoutWriter io.Writer
	StdoutPrefix string

	StderrWriter io.Writer
	StderrPrefix string

	AutoFlush bool

	DebugInfo string
}

func (e *execStreamer) recoverPanic(description string) {
	if r := recover(); r != nil {
		defer recover()
		fmt.Println(fmt.Sprintf("Exec-Stream-Recovery (%s - debug info: %s): %T %+v. Stack: %s\n------END STACK---------\n", description, e.DebugInfo, r, r, strings.Replace(string(debug.Stack()), "\n", "\\n", -1)))
	}
}

func (e *execStreamer) flushIfEnabled(writer io.Writer) {
	if e.AutoFlush && writer != nil {
		if flusher, ok := writer.(http.Flusher); ok {
			defer e.recoverPanic("flushIfEnabled")
			if flusher != nil {
				flusher.Flush()
			}
		}
	}
}

func (e *execStreamer) handleStdout(stdoutScanner *bufio.Scanner) {
	defer e.recoverPanic("handleStdout")
	defer e.stdOutAndErrWaitGroup.Done()
	for stdoutScanner.Scan() {
		fmt.Fprintf(e.StdoutWriter, "%s%s\n", e.StdoutPrefix, stdoutScanner.Text())
		e.flushIfEnabled(e.StdoutWriter)
	}
}

func (e *execStreamer) handleStderr(stderrScanner *bufio.Scanner) {
	defer e.recoverPanic("handleStderr")
	defer e.stdOutAndErrWaitGroup.Done()
	for stderrScanner.Scan() {
		fmt.Fprintf(e.StderrWriter, "%s%s\n", e.StderrPrefix, stderrScanner.Text())
		e.flushIfEnabled(e.StderrWriter)
	}
}

//StartExec will execute the command using the given executor and return (without waiting for completion) with the exec.Cmd
func (e *execStreamer) StartExec() (*exec.Cmd, error) {
	x, err := NewExecutorFromName(e.ExecutorName)
	if err != nil {
		return nil, err
	}

	cmd := x.GetCommand(e.Exe, e.Args...)

	if e.Dir != "" {
		cmd.Dir = e.Dir
	}
	if len(e.Env) > 0 {
		cmd.Env = e.Env
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	e.stdOutAndErrWaitGroup = &sync.WaitGroup{}

	e.stdOutAndErrWaitGroup.Add(1)
	stdoutScanner := bufio.NewScanner(stdout)
	go e.handleStdout(stdoutScanner)

	e.stdOutAndErrWaitGroup.Add(1)
	stderrScanner := bufio.NewScanner(stderr)
	go e.handleStderr(stderrScanner)

	return cmd, nil
}

//ExecAndWait will execute the command using the given executor and wait until completion
func (e *execStreamer) ExecAndWait() error {
	cmd, err := e.StartExec()
	if err != nil {
		return err
	}

	e.stdOutAndErrWaitGroup.Wait()

	err = cmd.Wait()
	if err != nil {
		return err
	}

	return nil
}
