package main

import (
	"html/template"
	"io"
)

type htmlTemplateRenderer struct {
}

func (t *htmlTemplateRenderer) Render(w io.Writer, name string, data interface{}) error {
	tpl := template.Must(template.ParseGlob("templates/*"))
	return tpl.ExecuteTemplate(w, name, data)
}
