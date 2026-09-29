package main

import (
	shortenerv1 "ShortenerContract/gen/go/shortener"
	"UrlShortener/internal/adapters/grpcclient"
	"UrlShortener/internal/adapters/grpcserver"
	"UrlShortener/internal/config"
	"UrlShortener/internal/repository"
	"UrlShortener/internal/service"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	fmt.Println("Shortener started")
	appConfig := config.Load()
	fmt.Println("Configuration loaded")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	redisClient := createRedisConn(appConfig.Redis)
	defer redisClient.Close()

	gpPool, err := createPostgresConn(ctx, appConfig.Postgres)
	if err != nil {
		log.Fatal(err)
	}
	defer gpPool.Close()

	rc := createRegistryClient(ctx, appConfig.LbAddr, fmt.Sprintf("%s:%s", appConfig.LocalAddr, appConfig.GrpcPort))
	defer rc.Close()

	cache := repository.NewRedis(redisClient, appConfig.Redis.TTL)
	persistent := repository.NewPostgres(gpPool)
	urlSrv := createUrlService(cache, persistent)

	shortenerSrv := grpcserver.NewShortenerServer(urlSrv)
	grpcSrv := grpc.NewServer()
	shortenerv1.RegisterShortenerServer(grpcSrv, shortenerSrv)
	grpcLis, _ := net.Listen("tcp", fmt.Sprintf(":%s", appConfig.GrpcPort))

	go func() {
		log.Printf("grpc registry server listening on %s", grpcLis.Addr())
		if err := grpcSrv.Serve(grpcLis); err != nil {
			log.Printf("grpc server error: %v\n", err)
		}
	}()

	<-ctx.Done()

	fmt.Println("Shutting down...")

	stopCtx, stopCancel := context.WithTimeout(context.Background(), time.Duration(10) * time.Second)
	defer stopCancel()

	gracefulDone := make(chan struct{})
	stopWg := &sync.WaitGroup{}

	stopWg.Go(func() {
		rc.Deregister(stopCtx)
	})

	stopWg.Go(func() {
		grpcSrv.GracefulStop()
	})

	go func() {
		stopWg.Wait()
		close(gracefulDone)
	}()

	select {
	case <-stopCtx.Done():
	case <-gracefulDone:
	}

	fmt.Println("Shortener stopped")
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

func createRegistryClient(ctx context.Context, lbAddr, localAddr string) *grpcclient.RegistryClient {
	conn, _ := grpc.NewClient(lbAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	rc := grpcclient.NewRegistryClient(conn, localAddr)
	rc.Register(ctx)
	rc.RunHeartbeat(ctx)
	return rc
}