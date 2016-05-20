package main

import (
	"fmt"

	"github.com/go-zero-boilerplate/osvisitors"

	"github.com/labstack/echo"

	"github.com/golang-devops/go-psexec/shared/dtos"
)

func (h *handler) handleGetOsTypeFunc(c echo.Context) error {
	runtimeOsType, err := osvisitors.GetRuntimeOsType()
	if err != nil {
		return fmt.Errorf("Unable to get runtime OsType, error: %s", err.Error())
	}

	visitor := &osvisitors.GoOSNameVisitor{}
	runtimeOsType.Accept(visitor)
	dto := &dtos.OsTypeDto{Name: visitor.Name}
	return c.JSON(200, dto)
}
