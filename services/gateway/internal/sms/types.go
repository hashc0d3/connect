package sms

import (
	"time"

	"github.com/shopspring/decimal"
)

type SendInput struct {
	Numbers     []string
	GroupIDs    []int64
	Text        string
	Sign        string
	DateSend    time.Time
	SmoothUntil time.Time
	CallbackURL string
	ShortLink   bool
	Translit    bool
}

type Item struct {
	ID           int64
	From         string
	Number       string
	Text         string
	Status       int
	ExtendStatus string
	Channel      string
	Cost         string
	DateCreate   time.Time
	DateSend     time.Time
	DateAnswer   *time.Time
}

type SendResult struct {
	Message     string
	Items       []Item
	Excluded    int
	Blacklisted int
	Cost        string
	Held        bool
}

type draft struct {
	Number          string
	Text            string
	Parts           int
	Cost            decimal.Decimal
	Operator        string
	When            time.Time
	ClientMessageID string
	UpstreamID      string
}
