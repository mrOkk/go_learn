package app

import (
	shortenerv1 "ShortenerContract/gen/go/shortener"
	"context"
	"time"

	"google.golang.org/grpc"
)

type RegistryServer interface {
	Register(ctx context.Context, request *shortenerv1.RegisterRequest) (*shortenerv1.RegisterResponse, error)
	Heartbeat(ctx context.Context, request *shortenerv1.HeartbeatRequest) (*shortenerv1.HeartbeatResponse, error)
	Deregister(ctx context.Context, request *shortenerv1.DeregisterRequest) (*shortenerv1.DeregisterResponse, error)
	mustEmbedUnimplementedRegistryServer()
}

type ShortenerClient interface {
	Shorten(ctx context.Context, request *shortenerv1.ShortenRequest, opts ...grpc.CallOption) (*shortenerv1.ShortenResponse, error)
	Resolve(ctx context.Context, request *shortenerv1.ResolveRequest, opts ...grpc.CallOption) (*shortenerv1.ResolveResponse, error)
}

type Registry interface {
	Register(address string) (string, time.Duration)
	Heartbeat(id string) error
	Deregister(id string)
}

type BackendCaller interface {
	Shorten(ctx context.Context, address, url string) (string, error)
	Resolve(ctx context.Context, address, code string) (string, error)
}