package main

import (
	"fmt"
	"strings"

	"github.com/go-zero-boilerplate/path_utils"

	"github.com/labstack/echo"

	"github.com/golang-devops/go-psexec/shared/tar_io"
)

func (h *handler) handleDownloadTarFunc(c echo.Context) error {
	pathToSend := c.QueryParam("path")
	fileFilter := c.QueryParam("file-filter")
	if strings.TrimSpace(pathToSend) == "" {
		return fmt.Errorf("Request does not contain query 'path' value")
	}

	var tarProvider tar_io.TarProvider
	if isDir, err := path_utils.DirectoryExists(pathToSend); err != nil {
		return fmt.Errorf("Unable to determine if path '%s' is a directory, error: %s", pathToSend, err.Error())
	} else if isDir {
		tarProvider = tar_io.Factories.TarProvider.Dir(pathToSend, fileFilter)
		h.logger.Infof("Now starting to send dir '%s'", pathToSend)
	} else if isFile, err := path_utils.FileExists(pathToSend); err != nil {
		return fmt.Errorf("Unable to determine if path '%s' is a file, error: %s", pathToSend, err.Error())
	} else if isFile {
		tarProvider = tar_io.Factories.TarProvider.File(pathToSend)
		h.logger.Infof("Now starting to send file '%s'", pathToSend)
	} else {
		return fmt.Errorf("Path '%s' is not an existing file or directory", pathToSend)
	}

	handler := &sendTarHandler{writer: c.Response()}
	err := tar_io.UploadProvider(tarProvider, handler)
	if err != nil {
		return fmt.Errorf("Unable to send file, error: %s", err.Error())
	}

	return nil
}
