package grpcclient

import (
	shortenerv1 "ShortenerContract/gen/go/shortener"
	"context"
	"log"
	"sync"
	"time"

	"google.golang.org/grpc"
)

type RegistryClient struct {
	id      string
	addr    string
	cc      *grpc.ClientConn
	ttl     time.Duration
	closing chan struct{}
	wg      sync.WaitGroup
}

func NewRegistryClient(cc *grpc.ClientConn, addr string) *RegistryClient {
	return &RegistryClient{
		cc:   cc,
		addr: addr,
	}
}

func (rc *RegistryClient) Register(ctx context.Context) {
	log.Printf("Register service: %v \n", rc.addr)
	client := shortenerv1.NewRegistryClient(rc.cc)
	resp, err := client.Register(ctx, &shortenerv1.RegisterRequest{Address: rc.addr})
	if err != nil {
		log.Fatalf("Registration not successful: %v", err)
	}
	rc.id = resp.InstanceId
	rc.ttl = time.Duration(resp.TtlSeconds) * time.Second
	rc.closing = make(chan struct{})
}

func (rc *RegistryClient) Deregister(ctx context.Context) {
	rc.wg.Add(1)
	defer rc.wg.Done()
	client := shortenerv1.NewRegistryClient(rc.cc)
	_, err := client.Deregister(ctx, &shortenerv1.DeregisterRequest{InstanceId: rc.id})
	if err != nil {
		log.Printf("Deregister not successful: %v", err)
	}
	rc.id = ""
}

func (rc *RegistryClient) Close() {
	close(rc.closing)
	rc.wg.Wait()
	rc.cc.Close()
}

func (rc *RegistryClient) RunHeartbeat(ctx context.Context) {
	rc.wg.Add(1)
	go rc.Heartbeat(ctx)
}

func (rc *RegistryClient) Heartbeat(ctx context.Context) {
	defer rc.wg.Done()

	ticker := time.NewTicker(rc.ttl / 3)
	defer ticker.Stop()
	reg := shortenerv1.NewRegistryClient(rc.cc)

	for {
		select {
		case <-ctx.Done():
			return
		case <-rc.closing:
			return
		case <-ticker.C:
			_, err := reg.Heartbeat(ctx, &shortenerv1.HeartbeatRequest{InstanceId: rc.id})
			if err != nil {
				// DO something
				log.Printf("heartbeat failed: %v", err)
			}
		}
	}
}
