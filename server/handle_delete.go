package main

import (
	"fmt"
	"os"

	"github.com/labstack/echo"

	"github.com/golang-devops/go-psexec/shared/dtos"
)

func (h *handler) handleDeleteFunc(c *echo.Context) error {
	dto := &dtos.FsDeleteDto{}
	err := h.getDto(c, dto)
	if err != nil {
		return err
	}

	err = os.RemoveAll(dto.Path)
	if err != nil {
		return fmt.Errorf("Unable to delete path '%s', error: %s", dto.Path, err.Error())
	}

	return nil
}
