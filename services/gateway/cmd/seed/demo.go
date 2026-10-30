package main

import (
	"context"
	"fmt"

	"github.com/connect/gateway/internal/authn"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

func seed(ctx context.Context, pool *pgxpool.Pool) error {
	settings := map[string]string{
		"min_text_length": "4", "max_receivers": "1000", "markup_percent": "20",
		"require_approved_sign": "true", "require_approved_contract": "false",
		"block_blacklist": "true", "prefix_sign_in_text": "false",
		"promo_percent": "10", "promo_min_sum": "3000", "promo_code": "CONNECT10",
	}
	for key, value := range settings {
		if _, err := pool.Exec(ctx, `
			INSERT INTO platform_settings (key, value) VALUES ($1,$2)
			ON CONFLICT (key) DO UPDATE SET value = EXCLUDED.value
		`, key, value); err != nil {
			return err
		}
	}
	tariffs := [][3]string{
		{"PAY_SIGN", "MEGAFON", "2.20"}, {"PAY_SIGN", "MTS", "2.20"}, {"PAY_SIGN", "BEELINE", "2.30"}, {"PAY_SIGN", "t2", "2.10"}, {"PAY_SIGN", "OTHER", "2.80"},
		{"FREE_SIGN", "MEGAFON", "4.50"}, {"FREE_SIGN", "MTS", "4.20"}, {"FREE_SIGN", "BEELINE", "4.40"}, {"FREE_SIGN", "t2", "4.10"}, {"FREE_SIGN", "OTHER", "4.80"},
		{"SERVICE", "MEGAFON", "1.70"}, {"SERVICE", "MTS", "1.60"}, {"SERVICE", "BEELINE", "1.70"}, {"SERVICE", "t2", "1.55"}, {"SERVICE", "OTHER", "2.10"},
		{"AUTH", "MEGAFON", "1.50"}, {"AUTH", "MTS", "1.50"}, {"AUTH", "BEELINE", "1.60"}, {"AUTH", "t2", "1.45"}, {"AUTH", "OTHER", "1.90"},
	}
	for _, row := range tariffs {
		if _, err := pool.Exec(ctx, `
			INSERT INTO tariffs (channel, operator, price) VALUES ($1,$2,$3::numeric)
			ON CONFLICT (channel, operator) DO UPDATE SET price = EXCLUDED.price
		`, row[0], row[1], row[2]); err != nil {
			return err
		}
	}
	demoHash, err := bcrypt.GenerateFromPassword([]byte("demo12345"), 10)
	if err != nil {
		return err
	}
	adminHash, err := bcrypt.GenerateFromPassword([]byte("admin12345"), 10)
	if err != nil {
		return err
	}
	var userID string
	if err := pool.QueryRow(ctx, `
		INSERT INTO users (email, password_hash, name, role, balance)
		VALUES ('demo@connect.local', $1, 'Демо', 'customer', 1500)
		ON CONFLICT (email) DO UPDATE SET password_hash = EXCLUDED.password_hash, balance = 1500
		RETURNING id::text
	`, string(demoHash)).Scan(&userID); err != nil {
		return err
	}
	if _, err := pool.Exec(ctx, `
		INSERT INTO users (email, password_hash, name, role, balance)
		VALUES ('admin@connect.local', $1, 'Администратор', 'admin', 0)
		ON CONFLICT (email) DO UPDATE SET password_hash = EXCLUDED.password_hash, role = 'admin'
	`, string(adminHash)); err != nil {
		return err
	}
	key := "cnt_demo_7f3a9c2e1b"
	if _, err := pool.Exec(ctx, `
		INSERT INTO api_keys (user_id, name, prefix, key_hash)
		VALUES ($1, 'Демо', $2, $3) ON CONFLICT (key_hash) DO NOTHING
	`, userID, key[:12], authn.HashKey(key)); err != nil {
		return err
	}
	if _, err := pool.Exec(ctx, `
		INSERT INTO signs (user_id, name, channel, status, extend_status)
		VALUES ($1, 'CONNECT', 'PAY_SIGN', 1, 'active')
		ON CONFLICT (user_id, name) DO UPDATE SET status = 1, extend_status = 'active'
	`, userID); err != nil {
		return err
	}
	if _, err := pool.Exec(ctx, `
		INSERT INTO billing_profiles (user_id, type, title, inn, status)
		SELECT $1, 'legal', 'Connect', '7700000000', 'approved'
		WHERE NOT EXISTS (SELECT 1 FROM billing_profiles WHERE user_id = $1 AND title = 'Connect')
	`, userID); err != nil {
		return err
	}
	var groupID int
	if err := pool.QueryRow(ctx, `
		INSERT INTO "groups" (user_id, name) VALUES ($1, 'Клиенты')
		ON CONFLICT (user_id, name) DO UPDATE SET name = EXCLUDED.name
		RETURNING id
	`, userID).Scan(&groupID); err != nil {
		return err
	}
	for _, row := range []struct{ number, fname, lname string }{
		{"79991234567", "Иван", "Иванов"},
		{"79997654321", "Анна", "Соколова"},
	} {
		if _, err := pool.Exec(ctx, `
			INSERT INTO contacts (user_id, number, fname, lname, operator_code, operator_name, personal)
			VALUES ($1,$2,$3,$4,6,'t2', false)
			ON CONFLICT (user_id, number) DO NOTHING
		`, userID, row.number, row.fname, row.lname); err != nil {
			return err
		}
		if _, err := pool.Exec(ctx, `
			INSERT INTO contact_groups (contact_id, group_id)
			SELECT id, $2 FROM contacts WHERE user_id = $1 AND number = $3
			ON CONFLICT DO NOTHING
		`, userID, groupID, row.number); err != nil {
			return fmt.Errorf("group link: %w", err)
		}
	}
	if _, err := pool.Exec(ctx, `
		INSERT INTO templates (user_id, name, body)
		SELECT $1, 'Визит', '{Имя}, ждём вас в нашем магазине'
		WHERE NOT EXISTS (SELECT 1 FROM templates WHERE user_id = $1 AND name = 'Визит')
	`, userID); err != nil {
		return err
	}
	_, err = pool.Exec(ctx, `
		INSERT INTO cards (user_id, brand, last4)
		SELECT $1, 'Visa', '4242'
		WHERE NOT EXISTS (SELECT 1 FROM cards WHERE user_id = $1 AND last4 = '4242')
	`, userID)
	return err
}
