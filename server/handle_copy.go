package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/labstack/echo"

	"github.com/golang-devops/go-psexec/shared/dtos"
	"github.com/termie/go-shutil"
)

func (h *handler) handleCopyFunc(c echo.Context) error {
	dto := &dtos.FsCopyDto{}
	err := h.getDto(c, dto)
	if err != nil {
		return err
	}

	destParentDir := filepath.Dir(dto.DestRemotePath)
	err = os.MkdirAll(destParentDir, 0755)
	if err != nil {
		return fmt.Errorf("Unable to copy path '%s' to '%s'. The parent dir ('%s') could not be created, error: %s", dto.SrcRemotePath, dto.DestRemotePath, destParentDir, err.Error())
	}

	srcStats, err := os.Stat(dto.SrcRemotePath)
	if err != nil {
		return fmt.Errorf("Unable to get stats of source path '%s' to copy, error: %s", dto.SrcRemotePath, err.Error())
	}

	if srcStats.IsDir() {
		if err = shutil.CopyTree(dto.SrcRemotePath, dto.DestRemotePath, nil); err != nil {
			return fmt.Errorf("Unable to copy tree from '%s' to '%s', error: %s", dto.SrcRemotePath, dto.DestRemotePath, err.Error())
		}
	} else {
		if err = shutil.CopyFile(dto.SrcRemotePath, dto.DestRemotePath, false); err != nil {
			return fmt.Errorf("Unable to copy file from '%s' to '%s', error: %s", dto.SrcRemotePath, dto.DestRemotePath, err.Error())
		}
	}

	return nil
}
