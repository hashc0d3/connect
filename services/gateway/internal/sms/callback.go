package sms

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/connect/gateway/internal/policy"
	"github.com/connect/gateway/internal/store"
)

func (s *Service) callback(msg store.Message, status int) {
	if strings.TrimSpace(msg.CallbackURL) == "" || s.Upstream.HTTP == nil {
		return
	}
	body := fmt.Sprintf(`{"id":%d,"status":%d,"extendStatus":"%s"}`, msg.ID, status, policy.ExtendStatus(status))
	url := msg.CallbackURL
	client := s.Upstream.HTTP
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBufferString(body))
		if err != nil {
			return
		}
		req.Header.Set("Content-Type", "application/json")
		res, err := client.Do(req)
		if err != nil {
			return
		}
		res.Body.Close()
	}()
}
