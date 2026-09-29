package main

import (
	"LoadBalancer/internal/adapters/grpcclient"
	"LoadBalancer/internal/adapters/grpcserver"
	"LoadBalancer/internal/adapters/httpfront"
	"LoadBalancer/internal/app"
	"LoadBalancer/internal/config"
	"LoadBalancer/internal/registry"
	shortenerv1 "ShortenerContract/gen/go/shortener"
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	reg := registry.NewRegistry(cfg.GRPC.TTL)
	go reg.SweepJob(ctx, cfg.GRPC.TTL/3)

	picker := app.NewRoundRobin(reg)

	pool := grpcclient.NewConnPool()
	caller := grpcclient.NewCaller(pool)

	handler := httpfront.NewHTTPHandler("templates/index.html", picker, caller)
	mux := httpfront.NewRouter(handler)
	httpSrv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: mux,
	}

	registrySrv := grpcserver.NewServer(reg)
	grpcSrv := grpc.NewServer()
	shortenerv1.RegisterRegistryServer(grpcSrv, registrySrv)

	grpcLis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPC.Port))
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		log.Printf("Listening on %v\n", httpSrv.Addr)
		if err := httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("Http server error: %v\n", err)
		}
	}()

	go func() {
		log.Printf("grpc registry server listening on %s", grpcLis.Addr())
		if err := grpcSrv.Serve(grpcLis); err != nil {
			log.Printf("grpc server error: %v\n", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpSrv.Shutdown(shutdownCtx); err != nil {
		log.Printf("http shutdown error: %v\n", err)
	}
	grpcSrv.GracefulStop()

	log.Println("Shutdown complete")
}
