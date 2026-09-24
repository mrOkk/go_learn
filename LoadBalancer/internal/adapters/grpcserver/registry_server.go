package grpcserver

import (
	"LoadBalancer/internal/app"
	shortenerv1 "ShortenerContract/gen/go/shortener"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RegistryServer struct {
	shortenerv1.UnimplementedRegistryServer
	registry app.Registry
}

func NewServer(registry app.Registry) *RegistryServer {
	return &RegistryServer{registry: registry}
}

func (r* RegistryServer) Run() {

}

func (r* RegistryServer) Register(ctx context.Context, request *shortenerv1.RegisterRequest) (*shortenerv1.RegisterResponse, error) {
	id, ttl := r.registry.Register(request.Address)
	return &shortenerv1.RegisterResponse{
		InstanceId: id,
		TtlSeconds: int64(ttl.Seconds()),
		}, nil
}

func (r* RegistryServer) Heartbeat(ctx context.Context, request *shortenerv1.HeartbeatRequest) (*shortenerv1.HeartbeatResponse, error){
	err := r.registry.Heartbeat(request.InstanceId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "instance %q not found %v", request.InstanceId, err)
	}
	return &shortenerv1.HeartbeatResponse{}, nil
}

func (r* RegistryServer) Deregister(ctx context.Context, request *shortenerv1.DeregisterRequest) (*shortenerv1.DeregisterResponse, error) {
	r.registry.Deregister(request.InstanceId)
	return &shortenerv1.DeregisterResponse{}, nil
}