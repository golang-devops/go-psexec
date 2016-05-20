package main

import (
	"encoding/json"

	"github.com/labstack/echo"

	"github.com/golang-devops/go-psexec/shared/dtos"
)

func (h *handler) getDto(c echo.Context, jsonDest interface{}) error {
	req := c.Request()

	sessionToken, err := h.getAuthenticatedSessionToken(c)
	if err != nil {
		return err
	}

	encryptedJsonContainer := &dtos.EncryptedJsonContainer{}
	err = h.deserializeBody(req.Body(), encryptedJsonContainer)
	if err != nil {
		return err
	}

	decryptedJson, err := sessionToken.DecryptWithSessionToken(encryptedJsonContainer.EncryptedJson)
	if err != nil {
		return err
	}

	return json.Unmarshal(decryptedJson, jsonDest)
}
