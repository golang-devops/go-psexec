package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"errors"
	"github.com/labstack/echo"
	"net/http"

	"github.com/golang-devops/go-psexec/shared"
)

func (h *handler) handleGenerateTokenFunc(c *echo.Context) error {
	dto := &shared.GetTokenRequestDto{}
	h.deserializeBody(c.Request().Body, dto)

	pubKeyInterface, err := x509.ParsePKIXPublicKey(dto.ClientPubPKIXBytes)
	if err != nil {
		return err
	}

	if clientPubKey, ok := pubKeyInterface.(*rsa.PublicKey); !ok {
		return errors.New("The public key received is in an incorrect format")
	} else {
		if !checkPubKeyAllowed(clientPubKey) {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}

		serverPubPKIXBytes, err := x509.MarshalPKIXPublicKey(&h.key.PublicKey)
		if err != nil {
			return err
		}

		sessionToken := newSessionToken()
		dto := &shared.GenTokenResponseDto{serverPubPKIXBytes, sessionToken}

		jsonBytes, err := json.Marshal(dto)
		if err != nil {
			return err
		}

		encryptedJson, err := shared.EncryptWithPublicKey(clientPubKey, jsonBytes)
		if err != nil {
			return err
		}

		c.Response().WriteHeader(http.StatusOK)
		c.Response().Write(encryptedJson)
	}

	return nil
}
