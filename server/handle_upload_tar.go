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

	info, err := os.Stat(basePath)
	if err != nil {
		return fmt.Errorf("Unable to obtain stats of path '%s', error: %s", basePath, err.Error())
	}

	if info.IsDir() {
		h.logger.Infof("Now starting to receive dir '%s'", basePath)
		tarReceiver := tar_io.Factories.TarReceiver.Dir(basePath)
		return tar_io.SaveTarToReceiver(req.Body, tarReceiver)
	} else {
		h.logger.Infof("Now starting to receive file '%s'", basePath)
		tarReceiver := tar_io.Factories.TarReceiver.File(basePath)
		return tar_io.SaveTarToReceiver(req.Body, tarReceiver)
	}
}
