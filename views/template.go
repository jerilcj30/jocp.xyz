package views

import (
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"

	"github.com/jerilcj30/jocp/context"

	"github.com/gorilla/csrf"
	"github.com/jerilcj30/jocp/models"
)

// will have all the tools to render the template
// code which is used to render a template
// code to make the tempplate render on a http.response writer

type Template struct {
	HTMLTpl *template.Template
}

func Must(t Template, err error) Template {
	if err != nil {
		panic(err)
	}
	return t
}

func ParseFS(fs fs.FS, patterns ...string) (Template, error) {
	tpl := template.New(patterns[0])
	tpl = tpl.Funcs(
		template.FuncMap{
			"csrfField": func() (template.HTML, error) {
				return "", fmt.Errorf("csrf not implemented")
			},
			"currentUser": func() (template.HTML, error) {
				return "", fmt.Errorf("currentUser not implemented")
			},
			"successMessage": func() (string, error) {
				return "", fmt.Errorf("success message not implemented ")
			},
		},
	)
	tpl, err := tpl.ParseFS(fs, patterns...)
	if err != nil {
		return Template{}, fmt.Errorf("parsing template: %w", err)
	}

	return Template{
		HTMLTpl: tpl,
	}, nil
}

func (t Template) Execute(w http.ResponseWriter, r *http.Request, data interface{}) {
	// Read the success message from the cookie
	successMessage := ""
	cookie, err := r.Cookie("successMessage")
	if err == nil {
		successMessage = cookie.Value
	}
	tpl := t.HTMLTpl
	tpl = tpl.Funcs(
		template.FuncMap{
			"csrfField": func() template.HTML {
				return csrf.TemplateField(r)
			},
			"currentUser": func() *models.UserTable {
				return context.User(r.Context())
			},
			"successMessage": func() string {
				return successMessage
			},
		},
	)
	w.Header().Set("Contect-Type", "text/html; charset=utf-8")
	err = tpl.Execute(w, data)
	if err != nil {
		log.Printf("executing template: %v", err)
		http.Error(w, "There was an error executing the template.", http.StatusInternalServerError)
		return
	}
}
