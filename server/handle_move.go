package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/labstack/echo"

	"github.com/golang-devops/go-psexec/shared/dtos"
)

func (h *handler) handleMoveFunc(c *echo.Context) error {
	dto := &dtos.FsMoveDto{}
	err := h.getDto(c, dto)
	if err != nil {
		return err
	}

	newParentDir := filepath.Dir(dto.NewRemotePath)
	err = os.MkdirAll(newParentDir, 0755)
	if err != nil {
		return fmt.Errorf("Unable to move(rename) path '%s' to '%s'. The parent dir ('%s') could not be created, error: %s", dto.OldRemotePath, dto.NewRemotePath, newParentDir, err.Error())
	}

	err = os.Rename(dto.OldRemotePath, dto.NewRemotePath)
	if err != nil {
		return fmt.Errorf("Unable to move(rename) path '%s' to '%s', error: %s", dto.OldRemotePath, dto.NewRemotePath, err.Error())
	}

	return nil
}
