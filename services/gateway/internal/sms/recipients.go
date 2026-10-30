package sms

import (
	"context"
	"strings"

	"github.com/connect/gateway/internal/policy"
	"github.com/connect/gateway/internal/store"
)

func (s *Service) recipients(ctx context.Context, userID string, in SendInput, rules store.Rules) ([]string, int, int, error) {
	numbers, excluded, err := s.collectNumbers(ctx, userID, in)
	if err != nil {
		return nil, 0, 0, err
	}
	blacklisted := 0
	if rules.BlockBlacklist {
		blocked, err := s.Store.Blacklist(ctx, userID)
		if err != nil {
			return nil, 0, 0, err
		}
		filtered := numbers[:0]
		for _, number := range numbers {
			if _, ok := blocked[number]; ok {
				blacklisted++
				continue
			}
			filtered = append(filtered, number)
		}
		numbers = filtered
	}
	if len(numbers) == 0 {
		return nil, 0, 0, Fail(422, "Нет получателей после проверки списка и чёрного списка")
	}
	if len(numbers) > rules.MaxReceivers {
		return nil, 0, 0, Fail(422, "Получателей больше допустимого лимита")
	}
	return numbers, excluded, blacklisted, nil
}

func (s *Service) collectNumbers(ctx context.Context, userID string, in SendInput) ([]string, int, error) {
	seen := map[string]struct{}{}
	var out []string
	rawCount := 0
	add := func(raw string) error {
		rawCount++
		number, err := policy.NormalizePhone(raw)
		if err != nil {
			return Fail(422, raw+": "+err.Error())
		}
		if _, ok := seen[number]; ok {
			return nil
		}
		seen[number] = struct{}{}
		out = append(out, number)
		return nil
	}
	for _, number := range in.Numbers {
		if strings.TrimSpace(number) == "" {
			continue
		}
		if err := add(number); err != nil {
			return nil, 0, err
		}
	}
	if len(in.GroupIDs) == 0 {
		return out, rawCount - len(out), nil
	}
	contacts, err := s.Store.ContactsByGroups(ctx, userID, in.GroupIDs)
	if err != nil {
		return nil, 0, err
	}
	for _, contact := range contacts {
		if err := add(contact.Number); err != nil {
			return nil, 0, err
		}
	}
	excluded := rawCount - len(out)
	if excluded < 0 {
		excluded = 0
	}
	return out, excluded, nil
}
