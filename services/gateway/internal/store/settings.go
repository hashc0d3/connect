package store

import (
	"context"
	"fmt"
	"strings"

	"github.com/shopspring/decimal"
)

type Rules struct {
	MinText         int
	MaxReceivers    int
	MarkupPercent   decimal.Decimal
	RequireSign     bool
	RequireContract bool
	BlockBlacklist  bool
	PrefixSign      bool
	PromoPercent    decimal.Decimal
	PromoMin        decimal.Decimal
	PromoCode       string
}

func (s *Store) Rules(ctx context.Context) (Rules, error) {
	rows, err := s.pool.Query(ctx, `SELECT key, value FROM platform_settings`)
	if err != nil {
		return Rules{}, err
	}
	defer rows.Close()
	raw := map[string]string{}
	for rows.Next() {
		var k, v string
		if err := rows.Scan(&k, &v); err != nil {
			return Rules{}, err
		}
		raw[k] = v
	}
	r := Rules{
		MinText:         atoiDefault(raw["min_text_length"], 4),
		MaxReceivers:    atoiDefault(raw["max_receivers"], 1000),
		MarkupPercent:   decDefault(raw["markup_percent"], "20"),
		RequireSign:     raw["require_approved_sign"] != "false",
		RequireContract: raw["require_approved_contract"] == "true",
		BlockBlacklist:  raw["block_blacklist"] != "false",
		PrefixSign:      raw["prefix_sign_in_text"] == "true",
		PromoPercent:    decDefault(raw["promo_percent"], "10"),
		PromoMin:        decDefault(raw["promo_min_sum"], "3000"),
		PromoCode:       strings.TrimSpace(raw["promo_code"]),
	}
	if r.MaxReceivers <= 0 || r.MaxReceivers > 1000 {
		r.MaxReceivers = 1000
	}
	return r, rows.Err()
}

func (s *Store) Tariffs(ctx context.Context) (map[string]decimal.Decimal, error) {
	rows, err := s.pool.Query(ctx, `SELECT channel, operator, price::text FROM tariffs`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := map[string]decimal.Decimal{}
	for rows.Next() {
		var channel, operator, price string
		if err := rows.Scan(&channel, &operator, &price); err != nil {
			return nil, err
		}
		out[channel+"|"+operator], err = decimal.NewFromString(price)
		if err != nil {
			return nil, err
		}
	}
	return out, rows.Err()
}

func Price(tariffs map[string]decimal.Decimal, channel, operator string) decimal.Decimal {
	if v, ok := tariffs[channel+"|"+operator]; ok {
		return v
	}
	if v, ok := tariffs[channel+"|OTHER"]; ok {
		return v
	}
	if v, ok := tariffs["PAY_SIGN|OTHER"]; ok {
		return v
	}
	return decimal.RequireFromString("3.00")
}

func atoiDefault(v string, def int) int {
	if v == "" {
		return def
	}
	var n int
	_, err := fmt.Sscan(v, &n)
	if err != nil {
		return def
	}
	return n
}

func decDefault(v, def string) decimal.Decimal {
	if strings.TrimSpace(v) == "" {
		v = def
	}
	d, err := decimal.NewFromString(v)
	if err != nil {
		d, _ = decimal.NewFromString(def)
	}
	return d
}
