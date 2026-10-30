package grpcapi

import (
	"errors"

	"github.com/connect/gateway/internal/sms"
	connectv1 "github.com/connect/gateway/gen/connect/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func grpcErr(err error) error {
	var api *sms.APIError
	if errors.As(err, &api) {
		switch api.HTTP {
		case 401:
			return status.Error(codes.Unauthenticated, api.Message)
		case 402:
			return status.Error(codes.FailedPrecondition, api.Message)
		case 404:
			return status.Error(codes.NotFound, api.Message)
		case 422, 400:
			return status.Error(codes.InvalidArgument, api.Message)
		case 429:
			return status.Error(codes.ResourceExhausted, api.Message)
		case 502:
			return status.Error(codes.Unavailable, api.Message)
		default:
			return status.Error(codes.Internal, api.Message)
		}
	}
	return status.Error(codes.Internal, "Внутренняя ошибка")
}

func toProto(item sms.Item) *connectv1.Message {
	msg := &connectv1.Message{
		Id: item.ID, From: item.From, Number: item.Number, Text: item.Text,
		Status: int32(item.Status), ExtendStatus: item.ExtendStatus, Channel: item.Channel,
		Cost: item.Cost, DateCreate: item.DateCreate.Unix(), DateSend: item.DateSend.Unix(),
	}
	if item.DateAnswer != nil {
		msg.DateAnswer = item.DateAnswer.Unix()
	}
	return msg
}
