package controllers

import (
	"net/http"

	"github.com/jerilcj30/jocp/views"
)

type StaticHandler struct {
	Templates struct {
		New views.Template
	}
}

func (h StaticHandler) New(w http.ResponseWriter, r *http.Request) {
	h.Templates.New.Execute(w, r, nil)
}
