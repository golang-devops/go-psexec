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

		sessionId, sessionToken, err := newSessionToken(clientPubKey)
		if err != nil {
			return err
		}

		encryptedSessionToken, encryptedTokenSignature, err := shared.EncryptWithPublicKey(clientPubKey, h.key, sessionToken)
		if err != nil {
			return err
		}

		serverPubKeyBytes, err := h.getPublicKeyBytes()
		if err != nil {
			return err
		}

		jsonMessage, err := json.Marshal(&shared.GenTokenResponseMessage{
			sessionId,
			encryptedTokenSignature,
			serverPubKeyBytes,
		})
		if err != nil {
			return err
		}

		encryptedJsonMessage, err := shared.EncryptSymmetric(sessionToken, jsonMessage)
		if err != nil {
			return err
		}

		return c.JSON(
			200,
			&shared.GenTokenResponseDto{
				encryptedSessionToken,
				encryptedJsonMessage,
			},
		)
	}

	return nil
}