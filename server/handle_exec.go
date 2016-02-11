package main

import (
	execstreamer "github.com/golang-devops/go-exec-streamer"
	"github.com/labstack/echo"

	"github.com/golang-devops/go-psexec/shared"
)

func (h *handler) handleExecFunc(c *echo.Context) error {
	dto := &shared.ExecDto{}
	h.deserializeBody(c.Request().Body, dto)

	h.logger.Infof("Starting command, exe = '%s', args = '%#v'", dto.Exe, dto.Args)

	streamer, err := execstreamer.NewExecStreamerBuilder().
		ExecutorName(dto.Executor).
		Exe(dto.Exe).
		Args(dto.Args...).
		Writers(c.Response()).
		StdoutPrefix("OUT1:").
		StderrPrefix("ERR1:").
		AutoFlush().
		Build()

	if err != nil {
		panic(err)
	}

	return streamer.ExecAndWait()
}
