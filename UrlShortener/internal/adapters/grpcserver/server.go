package grpcserver

import (
	shortenerv1 "ShortenerContract/gen/go/shortener"
	"UrlShortener/internal/service"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ShortenerServer struct {
	shortenerv1.UnimplementedShortenerServer
	urlSrv *service.UrlService
}

func NewShortenerServer(urlSrv *service.UrlService) *ShortenerServer {
	return &ShortenerServer{urlSrv: urlSrv}
}

func (s *ShortenerServer) Shorten(ctx context.Context, request *shortenerv1.ShortenRequest) (*shortenerv1.ShortenResponse, error) {
	code, err := s.urlSrv.Put(ctx, request.Url)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}
	return &shortenerv1.ShortenResponse{Code: code}, nil
}

func (s *ShortenerServer) Resolve(ctx context.Context, request *shortenerv1.ResolveRequest) (*shortenerv1.ResolveResponse, error) {
	url, err := s.urlSrv.Get(ctx, request.Code)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "%v", err)
	}
	return &shortenerv1.ResolveResponse{Url: url}, nil
}