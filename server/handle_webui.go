package main

import (
	"github.com/labstack/echo"
)

func (h *handler) handleWebUIFunc(c *echo.Context) error {
	err := c.Render(200, "webui.gohtml", nil)
	return err
}
