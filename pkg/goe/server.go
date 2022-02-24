package goe

import (
	context "context"
	"fmt"
	"log"
)

// NewRandomerServer returns a gRPC RandomerServer implementation.
func NewRandomerServer(logger *log.Logger, store Store) RandomerServer {
	return &randomerServer{
		logger: logger,
		store:  store,
	}
}

type randomerServer struct {
	UnimplementedRandomerServer

	logger *log.Logger
	store  Store
}

func (r *randomerServer) GenerateRandom(ctx context.Context, req *RandomRequest) (*RandomReply, error) {
	r.logger.Printf(
		"received timestamp: %d, id: %x, data: %x",
		req.CreatedAt.AsTime().UnixNano(), req.Id, req.Data)

	if err := r.store.Add(ctx, req.CreatedAt.AsTime(), req.Id, req.Data); err != nil {
		return nil, fmt.Errorf("store: %w", err)
	}
	return &RandomReply{}, nil
}
