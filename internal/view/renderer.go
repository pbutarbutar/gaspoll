package view

import (
	"embed"
	"html/template"
	"io"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

//go:embed templates/*.html
var templateFS embed.FS

type Renderer struct {
	tmpl  *template.Template
	funcs template.FuncMap
}

func NewRenderer() *Renderer {
	funcs := template.FuncMap{
		"currency": func(amount int) string {
			return formatIDR(amount)
		},
		"discounted": func(base int, pct int) int {
			if base <= 0 || pct <= 0 {
				if pct >= 100 {
					return 0
				}
				return base
			}
			d := base * (100 - pct) / 100
			return d
		},
		"slice": func(v ...interface{}) []interface{} {
			return v
		},
		"percent": func(cur, goal int) int {
			if goal == 0 {
				return 0
			}
			pct := cur * 100 / goal
			if pct > 100 {
				return 100
			}
			return pct
		},
		"formatDate": func(t time.Time) string {
			if t.IsZero() {
				return "-"
			}
			return t.Format("02 Jan 2006 15:04")
		},
	}
	t := template.New("").Funcs(funcs)
	// Parse layout first so page `{{ define "content" }}` blocks parsed later
	// will override the layout's content block.
	t = template.Must(t.ParseFS(templateFS, "templates/layout.html"))

	// Read template directory and parse other templates deterministically.
	entries, err := templateFS.ReadDir("templates")
	if err == nil {
		for _, e := range entries {
			name := e.Name()
			if name == "layout.html" {
				continue
			}
			t = template.Must(t.ParseFS(templateFS, "templates/"+name))
		}
	}

	return &Renderer{tmpl: t, funcs: funcs}
}

func (r *Renderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	// Build a fresh template set for this request containing only layout and
	// the requested page so the page's `content` overrides deterministically.
	t := template.New("").Funcs(r.funcs)
	if _, err := t.ParseFS(templateFS, "templates/layout.html"); err != nil {
		return err
	}
	if _, err := t.ParseFS(templateFS, "templates/"+name); err != nil {
		return err
	}
	return t.ExecuteTemplate(w, "layout", data)
}

func formatIDR(amount int) string {
	if amount == 0 {
		return "Rp0"
	}
	s := []byte{}
	tmp := amount
	for tmp > 0 {
		s = append([]byte{byte('0' + tmp%10)}, s...)
		tmp /= 10
	}
	str := string(s)
	n := len(str)
	if n <= 3 {
		return "Rp" + str
	}
	parts := []string{}
	for n > 3 {
		parts = append([]string{str[n-3:]}, parts...)
		str = str[:n-3]
		n = len(str)
	}
	parts = append([]string{str}, parts...)
	return "Rp" + strings.Join(parts, ".")
}
