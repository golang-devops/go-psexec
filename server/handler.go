package main

import (
	"encoding/json"
	"github.com/ayufan/golang-kardianos-service"
	"io"

	"github.com/golang-devops/go-psexec/shared"
)

type handler struct {
	logger service.Logger
}

func (h *handler) deserializeBody(body io.Reader) *shared.Dto {
	decoder := json.NewDecoder(body)

	dto := &shared.Dto{}
	err := decoder.Decode(dto)
	checkError(err)

	return dto
}
