package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/labstack/echo"

	"github.com/golang-devops/go-psexec/shared/dtos"
)

func (h *handler) handleSymlinkFunc(c echo.Context) error {
	dto := &dtos.FsSymlinkDto{}
	err := h.getDto(c, dto)
	if err != nil {
		return err
	}

	destParentDir := filepath.Dir(dto.DestRemoteSymlinkPath)
	err = os.MkdirAll(destParentDir, 0755)
	if err != nil {
		return fmt.Errorf("Unable to copy path '%s' to '%s'. The parent dir ('%s') could not be created, error: %s", dto.SrcRemotePath, dto.DestRemoteSymlinkPath, destParentDir, err.Error())
	}

	if err = os.Symlink(dto.SrcRemotePath, dto.DestRemoteSymlinkPath); err != nil {
		return fmt.Errorf("Cannot create symlink from '%s' to '%s', error: %s", dto.SrcRemotePath, dto.DestRemoteSymlinkPath, err.Error())
	}

	return nil
}
