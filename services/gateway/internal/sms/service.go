package sms

import (
	"sync"
	"time"

	"github.com/connect/gateway/internal/store"
	"github.com/connect/gateway/internal/upstream"
)

type Service struct {
	Store    *store.Store
	Upstream *upstream.Client
	Links    string
	Reveal   bool

	mu      sync.Mutex
	buckets map[string]bucket
}

type bucket struct {
	start time.Time
	count int
}

type APIError struct {
	HTTP    int
	Message string
	Fields  map[string][]string
}

func (e *APIError) Error() string { return e.Message }

func Fail(code int, msg string) *APIError {
	return &APIError{HTTP: code, Message: msg}
}

func (s *Service) limit(userID string, max int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.buckets == nil {
		s.buckets = map[string]bucket{}
	}
	now := time.Now()
	current := s.buckets[userID]
	if now.Sub(current.start) > time.Minute {
		current = bucket{start: now}
	}
	current.count++
	s.buckets[userID] = current
	if current.count > max {
		return Fail(429, "Слишком много запросов, повторите через минуту")
	}
	return nil
}
