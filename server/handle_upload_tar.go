package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/labstack/echo"

	"github.com/golang-devops/go-psexec/shared"
	"github.com/golang-devops/go-psexec/shared/tar_io"
)

func (h *handler) handleUploadTarFunc(c *echo.Context) error {
	req := c.Request()
	basePath := req.Header.Get(shared.BASE_PATH_HTTP_HEADER_NAME)
	if strings.TrimSpace(basePath) == "" {
		return fmt.Errorf("Request does not contain header '%s'", shared.BASE_PATH_HTTP_HEADER_NAME)
	}
	isDirVal := req.Header.Get(shared.IS_DIR_HTTP_HEADER_NAME)
	if strings.TrimSpace(isDirVal) == "" {
		return fmt.Errorf("Request does not contain header '%s'", shared.IS_DIR_HTTP_HEADER_NAME)
	}
	isDir := strings.TrimSpace(isDirVal) == "1"

	if isDir {
		//TODO: Is permission fine?
		if err := os.MkdirAll(basePath, 0755); err != nil {
			return fmt.Errorf("Unable to create base path %s, error: %s", basePath, err.Error())
		}
	}

	if isDir {
		h.logger.Infof("Now starting to receive dir '%s'", basePath)
		tarReceiver := tar_io.Factories.TarReceiver.Dir(basePath)
		return tar_io.SaveTarToReceiver(req.Body, tarReceiver)
	} else {
		h.logger.Infof("Now starting to receive file '%s'", basePath)
		tarReceiver := tar_io.Factories.TarReceiver.File(basePath)
		return tar_io.SaveTarToReceiver(req.Body, tarReceiver)
	}
}
