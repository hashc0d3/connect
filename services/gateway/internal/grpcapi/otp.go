package grpcapi

import (
	"context"

	connectv1 "github.com/connect/gateway/gen/connect/v1"
)

func (s *Server) CreateOtp(ctx context.Context, req *connectv1.CreateOtpRequest) (*connectv1.StatusResponse, error) {
	id, err := userID(ctx)
	if err != nil {
		return nil, err
	}
	message, err := s.SMS.CreateOTP(ctx, id, req.GetPhone())
	if err != nil {
		return nil, grpcErr(err)
	}
	return &connectv1.StatusResponse{Success: true, Message: message}, nil
}

func (s *Server) ValidateOtp(ctx context.Context, req *connectv1.ValidateOtpRequest) (*connectv1.StatusResponse, error) {
	id, err := userID(ctx)
	if err != nil {
		return nil, err
	}
	message, err := s.SMS.ValidateOTP(ctx, id, req.GetPhone(), req.GetCode())
	if err != nil {
		return nil, grpcErr(err)
	}
	return &connectv1.StatusResponse{Success: true, Message: message}, nil
}

func (s *Server) AuthCheck(ctx context.Context, _ *connectv1.Empty) (*connectv1.StatusResponse, error) {
	if _, err := userID(ctx); err != nil {
		return nil, err
	}
	return &connectv1.StatusResponse{Success: true, Message: "Successful authorization."}, nil
}
