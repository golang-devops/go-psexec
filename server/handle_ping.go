package main

import (
	"github.com/labstack/echo"

	"github.com/golang-devops/go-psexec/shared/dtos"
)

func (h *handler) handlePingFunc(c echo.Context) error {
	dto := &dtos.PingDto{Ping: "pong"}
	return c.JSON(200, dto)
}
