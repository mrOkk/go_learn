package main

import (
	handler2 "UrlShortener/internal/handler"
	"UrlShortener/internal/repository"
	"UrlShortener/internal/service"
	"errors"
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
	cache := repository.NewStubCache()
	persistent := repository.NewStubStorage()
	repoSrv := service.NewRepositoryService(cache, persistent)
	urlSrv := service.NewUrlService(repoSrv)
	h := handler2.NewHandler(tmpl, urlSrv)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", h.Home)
	mux.HandleFunc("POST /shorten", h.Shorten)
	mux.HandleFunc("GET /favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})
	mux.HandleFunc("GET /{code}", h.RedirectByPath)
	mux.HandleFunc("GET /redirect", h.RedirectByQuery)
	srv := http.Server{Addr: ":8080", Handler: mux}
	srvErr := srv.ListenAndServe()

	if errors.Is(srvErr, http.ErrServerClosed) {
		fmt.Println("Shortener server closed, TODO graceful shutdown")
	} else if srvErr != nil {
		fmt.Println("Shortener server get unhandled error", srvErr)
	}
}
