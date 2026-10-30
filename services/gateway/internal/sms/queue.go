package sms

import (
	"context"
	"fmt"
	"log"

	"github.com/connect/gateway/internal/policy"
	"github.com/connect/gateway/internal/upstream"
)

func (s *Service) DispatchDue(ctx context.Context) error {
	batch, err := s.Store.ClaimDue(ctx, 100)
	if err != nil {
		return err
	}
	for _, msg := range batch {
		clientID := fmt.Sprintf("msg_%d", msg.ID)
		accepted, err := s.Upstream.SendMessage(ctx, upstream.Outbound{
			ClientMessageID: clientID,
			Messenger:       "sms",
			Originator:      msg.From,
			Recipient:       msg.Number,
			Text:            msg.Text,
		})
		if err != nil {
			log.Printf("отложенная отправка %d: %v", msg.ID, err)
			_ = s.Store.MarkUpstream(ctx, msg.ID, "failed", policy.StatusReject, policy.ExtendStatus(policy.StatusReject), err.Error(), true)
			_ = s.Store.Refund(ctx, msg.UserID, msg.Cost, msg.ID, "возврат за недоставленную отправку")
			s.callback(msg, policy.StatusReject)
			continue
		}
		_ = s.Store.MarkAccepted(ctx, msg.ID, accepted.ClientMessageID, accepted.MessageID)
		s.callback(msg, policy.StatusSent)
	}
	return nil
}
