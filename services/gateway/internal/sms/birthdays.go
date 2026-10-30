package sms

import (
	"context"
	"log"
	"strings"
)

func (s *Service) RunBirthdays(ctx context.Context) {
	rules, err := s.Store.DueAutoRules(ctx)
	if err != nil {
		log.Printf("авторассылка: %v", err)
		return
	}
	for _, rule := range rules {
		s.runBirthday(ctx, rule.UserID, rule.TemplateID, rule.SignName)
	}
}

func (s *Service) runBirthday(ctx context.Context, userID string, templateID int, signName string) {
	_ = s.Store.MarkAutoRan(ctx, userID)
	body, err := s.Store.TemplateBody(ctx, userID, templateID)
	if err != nil || strings.TrimSpace(body) == "" {
		return
	}
	numbers, err := s.Store.BirthdayNumbers(ctx, userID)
	if err != nil || len(numbers) == 0 {
		return
	}
	if _, err := s.Send(ctx, userID, SendInput{Numbers: numbers, Text: body, Sign: signName}); err != nil {
		log.Printf("авторассылка %s: %v", userID, err)
	}
}
