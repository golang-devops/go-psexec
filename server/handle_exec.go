package main

import (
	"bufio"
	"fmt"
	"github.com/labstack/echo"
)

func (h *handler) handleExecFunc(c *echo.Context) error {
	dto := h.deserializeBody(c.Request().Body)

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
			fmt.Fprintf(c.Response(), "O:%s\n", stdoutScanner.Text())
			c.Response().Flush()
		}
	}()

	stderrScanner := bufio.NewScanner(stderr)
	go func() {
		for stderrScanner.Scan() {
			fmt.Fprintf(c.Response(), "E:%s\n", stderrScanner.Text())
			c.Response().Flush()
		}
	}()

	err = cmd.Start()
	checkError(err)

	err = cmd.Wait()
	checkError(err)

	return nil
}
