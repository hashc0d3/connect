package upstream

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type apiError struct {
	Code    string `json:"error_code"`
	Message string `json:"error_message"`
}

type ValidationError struct {
	Message string
	Fields  map[string][]string
}

func (e *ValidationError) Error() string { return e.Message }

func (c *Client) post(ctx context.Context, path string, payload any) (int, []byte, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return 0, nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+path, bytes.NewReader(body))
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json")
	res, err := c.HTTP.Do(req)
	if err != nil {
		return 0, nil, fmt.Errorf("EasySMS недоступен: %w", err)
	}
	defer res.Body.Close()
	raw, _ := io.ReadAll(io.LimitReader(res.Body, 1<<20))
	return res.StatusCode, raw, nil
}

func (c *Client) envelope(payload map[string]any) map[string]any {
	return map[string]any{
		"Header":  map[string]string{"login": c.Login, "password": c.Password},
		"Payload": payload,
	}
}

func explain(status int, raw []byte) error {
	if status == http.StatusForbidden {
		return fmt.Errorf("логин или пароль EasySMS отклонены")
	}
	var body struct {
		Error   *apiError `json:"error"`
		Message struct {
			Error *apiError `json:"error"`
		} `json:"message"`
	}
	_ = json.Unmarshal(raw, &body)
	api := body.Error
	if api == nil {
		api = body.Message.Error
	}
	if api != nil && (api.Code != "" || api.Message != "") {
		text := api.Message
		if text == "" {
			text = "EasySMS отклонил сообщение"
		}
		if api.Code != "" {
			text = api.Code + ": " + text
		}
		return &ValidationError{Message: text}
	}
	if status >= 500 {
		return fmt.Errorf("EasySMS временно недоступен (%d)", status)
	}
	if status >= 400 {
		return fmt.Errorf("EasySMS отклонил запрос (%d)", status)
	}
	return nil
}
