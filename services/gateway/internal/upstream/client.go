package upstream

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	BaseURL      string
	ConnectID    string
	Login        string
	Password     string
	WebhookToken string
	Mock         bool
	RevealOTP    bool
	HTTP         *http.Client
}

func New(baseURL, connectID, login, password, webhookToken string, mock, reveal bool) *Client {
	login = strings.TrimSpace(login)
	password = strings.TrimSpace(password)
	connectID = strings.Trim(strings.TrimSpace(connectID), "/")
	if login == "" || password == "" || connectID == "" {
		mock = true
	}
	return &Client{
		BaseURL:      strings.TrimRight(baseURL, "/"),
		ConnectID:    connectID,
		Login:        login,
		Password:     password,
		WebhookToken: strings.TrimSpace(webhookToken),
		Mock:         mock,
		RevealOTP:    reveal,
		HTTP:         &http.Client{Timeout: 20 * time.Second},
	}
}

func ClientMessageID() (string, error) {
	var buf [8]byte
	if _, err := rand.Read(buf[:]); err != nil {
		return "", err
	}
	return "cmsg_" + hex.EncodeToString(buf[:]), nil
}
