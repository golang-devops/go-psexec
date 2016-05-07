package main

import (
	"github.com/golang-devops/go-psexec/shared/dtos"

	"github.com/labstack/echo"
)

func (h *handler) handleVersionFunc(c *echo.Context) error {
	dto := &dtos.VersionDto{Version: TempVersion}
	return c.JSON(200, dto)
}
