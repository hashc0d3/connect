package grpcapi

import (
	"context"

	connectv1 "github.com/connect/gateway/gen/connect/v1"
)

func (s *Server) GetStatus(ctx context.Context, req *connectv1.GetStatusRequest) (*connectv1.Message, error) {
	id, err := userID(ctx)
	if err != nil {
		return nil, err
	}
	item, err := s.SMS.Status(ctx, id, req.GetId())
	if err != nil {
		return nil, grpcErr(err)
	}
	return toProto(*item), nil
}

func (s *Server) ListMessages(ctx context.Context, req *connectv1.ListMessagesRequest) (*connectv1.ListMessagesResponse, error) {
	id, err := userID(ctx)
	if err != nil {
		return nil, err
	}
	items, total, err := s.SMS.List(ctx, id, req.GetNumber(), req.GetText(), int(req.GetPage()))
	if err != nil {
		return nil, grpcErr(err)
	}
	data := make([]*connectv1.Message, 0, len(items))
	for _, item := range items {
		data = append(data, toProto(item))
	}
	return &connectv1.ListMessagesResponse{Success: true, Data: data, Page: req.GetPage(), Total: int32(total)}, nil
}

func (s *Server) GetBalance(ctx context.Context, _ *connectv1.Empty) (*connectv1.BalanceResponse, error) {
	id, err := userID(ctx)
	if err != nil {
		return nil, err
	}
	balance, err := s.SMS.Balance(ctx, id)
	if err != nil {
		return nil, grpcErr(err)
	}
	return &connectv1.BalanceResponse{Success: true, Balance: balance}, nil
}
