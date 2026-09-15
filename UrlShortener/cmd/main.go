package main

import (
	handler2 "UrlShortener/internal/handler"
	"fmt"
	"html/template"
	"net/http"
)

func main() {
	fmt.Println("Shortener started")
	//ctx := context.Background()
	tmpl, parseErr := template.ParseFiles("templates/index.html")
	if parseErr != nil {
		panic(parseErr)
	}
	h := handler2.NewHandler(tmpl)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", h.Home)
	mux.HandleFunc("POST /shorten", h.Shorten)
	mux.HandleFunc("GET /{code}", h.RedirectByPath)
	mux.HandleFunc("GET /redirect", h.RedirectByQuery)
	srv := http.Server{Addr: ":8080", Handler: mux}
	srv.ListenAndServe()
}

func handleIndex(w http.ResponseWriter, r *http.Request) {

}
