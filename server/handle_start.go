package main

import (
	"fmt"
	"net/http"

	"github.com/go-zero-boilerplate/path_utils"
	execstreamer "github.com/golang-devops/go-exec-streamer"
	"github.com/labstack/echo"

	"github.com/golang-devops/go-psexec/shared"
	"github.com/golang-devops/go-psexec/shared/dtos"
)

func (h *handler) handleStartFunc(c echo.Context) error {
	req := c.Request()
	resp := c.Response()

	dto := &dtos.ExecDto{}
	err := h.getDto(c, dto)
	if err != nil {
		return err
	}

	if exeExists, err := path_utils.FileExists(dto.Exe); err != nil {
		return fmt.Errorf("Unable to check if Exe path '%s' exists, error: %s", dto.Exe, err.Error())
	} else if !exeExists {
		return c.String(http.StatusBadRequest, fmt.Sprintf("Exe path '%s' does not exist. Please specify the full Exe path to run.", dto.Exe))
	}

	ip := getIPFromRequest(req)
	hostNames, err := getHostNamesFromIP(ip)
	if err != nil {
		h.logger.Warningf("Unable to find hostname(s) for IP '%s', error: %s", ip, err.Error())
	}

	h.logger.Infof(
		"Starting command (remote ip %s, hostnames = %+v), exe = '%s', args = '%#v' (working dir '%s')",
		ip, hostNames, dto.Exe, dto.Args, dto.WorkingDir)

	executor, err := execstreamer.NewExecutorFromName(dto.Executor)
	if err != nil {
		return err
	}

	cmd := executor.GetCommand(dto.Exe, dto.Args...)
	cmd.Dir = dto.WorkingDir

	err = cmd.Start()
	if err != nil {
		return err
	}

	resp.Header().Set(shared.PROCESS_ID_HTTP_HEADER_NAME, fmt.Sprintf("%d", cmd.Process.Pid))
	return c.String(200, "The command was successfully started. Pid in header: "+shared.PROCESS_ID_HTTP_HEADER_NAME)
}
