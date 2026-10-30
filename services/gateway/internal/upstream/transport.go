package upstream

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type envelope struct {
	Status  bool                `json:"status"`
	Message string              `json:"message"`
	Data    json.RawMessage     `json:"data"`
	Errors  map[string][]string `json:"errors"`
}

type ValidationError struct {
	Message string
	Fields  map[string][]string
}

func (e *ValidationError) Error() string { return e.Message }

func (c *Client) post(ctx context.Context, path string, payload any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+path, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	res, err := c.HTTP.Do(req)
	if err != nil {
		return fmt.Errorf("EasySMS недоступен: %w", err)
	}
	defer res.Body.Close()
	raw, _ := io.ReadAll(io.LimitReader(res.Body, 1<<20))

	var env envelope
	if json.Unmarshal(raw, &env) != nil {
		return fmt.Errorf("EasySMS вернул непонятный ответ (%d)", res.StatusCode)
	}
	if res.StatusCode == http.StatusUnauthorized || strings.Contains(strings.ToLower(env.Message), "unauthenticated") {
		return fmt.Errorf("ключ EasySMS отклонён")
	}
	if res.StatusCode >= 500 {
		return fmt.Errorf("EasySMS временно недоступен (%d)", res.StatusCode)
	}
	if env.Status {
		return nil
	}
	msg := env.Message
	if msg == "" {
		msg = "EasySMS отклонил запрос"
	}
	if len(env.Errors) > 0 {
		return &ValidationError{Message: msg, Fields: env.Errors}
	}
	return fmt.Errorf("%s", msg)
}
