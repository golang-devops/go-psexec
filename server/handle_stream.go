package main

import (
	"fmt"
	"net/http"

	execstreamer "github.com/golang-devops/go-exec-streamer"
	"github.com/labstack/echo"

	"github.com/golang-devops/go-psexec/shared"
	"github.com/golang-devops/go-psexec/shared/dtos"
)

func (h *handler) handleStreamFunc(c *echo.Context) error {
	req := c.Request()
	resp := c.Response()

	dto := &dtos.ExecDto{}
	err := h.getDto(c, dto)
	if err != nil {
		return err
	}

	ip := getIPFromRequest(req)
	hostNames, err := getHostNamesFromIP(ip)
	if err != nil {
		h.logger.Warningf("Unable to find hostname(s) for IP '%s', error: %s", ip, err.Error())
	}

	h.logger.Infof(
		"Starting to stream command (remote ip %s, hostnames = %+v), exe = '%s', args = '%#v' (working dir '%s')",
		ip, hostNames, dto.Exe, dto.Args, dto.WorkingDir)

	resp.Header().Set("Content-Type", "application/octet-stream")
	resp.Header().Set("Transfer-Encoding", "chunked")

	onStarted := func(startedDetails *execstreamer.StartedDetails) {
		resp.Header().Set(shared.PROCESS_ID_HTTP_HEADER_NAME, fmt.Sprintf("%d", startedDetails.Pid))
		resp.WriteHeader(http.StatusOK)
		resp.Flush()
	}

	streamer, err := execstreamer.NewExecStreamerBuilder().
		ExecutorName(dto.Executor).
		Exe(dto.Exe).
		Args(dto.Args...).
		Dir(dto.WorkingDir).
		Writers(resp). //Writers(sessionToken.NewEncryptedWriter(resp)).
		// StdoutPrefix("OUT1:").
		StderrPrefix("ERROR:").
		AutoFlush().
		DebugInfo(fmt.Sprintf("ARGS=%+v", dto.Args)).
		OnStarted(onStarted).
		Build()

	if err != nil {
		return err
	}

	err = streamer.ExecAndWait()
	if err != nil {
		h.logger.Warningf("Unable to execute command, error was: %s", err.Error())
		return err
	}

	resp.Flush()
	resp.Write([]byte(shared.RESPONSE_EOF))
	resp.Flush()

	return nil
}
