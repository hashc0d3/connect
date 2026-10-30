package upstream

import (
	"net/http"
	"strings"
	"time"
)

type Client struct {
	BaseURL   string
	Token     string
	Mock      bool
	RevealOTP bool
	HTTP      *http.Client
}

func New(baseURL, token string, mock, reveal bool) *Client {
	if strings.TrimSpace(token) == "" {
		mock = true
	}
	return &Client{
		BaseURL:   strings.TrimRight(baseURL, "/"),
		Token:     token,
		Mock:      mock,
		RevealOTP: reveal,
		HTTP:      &http.Client{Timeout: 20 * time.Second},
	}
}

type Receiver struct {
	Phone string `json:"phone"`
}
