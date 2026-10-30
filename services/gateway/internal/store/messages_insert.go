package store

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"
)

type Message struct {
	ID           int64
	UserID       string
	From         string
	Number       string
	Text         string
	Status       int
	ExtendStatus string
	Channel      string
	Cost         decimal.Decimal
	DateCreate   time.Time
	DateSend     time.Time
	DateAnswer   *time.Time
	Upstream     string
	CallbackURL  string
}

type NewMessage struct {
	UserID      string
	CampaignID  int64
	From        string
	Number      string
	Text        string
	Status      int
	Extend      string
	Channel     string
	Cost        decimal.Decimal
	Parts       int
	DateCreate  time.Time
	DateSend    time.Time
	Upstream    string
	CallbackURL string
	ClientID    string
	UpstreamID  string
}

func (s *Store) InsertCampaign(ctx context.Context, tx pgx.Tx, userID, original, text, sign, status string, recipients, excluded int, cost decimal.Decimal, scheduled, smooth *time.Time) (int64, error) {
	var id int64
	err := tx.QueryRow(ctx, `
		INSERT INTO campaigns (
			user_id, original_text, text, sign_name, status, recipient_count,
			excluded_duplicates, cost, scheduled_at, smooth_until
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8::numeric,$9,$10)
		RETURNING id
	`, userID, original, text, sign, status, recipients, excluded, cost.StringFixed(4), scheduled, smooth).Scan(&id)
	return id, err
}

func (s *Store) InsertMessage(ctx context.Context, tx pgx.Tx, m NewMessage) (int64, error) {
	var id int64
	err := tx.QueryRow(ctx, `
		INSERT INTO messages (
			user_id, campaign_id, from_sign, number, text, status, extend_status,
			channel, cost, parts, date_create, date_send, upstream_state, upstream_error, callback_url,
			client_message_id, upstream_message_id
		) VALUES (
			$1,$2,$3,$4,$5,$6,$7,$8,$9::numeric,$10,$11,$12,$13,'',$14,$15,$16
		) RETURNING id
	`, m.UserID, m.CampaignID, m.From, m.Number, m.Text, m.Status, m.Extend, m.Channel,
		m.Cost.StringFixed(4), m.Parts, m.DateCreate, m.DateSend, m.Upstream, m.CallbackURL, m.ClientID, m.UpstreamID).Scan(&id)
	return id, err
}
