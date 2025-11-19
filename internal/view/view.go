package view

import (
	"embed"
	"html/template"
	"io"
	"strings"

	"github.com/labstack/echo/v4"
)

//go:embed templates/*.html templates/**/*.html
var templateFS embed.FS

type Data struct {
	Title   string
	Current any
	Data    map[string]any
}

type Engine struct {
	templates *template.Template
}

func NewEngine() (*Engine, error) {
	funcs := template.FuncMap{
		"upper": strings.ToUpper,
	}

	tmpl, err := template.New("base.html").Funcs(funcs).ParseFS(templateFS, "templates/*.html", "templates/**/*.html")
	if err != nil {
		return nil, err
	}

	return &Engine{templates: tmpl}, nil
}

func (e *Engine) Render(w io.Writer, name string, data interface{}, _ echo.Context) error {
	return e.templates.ExecuteTemplate(w, name, data)
}
