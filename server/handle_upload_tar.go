package main

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/labstack/echo"

	"github.com/golang-devops/go-psexec/shared"
)

func (h *handler) handleUploadTarFunc(c *echo.Context) error {
	req := c.Request()

	basePath := req.Header.Get(shared.BASE_PATH_HTTP_HEADER_NAME)
	if strings.TrimSpace(basePath) == "" {
		return fmt.Errorf("Request does not contain header '%s'", shared.BASE_PATH_HTTP_HEADER_NAME)
	}

	/*isDirStr := req.Header.Get(shared.IS_DIR_HTTP_HEADER_NAME)
	if strings.TrimSpace(isDirStr) == "" {
		return fmt.Errorf("Request does not contain header '%s'", shared.IS_DIR_HTTP_HEADER_NAME)
	}
	isDir := isDirStr == "1"*/

	tarProvider := tar.NewReader(req.Body)

	foundEndOfTar := false
	for {
		hdr, err := tarProvider.Next()
		if err != nil {
			if err == io.EOF {
				// end of tar archive
				break
			} else {
				return fmt.Errorf("Unable to read next tar chunk, error: %s", err.Error())
			}
		}

		if hdr.Name == shared.END_OF_TAR_FILENAME {
			foundEndOfTar = true
			continue
		}

		relativePath := hdr.Name

		if hdr.FileInfo().IsDir() {
			fullDestPath := filepath.Join(basePath, relativePath)
			defer os.Chtimes(fullDestPath, hdr.AccessTime, hdr.ModTime)

			os.MkdirAll(fullDestPath, os.FileMode(hdr.Mode))
		} else {
			fullDestPath := filepath.Join(basePath, relativePath)
			if val, ok := hdr.Xattrs["SINGLE_FILE_ONLY"]; ok && val == "1" {
				fullDestPath = basePath
			}

			parentDir := filepath.Dir(fullDestPath)
			err = os.MkdirAll(parentDir, os.FileMode(hdr.Mode))
			if err != nil {
				return fmt.Errorf("Unable to dir '%s', error: %s", parentDir, err.Error())
			}

			file, err := os.OpenFile(fullDestPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(hdr.Mode))
			if err != nil {
				return fmt.Errorf("Unable to open file '%s', error: %s", fullDestPath, err.Error())
			}

			defer func() {
				file.Close()
				os.Chtimes(fullDestPath, hdr.AccessTime, hdr.ModTime)
			}()

			_, err = io.Copy(file, tarProvider)
			if err != nil {
				return fmt.Errorf("Unable to copy stream to file '%s', error: %s", fullDestPath, err.Error())
			}
		}
	}

	if !foundEndOfTar {
		return fmt.Errorf("TAR EOF not found. Stream validation failed, something has gone wrong during the transfer.")
	}

	return nil
}
