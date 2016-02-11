package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"github.com/ayufan/golang-kardianos-service"
	"io"
)

type handler struct {
	logger     service.Logger
	privateKey *rsa.PrivateKey
}

func (h *handler) deserializeBody(body io.Reader, dest interface{}) {
	decoder := json.NewDecoder(body)
	err := decoder.Decode(dest)
	checkError(err)
}

func (h *handler) getPublicKeyBytes() ([]byte, error) {
	return x509.MarshalPKIXPublicKey(&h.privateKey.PublicKey)
}
