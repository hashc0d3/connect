package sms

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/connect/gateway/internal/policy"
	"github.com/connect/gateway/internal/upstream"
	"github.com/shopspring/decimal"
)

func (s *Service) deliverNow(ctx context.Context, userID string, in SendInput, signName, channel string, drafts []draft, excluded, blacklisted int) (*SendResult, error) {
	failed := s.pushUpstream(ctx, drafts)
	okDrafts, badDrafts, charged, reason := splitDrafts(drafts, failed)

	var items []Item
	if len(okDrafts) > 0 {
		saved, err := s.persist(ctx, persistRequest{
			UserID: userID, Original: in.Text, Text: in.Text, Sign: signName,
			CampaignStatus: "sent", Channel: channel, Drafts: okDrafts, Excluded: excluded,
			Charge: charged, Status: policy.StatusSent, Upstream: "sent", Debit: true, Callback: in.CallbackURL,
		})
		if err != nil {
			return nil, err
		}
		items = append(items, saved...)
	}
	if len(badDrafts) > 0 {
		saved, err := s.persist(ctx, persistRequest{
			UserID: userID, Original: in.Text, Text: in.Text, Sign: signName,
			CampaignStatus: "reject", Channel: channel, Drafts: badDrafts,
			Status: policy.StatusReject, Upstream: "failed", Callback: in.CallbackURL,
		})
		if err != nil {
			return nil, err
		}
		items = append(items, saved...)
	}
	if len(okDrafts) == 0 {
		if reason == "" {
			reason = "EasySMS отклонил отправку"
		}
		return &SendResult{Message: reason, Items: items, Excluded: excluded, Blacklisted: blacklisted, Cost: "0.0000"}, Fail(502, reason)
	}
	note := "Messages sent."
	if len(badDrafts) > 0 {
		note = "Часть сообщений отправлена, часть отклонена поставщиком"
	}
	if blacklisted > 0 {
		note += fmt.Sprintf(" Исключено из чёрного списка: %d.", blacklisted)
	}
	return &SendResult{
		Message: note, Items: items, Excluded: excluded, Blacklisted: blacklisted, Cost: charged.StringFixed(4),
	}, nil
}

func (s *Service) pushUpstream(ctx context.Context, drafts []draft) map[int]string {
	grouped := map[string][]int{}
	var order []string
	for i, item := range drafts {
		if _, ok := grouped[item.Text]; !ok {
			order = append(order, item.Text)
		}
		grouped[item.Text] = append(grouped[item.Text], i)
	}
	failed := map[int]string{}
	for _, text := range order {
		indexes := grouped[text]
		phones := make([]string, len(indexes))
		for i, idx := range indexes {
			phones[i] = drafts[idx].Number
		}
		for _, chunk := range chunks(phones, 1000) {
			if err := s.Upstream.SendSMS(ctx, chunk, text); err != nil {
				msg := upstreamMessage(err)
				present := map[string]struct{}{}
				for _, phone := range chunk {
					present[phone] = struct{}{}
				}
				for _, idx := range indexes {
					if _, ok := present[drafts[idx].Number]; ok {
						failed[idx] = msg
					}
				}
			}
		}
	}
	return failed
}

func splitDrafts(drafts []draft, failed map[int]string) (ok, bad []draft, charged decimal.Decimal, reason string) {
	for i, item := range drafts {
		if msg, rejected := failed[i]; rejected {
			bad = append(bad, item)
			if reason == "" {
				reason = msg
			}
			continue
		}
		ok = append(ok, item)
		charged = charged.Add(item.Cost)
	}
	return ok, bad, charged, reason
}

func upstreamMessage(err error) string {
	var validation *upstream.ValidationError
	if errors.As(err, &validation) {
		return validation.Message
	}
	log.Printf("easysms send: %v", err)
	if err.Error() == "" {
		return "поставщик отклонил сообщение"
	}
	return err.Error()
}
