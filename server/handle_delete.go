package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/labstack/echo"

	"github.com/golang-devops/go-psexec/shared/dtos"
)

func (h *handler) handleDeleteFunc(c *echo.Context) error {
	dto := &dtos.FsDeleteDto{}
	err := h.getDto(c, dto)
	if err != nil {
		return err
	}

	sanitizedPath := strings.Replace(strings.Trim(dto.Path, " /\\"), "\\", "/", -1)
	if sanitizedPath == "" || sanitizedPath == "." || sanitizedPath == ".." || sanitizedPath == "/" ||
		strings.HasPrefix(sanitizedPath, "./") || strings.HasPrefix(sanitizedPath, "../") {
		return fmt.Errorf("Empty paths or relative paths are not allowed, cleaned path was '%s'", sanitizedPath)
	}

	err = os.RemoveAll(dto.Path)
	if err != nil {
		return fmt.Errorf("Unable to delete path '%s', error: %s", sanitizedPath, err.Error())
	}

	return nil
}
