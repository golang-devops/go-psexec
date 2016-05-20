package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"github.com/ayufan/golang-kardianos-service"
	"github.com/labstack/echo"
	"io"
)

type handler struct {
	logger     service.Logger
	privateKey *rsa.PrivateKey
    svcs *HandlerServices
}

func (h *handler) deserializeBody(body io.Reader, dest interface{}) error {
	decoder := json.NewDecoder(body)
	return decoder.Decode(dest)
}

func (h *handler) getPublicKeyBytes() ([]byte, error) {
	return x509.MarshalPKIXPublicKey(&h.privateKey.PublicKey)
}

func (h *handler) getAuthenticatedSessionToken(c echo.Context) (*sessionToken, error) {
	sessionIdInterface := c.Get("session-id")
	sessionId, ok := sessionIdInterface.(int)
	if !ok {
		return nil, fmt.Errorf("Context session-id invalid format '%#v'", sessionIdInterface)
	}

	token, ok := tokenStore.GetSessionToken(sessionId)
	if !ok {
		return nil, fmt.Errorf("Invalid token")
	}

	return token, nil
}
