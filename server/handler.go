package main

import (
	"crypto/rsa"
	"encoding/json"
	"github.com/ayufan/golang-kardianos-service"
	"io"
)

type handler struct {
	logger service.Logger
	key    *rsa.PrivateKey
}

func (h *handler) deserializeBody(body io.Reader, dest interface{}) {
	decoder := json.NewDecoder(body)
	err := decoder.Decode(dest)
	checkError(err)
}
