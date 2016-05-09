package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/labstack/echo"

	"github.com/golang-devops/go-psexec/shared/dtos"
)

func (h *handler) handlePathSummaryFunc(c *echo.Context) error {
	path := c.Query("path")
	if strings.TrimSpace(path) == "" {
		return fmt.Errorf("Request does not contain query 'path' value")
	}

	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			//TODO: Return empty dir instead of NotFound so that on the other side it will result in a "out of sync" status?
			//return c.JSON(http.StatusNotFound, &dtos.FilesystemSummaryDto{
			return c.JSON(200, &dtos.FilesystemSummaryDto{
				FlattenedFiles: []*dtos.FileSummary{
					&dtos.FileSummary{},
				},
			})
		}
		return fmt.Errorf("Unable to get stats of path '%s', error: %s", path, err.Error())
	}

	baseDir := ""
	flattenedFiles := []*dtos.FileSummary{}
	if info.IsDir() {
		dirSummary, err := h.svcs.FilePathSummaries.GetDirSummary(path)
		if err != nil {
			return fmt.Errorf("Cannot get dir '%s' summary, error: %s", path, err.Error())
		}

		baseDir = path
		for _, f := range dirSummary.FlattenedFileSummaries {
			relPath := strings.TrimLeft(f.FullPath[len(baseDir):], "\\/")
			flattenedFiles = append(flattenedFiles, &dtos.FileSummary{
				RelativePath: relPath,
				ModTime:      f.ModTime,
				ChecksumHex:  f.Checksum.HexString(),
			})
		}
	} else {
		fileSummary, err := h.svcs.FilePathSummaries.GetFileSummary(path)
		if err != nil {
			return fmt.Errorf("Cannot get file '%s' summary, error: %s", path, err.Error())
		}

		baseDir = ""
		relPath := fileSummary.FullPath //Use full path and keep base dir empty
		flattenedFiles = append(flattenedFiles, &dtos.FileSummary{
			RelativePath: relPath,
			ModTime:      fileSummary.ModTime,
			ChecksumHex:  fileSummary.Checksum.HexString(),
		})
	}
	returnDto := &dtos.FilesystemSummaryDto{
		BaseDir:        baseDir,
		FlattenedFiles: flattenedFiles,
	}
	return c.JSON(200, returnDto)
}
