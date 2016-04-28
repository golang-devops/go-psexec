package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"errors"
	"github.com/labstack/echo"
	"net/http"

	"github.com/golang-devops/go-psexec/shared"
	"github.com/golang-devops/go-psexec/shared/dtos"
)

func (h *handler) handleGenerateTokenFunc(c *echo.Context) error {
	dto := &dtos.GetTokenRequestDto{}
	err := h.deserializeBody(c.Request().Body, dto)
	if err != nil {
		return err
	}

	pubKeyInterface, err := x509.ParsePKIXPublicKey(dto.ClientPubPKIXBytes)
	if err != nil {
		return err
	}

	if clientPubKey, ok := pubKeyInterface.(*rsa.PublicKey); !ok {
		return errors.New("The public key received is in an incorrect format")
	} else {
		if !h.checkPubKeyAllowed(clientPubKey) {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}

		sessionId, sessionToken, err := tokenStore.NewSessionToken(clientPubKey)
		if err != nil {
			return err
		}

		encryptedSessionToken, err := shared.EncryptWithPublicKey(clientPubKey, sessionToken)
		if err != nil {
			return err
		}

		encryptedTokenSignature, err := shared.GenerateMessageSignature(h.privateKey, sessionToken)
		if err != nil {
			return err
		}

		serverPubKeyBytes, err := h.getPublicKeyBytes()
		if err != nil {
			return err
		}

		jsonMessage, err := json.Marshal(&dtos.GenTokenResponseMessage{
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
			&dtos.GenTokenResponseDto{
				encryptedSessionToken,
				encryptedJsonMessage,
			},
		)
	}
}
