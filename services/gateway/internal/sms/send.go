package sms

import (
	"context"
	"strings"
	"time"

	"github.com/connect/gateway/internal/policy"
	"github.com/shopspring/decimal"
)

func (s *Service) Send(ctx context.Context, userID string, in SendInput) (*SendResult, error) {
	if err := s.limit(userID, 60); err != nil {
		return nil, err
	}
	in.Text = strings.TrimSpace(in.Text)
	if in.Text == "" {
		return nil, Fail(422, "Текст сообщения пустой")
	}
	if err := policy.ValidateCallback(in.CallbackURL); err != nil {
		return nil, Fail(422, err.Error())
	}
	rules, err := s.Store.Rules(ctx)
	if err != nil {
		return nil, err
	}
	tariffs, err := s.Store.Tariffs(ctx)
	if err != nil {
		return nil, err
	}
	signName, channel, err := s.resolveSign(ctx, userID, in.Sign, rules)
	if err != nil {
		return nil, err
	}
	numbers, excluded, blacklisted, err := s.recipients(ctx, userID, in, rules)
	if err != nil {
		return nil, err
	}
	drafts, total, err := s.buildDrafts(ctx, userID, in, rules, tariffs, numbers, signName, channel)
	if err != nil {
		return nil, err
	}
	return s.route(ctx, userID, in, signName, channel, drafts, total, excluded, blacklisted)
}

func (s *Service) route(ctx context.Context, userID string, in SendInput, signName, channel string, drafts []draft, total decimal.Decimal, excluded, blacklisted int) (*SendResult, error) {
	approved, err := s.Store.HasApprovedProfile(ctx, userID)
	if err != nil {
		return nil, err
	}
	rules, err := s.Store.Rules(ctx)
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	start := in.DateSend
	if start.IsZero() {
		start = now
	}
	hold := rules.RequireContract && !approved
	deferred := !hold && (start.After(now.Add(2*time.Second)) || !in.SmoothUntil.IsZero())

	var scheduled, smooth *time.Time
	if deferred {
		scheduled = &start
		if !in.SmoothUntil.IsZero() {
			smooth = &in.SmoothUntil
		}
	}
	if hold {
		items, err := s.persist(ctx, persistRequest{
			UserID: userID, Original: in.Text, Text: in.Text, Sign: signName,
			CampaignStatus: "moderation", Channel: channel, Drafts: drafts, Excluded: excluded,
			Status: policy.StatusModeration, Upstream: "held", Callback: in.CallbackURL,
		})
		if err != nil {
			return nil, err
		}
		return &SendResult{
			Message: "Сообщения приняты и ждут модерации Connect", Items: items,
			Excluded: excluded, Blacklisted: blacklisted, Cost: "0.0000", Held: true,
		}, nil
	}
	if deferred {
		items, err := s.persist(ctx, persistRequest{
			UserID: userID, Original: in.Text, Text: in.Text, Sign: signName,
			CampaignStatus: "queue", Channel: channel, Drafts: drafts, Excluded: excluded,
			Charge: total, Scheduled: scheduled, Smooth: smooth, Status: policy.StatusQueue,
			Upstream: "pending", Debit: true, Callback: in.CallbackURL,
		})
		if err != nil {
			return nil, err
		}
		return &SendResult{
			Message: "Сообщения поставлены в очередь", Items: items,
			Excluded: excluded, Blacklisted: blacklisted, Cost: total.StringFixed(4),
		}, nil
	}
	return s.deliverNow(ctx, userID, in, signName, channel, drafts, excluded, blacklisted)
}
