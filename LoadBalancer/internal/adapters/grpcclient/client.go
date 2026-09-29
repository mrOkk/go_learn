package grpcclient

import (
	shortenerv1 "ShortenerContract/gen/go/shortener"
	"context"
)

type Caller struct {
	pool *Pool
}

func NewCaller(pool *Pool) *Caller {
	return &Caller{pool: pool}
}

func (c *Caller) Shorten(ctx context.Context, address, url string) (string, error) {
	conn, err := c.pool.get(address)
	if err != nil {
		return "", err
	}
	client := shortenerv1.NewShortenerClient(conn)
	resp, err := client.Shorten(ctx, &shortenerv1.ShortenRequest{Url: url})
	if err != nil {
		return "", err
	}
	return resp.Code, nil
}

func (c *Caller) Resolve(ctx context.Context, address, code string) (string, error) {
	conn, err := c.pool.get(address)
	if err != nil {
		return "", err
	}
	client := shortenerv1.NewShortenerClient(conn)
	resp, err := client.Resolve(ctx, &shortenerv1.ResolveRequest{Code: code})
	if err != nil {
		return "", err
	}
	return resp.Url, nil
}
