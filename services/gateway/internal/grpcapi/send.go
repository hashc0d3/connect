package grpcapi

import (
	"context"
	"time"

	"github.com/connect/gateway/internal/sms"
	connectv1 "github.com/connect/gateway/gen/connect/v1"
)

func (s *Server) SendSms(ctx context.Context, req *connectv1.SendSmsRequest) (*connectv1.SendSmsResponse, error) {
	id, err := userID(ctx)
	if err != nil {
		return nil, err
	}
	in := sms.SendInput{
		Text: req.GetText(), Sign: req.GetSign(), CallbackURL: req.GetCallbackUrl(),
		ShortLink: req.GetShortLink(), Translit: req.GetTranslit(), GroupIDs: req.GetGroupIds(),
	}
	for _, receiver := range req.GetReceivers() {
		in.Numbers = append(in.Numbers, receiver.GetPhone())
	}
	if req.GetDateSend() > 0 {
		in.DateSend = time.Unix(req.GetDateSend(), 0).UTC()
	}
	if req.GetSmoothUntil() > 0 {
		in.SmoothUntil = time.Unix(req.GetSmoothUntil(), 0).UTC()
	}
	result, err := s.SMS.Send(ctx, id, in)
	if err != nil && (result == nil || len(result.Items) == 0) {
		return nil, grpcErr(err)
	}
	data := make([]*connectv1.Message, 0, len(result.Items))
	for _, item := range result.Items {
		data = append(data, toProto(item))
	}
	return &connectv1.SendSmsResponse{
		Success: true, Message: result.Message, Data: data,
		ExcludedDuplicates: int32(result.Excluded), Cost: result.Cost, HeldForModeration: result.Held,
	}, nil
}
