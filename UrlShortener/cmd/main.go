package main

import (
	"UrlShortener/internal/config"
	"UrlShortener/internal/handler"
	"UrlShortener/internal/repository"
	"UrlShortener/internal/service"
	"context"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func main() {
	fmt.Println("Shortener started")
	ctx := context.Background()

	appConfig := config.Load()
	redisClient := createRedisConn(appConfig.Redis)
	defer redisClient.Close()

	gpPool, err := createPostgresConn(ctx, appConfig.Postgres)
	if err != nil {
		log.Fatal(err)
	}
	defer gpPool.Close()

	cache := repository.NewRedis(redisClient, appConfig.Redis.TTL)
	persistent := repository.NewPostgres(gpPool)
	urlSrv := createUrlService(cache, persistent)

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Fatal(err)
	}

	h := handler.NewHandler(tmpl, urlSrv)
	mux := http.NewServeMux()
	registerHandlers(mux, h)

	srv := http.Server{Addr: fmt.Sprintf(":%s", appConfig.Port), Handler: mux}
	err = srv.ListenAndServe()
	defer srv.Close()

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Println("Shortener server closed, TODO graceful shutdown")
	} else if err != nil {
		fmt.Println("Shortener server get unhandled error", err)
	}
}

func createRedisConn(config config.RedisConfig) *redis.Client {
	return redis.NewClient(config.GetOpts())
}

func createPostgresConn(ctx context.Context, config config.PostgresConfig) (*pgxpool.Pool, error) {
	return pgxpool.New(ctx, config.GetConnString())
}

func createUrlService(cache service.CacheStorage, persistent service.UrlStorage) *service.UrlService {
	repoSrv := service.NewRepositoryService(cache, persistent)
	return service.NewUrlService(repoSrv)
}

func registerHandlers(mux *http.ServeMux, h *handler.Handler) {
	mux.HandleFunc("GET /", h.Home)
	mux.HandleFunc("POST /shorten", h.Shorten)
	mux.HandleFunc("GET /favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})
	mux.HandleFunc("GET /{code}", h.RedirectByPath)
	mux.HandleFunc("GET /redirect", h.RedirectByQuery)
}