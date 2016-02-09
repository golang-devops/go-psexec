package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/ayufan/golang-kardianos-service"
	"io"
	"net/http"

	"github.com/golang-devops/go-psexec/shared"
)

type handler struct {
	logger service.Logger
}

func (h *handler) handleHttpPanic(w http.ResponseWriter) {
	if r := recover(); r != nil {
		errStr := ""
		switch t := r.(type) {
		case error:
			errStr = t.Error()
			break
		default:
			errStr = fmt.Sprintf("%#v", r)
		}

		http.Error(w, errStr, http.StatusInternalServerError)
		h.logger.Warning(errStr)
	}
}

func (h *handler) deserializeBody(body io.Reader) *shared.Dto {
	decoder := json.NewDecoder(body)

	dto := &shared.Dto{}
	err := decoder.Decode(dto)
	checkError(err)

	return dto
}

func (h *handler) handler(w http.ResponseWriter, r *http.Request) {
	defer h.handleHttpPanic(w)

	dto := h.deserializeBody(r.Body)

	flusher, hasFlusher := w.(http.Flusher)

	e := NewExecutorFromName(dto.Executor)
	h.logger.Infof("Starting command, exe = '%s', args = '%#v'", dto.Exe, dto.Args)
	cmd := e.GetCommand(dto.Exe, dto.Args...)

	stdout, err := cmd.StdoutPipe()
	checkError(err)

	stderr, err := cmd.StderrPipe()
	checkError(err)

	stdoutScanner := bufio.NewScanner(stdout)
	go func() {
		for stdoutScanner.Scan() {
			fmt.Fprintf(w, "OUT: %s\n", stdoutScanner.Text())
			if hasFlusher {
				flusher.Flush()
			}
		}
	}()

	stderrScanner := bufio.NewScanner(stderr)
	go func() {
		for stderrScanner.Scan() {
			fmt.Fprintf(w, "ERROR: %s\n", stderrScanner.Text())
			if hasFlusher {
				flusher.Flush()
			}
		}
	}()

	err = cmd.Start()
	checkError(err)

	err = cmd.Wait()
	checkError(err)
}
