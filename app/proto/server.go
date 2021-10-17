package proto

import (
	"context"
	"go-fib/app/internal"
)

type GRPCServer struct {
	Api *internal.Server
}

func (s *GRPCServer) mustEmbedUnimplementedReverseServer() {
	panic("implement me")
}

type Res struct {
	pos  uint32
	item uint32
}

func (s *GRPCServer) Fib(ctx context.Context, in *Request) (*Response, error) {

	res := s.Api.CalcFib(in.X, in.Y)
	var items []*Item
	for _, r := range res {
		item := &Item{
			Pos:  r.Position,
			Item: r.Item,
		}
		items = append(items, item)
	}

	return &Response{Items: items}, nil
}
