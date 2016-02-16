package main

import (
	"encoding/json"
	execstreamer "github.com/golang-devops/go-exec-streamer"
	"github.com/labstack/echo"

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

	descryptedJson, err := sessionToken.DecryptWithServerPrivateKey(h.privateKey, encryptedJsonContainer.EncryptedJson)
	if err != nil {
		return err
	}

	dto := &shared.ExecDto{}
	err = json.Unmarshal(descryptedJson, dto)
	if err != nil {
		return err
	}

	h.logger.Infof("Starting command, exe = '%s', args = '%#v'", dto.Exe, dto.Args)

	streamer, err := execstreamer.NewExecStreamerBuilder().
		ExecutorName(dto.Executor).
		Exe(dto.Exe).
		Args(dto.Args...).
		Writers(c.Response()).
		// StdoutPrefix("OUT1:").
		StderrPrefix("ERROR:").
		AutoFlush().
		Build()

	if err != nil {
		panic(err)
	}

	return streamer.ExecAndWait()
}
