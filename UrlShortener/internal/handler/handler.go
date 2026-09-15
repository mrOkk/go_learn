package handler

import (
	"fmt"
	"html/template"
	"net/http"
)

type Handler struct {
	tmpl *template.Template
}

func NewHandler(tmpl *template.Template) *Handler {
	return &Handler{tmpl: tmpl}
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Home page | ", r.URL)
	err := h.tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) Shorten(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Shorten page | ", r.URL)
}

func (h *Handler) RedirectByPath(w http.ResponseWriter, r *http.Request) {
	fmt.Println("RedirectByPath page | ", r.URL)
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func (h *Handler) RedirectByQuery(w http.ResponseWriter, r *http.Request) {
	fmt.Println("RedirectByQuery page | ", r.URL)
	http.Error(w, "not implemented", http.StatusNotImplemented)
}
