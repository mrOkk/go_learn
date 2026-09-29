package httpfront

import (
	"LoadBalancer/internal/app"
	"html/template"
	"log"
	"net/http"
)

const shorteningFailed = "Shortening failed"

type HTTPHandler struct {
	tmpl       *template.Template
	picker     app.Picker
	backCaller app.BackendCaller
}

func NewHTTPHandler(tmplPath string, picker app.Picker, backCaller app.BackendCaller) *HTTPHandler {
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		log.Fatal(err)
	}
	return &HTTPHandler{
		tmpl:       tmpl,
		picker:     picker,
		backCaller: backCaller,
	}
}

func NewRouter(h *HTTPHandler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", h.Home)
	mux.HandleFunc("POST /shorten", h.Shorten)
	mux.HandleFunc("GET /favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})
	mux.HandleFunc("GET /{code}", h.RedirectByPath)
	mux.HandleFunc("GET /redirect", h.RedirectByQuery)

	return mux
}

func (h *HTTPHandler) Home(w http.ResponseWriter, r *http.Request) {
	err := h.tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *HTTPHandler) Shorten(w http.ResponseWriter, r *http.Request) {
	if r.ParseForm() != nil {
		http.Error(w, shorteningFailed, http.StatusBadRequest)
		return
	}

	url := r.Form.Get("url")
	b, err := h.picker.Pick()

	if err != nil {
		http.Error(w, shorteningFailed, http.StatusInternalServerError)
		return
	}

	code, err := h.backCaller.Shorten(r.Context(), b.Address, url)
	if err != nil {
		http.Error(w, shorteningFailed, http.StatusInternalServerError)
		return
	}

	err = h.tmpl.Execute(w, PageData{ShortURL: code})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *HTTPHandler) RedirectByPath(w http.ResponseWriter, r *http.Request) {
	h.Redirect(w, r, r.PathValue("code"))
}

func (h *HTTPHandler) RedirectByQuery(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	h.Redirect(w, r, code)
}

func (h *HTTPHandler) Redirect(w http.ResponseWriter, r *http.Request, code string) {
	b, err := h.picker.Pick()
	if err != nil {
		http.Error(w, shorteningFailed, http.StatusInternalServerError)
		return
	}
	url, err := h.backCaller.Resolve(r.Context(), b.Address, code)
	if err != nil {
		http.Error(w, shorteningFailed, http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
