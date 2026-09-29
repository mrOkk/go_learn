package grpcclient

import (
	"fmt"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Pool struct {
	mu    sync.Mutex
	conns map[string]*grpc.ClientConn
}

func NewConnPool() *Pool {
	return &Pool{conns: make(map[string]*grpc.ClientConn)}
}

func (p *Pool) get(address string) (*grpc.ClientConn, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if conn, ok := p.conns[address]; ok {
		return conn, nil
	}

	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s: %w", address, err)
	}
	p.conns[address] = conn
	return conn, nil
}

func (p *Pool) Close(address string) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	conn, ok := p.conns[address]
	if !ok {
		return nil
	}
	delete(p.conns, address)
	return conn.Close()
}
