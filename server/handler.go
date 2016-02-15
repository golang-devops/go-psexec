package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"github.com/ayufan/golang-kardianos-service"
	"github.com/labstack/echo"
	"io"

	"github.com/golang-devops/go-psexec/shared"
)

type handler struct {
	logger            service.Logger
	privateKey        *rsa.PrivateKey
	AllowedPublicKeys []*shared.AllowedPublicKey
}

func (h *handler) deserializeBody(body io.Reader, dest interface{}) {
	decoder := json.NewDecoder(body)
	err := decoder.Decode(dest)
	checkError(err)
}

func (h *handler) getPublicKeyBytes() ([]byte, error) {
	return x509.MarshalPKIXPublicKey(&h.privateKey.PublicKey)
}

func (h *handler) getAuthenticatedSessionToken(c *echo.Context) (*sessionToken, error) {
	sessionIdInterface := c.Get("session-id")
	sessionId, ok := sessionIdInterface.(int)
	if !ok {
		return nil, fmt.Errorf("Context session-id invalid format '%#v'", sessionIdInterface)
	}

	token, ok := tmpTokens[sessionId]
	if !ok {
		return nil, fmt.Errorf("Invalid token")
	}

	return token, nil
}
