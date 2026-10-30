package sms

import (
	"context"
	"log"

	"github.com/connect/gateway/internal/policy"
	"github.com/connect/gateway/internal/store"
)

func (s *Service) DispatchDue(ctx context.Context) error {
	batch, err := s.Store.ClaimDue(ctx, 100)
	if err != nil {
		return err
	}
	for _, group := range groupByText(batch) {
		phones := make([]string, len(group))
		for i, msg := range group {
			phones[i] = msg.Number
		}
		err := s.Upstream.SendSMS(ctx, phones, group[0].Text)
		for _, msg := range group {
			s.finishDue(ctx, msg, err)
		}
	}
	return nil
}

func groupByText(batch []store.Message) [][]store.Message {
	type key struct{ user, text string }
	grouped := map[key][]store.Message{}
	var order []key
	for _, msg := range batch {
		k := key{msg.UserID, msg.Text}
		if _, ok := grouped[k]; !ok {
			order = append(order, k)
		}
		grouped[k] = append(grouped[k], msg)
	}
	out := make([][]store.Message, 0, len(order))
	for _, k := range order {
		out = append(out, grouped[k])
	}
	return out
}

func (s *Service) finishDue(ctx context.Context, msg store.Message, sendErr error) {
	if sendErr != nil {
		log.Printf("отложенная отправка %d: %v", msg.ID, sendErr)
		_ = s.Store.MarkUpstream(ctx, msg.ID, "failed", policy.StatusReject, policy.ExtendStatus(policy.StatusReject), sendErr.Error(), true)
		_ = s.Store.Refund(ctx, msg.UserID, msg.Cost, msg.ID, "возврат за недоставленную отправку")
		s.callback(msg, policy.StatusReject)
		return
	}
	_ = s.Store.MarkUpstream(ctx, msg.ID, "sent", policy.StatusSent, policy.ExtendStatus(policy.StatusSent), "", true)
	s.callback(msg, policy.StatusSent)
}
