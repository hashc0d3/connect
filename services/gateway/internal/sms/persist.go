package sms

import (
	"context"
	"errors"
	"time"

	"github.com/connect/gateway/internal/policy"
	"github.com/connect/gateway/internal/store"
	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"
)

type persistRequest struct {
	UserID         string
	Original       string
	Text           string
	Sign           string
	CampaignStatus string
	Channel        string
	Drafts         []draft
	Excluded       int
	Charge         decimal.Decimal
	Scheduled      *time.Time
	Smooth         *time.Time
	Status         int
	Upstream       string
	Debit          bool
	Callback       string
}

func (s *Service) persist(ctx context.Context, req persistRequest) ([]Item, error) {
	now := time.Now().UTC()
	var items []Item
	err := s.Store.Tx(ctx, func(tx pgx.Tx) error {
		if req.Debit {
			if err := s.Store.DebitTx(ctx, tx, req.UserID, req.Charge, nil, "рассылка"); err != nil {
				if err.Error() == "недостаточно средств" {
					return Fail(402, "Недостаточно средств на балансе")
				}
				return err
			}
		}
		campaignID, err := s.Store.InsertCampaign(ctx, tx, req.UserID, req.Original, req.Text, req.Sign, req.CampaignStatus, len(req.Drafts), req.Excluded, req.Charge, req.Scheduled, req.Smooth)
		if err != nil {
			return err
		}
		for _, item := range req.Drafts {
			id, err := s.Store.InsertMessage(ctx, tx, store.NewMessage{
				UserID: req.UserID, CampaignID: campaignID, From: req.Sign, Number: item.Number, Text: item.Text,
				Status: req.Status, Extend: policy.ExtendStatus(req.Status), Channel: req.Channel,
				Cost: item.Cost, Parts: item.Parts, DateCreate: now, DateSend: item.When,
				Upstream: req.Upstream, CallbackURL: req.Callback,
				ClientID: item.ClientMessageID, UpstreamID: item.UpstreamID,
			})
			if err != nil {
				return err
			}
			items = append(items, Item{
				ID: id, From: req.Sign, Number: item.Number, Text: item.Text, Status: req.Status,
				ExtendStatus: policy.ExtendStatus(req.Status), Channel: req.Channel,
				Cost: item.Cost.StringFixed(4), DateCreate: now, DateSend: item.When,
			})
		}
		return nil
	})
	var api *APIError
	if errors.As(err, &api) {
		return nil, api
	}
	return items, err
}
