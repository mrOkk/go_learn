package handler

import (
	"UrlShortener/internal/domain"
	"UrlShortener/internal/service"
	"fmt"
	"html/template"
	"net/http"
)

const shorteningFailed = "Shortening failed"

type Handler struct {
	tmpl   *template.Template
	urlSrv *service.UrlService
}

func NewHandler(tmpl *template.Template, urlSrv *service.UrlService) *Handler {
	return &Handler{
		tmpl:   tmpl,
		urlSrv: urlSrv,
	}
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	err := h.tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) Shorten(w http.ResponseWriter, r *http.Request) {
	if r.ParseForm() != nil {
		http.Error(w, shorteningFailed, http.StatusBadRequest)
		return
	}
	url := r.PostFormValue("url")
	fmt.Println(url)

	code, err := h.urlSrv.Put(r.Context(), url)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.tmpl.Execute(w, domain.PageData{ShortURL: code})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) RedirectByPath(w http.ResponseWriter, r *http.Request) {
	fmt.Println("RedirectByPath page | ", r.URL)
	h.Redirect(w, r, r.PathValue("code"))
}

func (h *Handler) RedirectByQuery(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	h.Redirect(w, r, code)
}

func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request, code string) {
	url, err := h.urlSrv.Get(r.Context(), code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, url, http.StatusSeeOther)
}
