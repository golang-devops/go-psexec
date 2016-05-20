package main

import (
	"html/template"
	"io"

	"github.com/labstack/echo"
)

type htmlTemplateRenderer struct {
}

func (t *htmlTemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	tpl := template.Must(template.ParseGlob("templates/*"))
	return tpl.ExecuteTemplate(w, name, data)
}
