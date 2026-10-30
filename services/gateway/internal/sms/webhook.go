package sms

import (
	"context"

	"github.com/connect/gateway/internal/policy"
)

type WebhookStatus struct {
	ClientMessageID string
	UpstreamID      string
	State           string
	Code            string
}

func (s *Service) ApplyWebhook(ctx context.Context, in WebhookStatus) error {
	if status, upstream, ok := policy.EasySMSStatus(in.State); ok && in.ClientMessageID != "" {
		if _, err := s.Store.ApplyDelivery(ctx, in.ClientMessageID, in.UpstreamID, status, policy.ExtendStatus(status), upstream); err != nil {
			return err
		}
	}
	if in.Code != "" && in.ClientMessageID != "" {
		if _, err := s.Store.SaveOTPCode(ctx, in.ClientMessageID, hashCode(in.Code)); err != nil {
			return err
		}
	}
	return nil
}
