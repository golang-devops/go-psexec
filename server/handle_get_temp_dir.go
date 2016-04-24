package main

import (
	"os"

	"github.com/labstack/echo"

	"github.com/golang-devops/go-psexec/shared/dtos"
)

func (h *handler) handleGetTempDirFunc(c *echo.Context) error {
	dto := &dtos.TempDirDto{TempDir: os.TempDir()}
	return c.JSON(200, dto)
}
