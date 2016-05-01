package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo"

	"github.com/golang-devops/go-psexec/shared/dtos"
)

func (h *handler) handleStatsFunc(c *echo.Context) error {
	path := c.Query("path")
	if strings.TrimSpace(path) == "" {
		return fmt.Errorf("Request does not contain query 'path' value")
	}

	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return c.JSON(http.StatusNotFound, &dtos.StatsDto{
				Path:   path,
				Exists: false,
			})
		}
		return fmt.Errorf("Unable to get stats of path '%s', error: %s", path, err.Error())
	}

	returnDto := &dtos.StatsDto{
		Path:    path,
		Exists:  true,
		IsDir:   info.IsDir(),
		ModTime: info.ModTime(),
		Mode:    info.Mode(),
		Size:    info.Size(),
	}
	return c.JSON(200, returnDto)
}
