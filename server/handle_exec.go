package main

import (
	"encoding/json"
	"fmt"
	execstreamer "github.com/golang-devops/go-exec-streamer"
	"github.com/labstack/echo"
	"net/http"

	"github.com/golang-devops/go-psexec/shared"
)

func (h *handler) handleExecFunc(c *echo.Context) error {
	sessionToken, err := h.getAuthenticatedSessionToken(c)
	if err != nil {
		return err
	}

	encryptedJsonContainer := &shared.EncryptedJsonContainer{}
	err = h.deserializeBody(c.Request().Body, encryptedJsonContainer)
	if err != nil {
		return err
	}

	decryptedJson, err := sessionToken.DecryptWithSessionToken(encryptedJsonContainer.EncryptedJson)
	if err != nil {
		return err
	}

	dto := &shared.ExecDto{}
	err = json.Unmarshal(decryptedJson, dto)
	if err != nil {
		return err
	}

	ip := getIPFromRequest(c.Request())
	hostNames, err := getHostNamesFromIP(ip)
	if err != nil {
		h.logger.Warningf("Unable to find hostname(s) for IP '%s', error: %s", ip, err.Error())
	}

	h.logger.Infof(
		"Starting command (remote ip %s, hostnames = %+v), exe = '%s', args = '%#v'",
		ip, hostNames, dto.Exe, dto.Args)

	c.Response().Header().Set("Content-Type", "application/octet-stream")
	c.Response().Header().Set("Transfer-Encoding", "chunked")
	c.Response().WriteHeader(http.StatusOK)

	streamer, err := execstreamer.NewExecStreamerBuilder().
		ExecutorName(dto.Executor).
		Exe(dto.Exe).
		Args(dto.Args...).
		Writers(c.Response()). //Writers(sessionToken.NewEncryptedWriter(c.Response())).
		// StdoutPrefix("OUT1:").
		StderrPrefix("ERROR:").
		AutoFlush().
		DebugInfo(fmt.Sprintf("ARGS=%+v", dto.Args)).
		Build()

	if err != nil {
		return err
	}

	err = streamer.ExecAndWait()
	if err != nil {
		h.logger.Warningf("Unable to execute command, error was: %s", err.Error())
		return err
	}

	c.Response().Flush()
	c.Response().Write([]byte(shared.RESPONSE_EOF))
	c.Response().Flush()

	return nil
}
