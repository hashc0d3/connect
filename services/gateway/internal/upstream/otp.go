package upstream

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
)

func (c *Client) CreateOTP(ctx context.Context, phone string) (string, error) {
	if c.Mock {
		return randomCode()
	}
	return "", c.post(ctx, "/api/create-otp", map[string]any{"phone": phone})
}

func (c *Client) ValidateOTP(ctx context.Context, phone, code string) error {
	if c.Mock {
		return nil
	}
	return c.post(ctx, "/api/validate-otp", map[string]any{
		"phone": phone,
		"code":  json.Number(code),
	})
}

func randomCode() (string, error) {
	var buf [3]byte
	if _, err := rand.Read(buf[:]); err != nil {
		return "", err
	}
	n := int(buf[0])<<16 | int(buf[1])<<8 | int(buf[2])
	return fmt.Sprintf("%06d", n%1000000), nil
}
