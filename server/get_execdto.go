package main

import (
	"encoding/json"

	"github.com/labstack/echo"

	"github.com/golang-devops/go-psexec/shared"
)

func (h *handler) getExecDto(c *echo.Context) (*shared.ExecDto, error) {
	req := c.Request()

	sessionToken, err := h.getAuthenticatedSessionToken(c)
	if err != nil {
		return nil, err
	}

	encryptedJsonContainer := &shared.EncryptedJsonContainer{}
	err = h.deserializeBody(req.Body, encryptedJsonContainer)
	if err != nil {
		return nil, err
	}

	decryptedJson, err := sessionToken.DecryptWithSessionToken(encryptedJsonContainer.EncryptedJson)
	if err != nil {
		return nil, err
	}

	dto := &shared.ExecDto{}
	err = json.Unmarshal(decryptedJson, dto)
	if err != nil {
		return nil, err
	}

	return dto, nil
}
