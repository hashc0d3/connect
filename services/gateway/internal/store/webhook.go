package store

import "context"

func (s *Store) ApplyDelivery(ctx context.Context, clientID, upstreamID string, status int, extend, upstream string) (bool, error) {
	tag, err := s.pool.Exec(ctx, `
		UPDATE messages
		SET status = $2, extend_status = $3, upstream_state = $4,
			upstream_message_id = CASE WHEN $5 = '' THEN upstream_message_id ELSE $5 END,
			date_answer = now()
		WHERE client_message_id = $1 AND client_message_id <> ''
	`, clientID, status, extend, upstream, upstreamID)
	if err != nil {
		return false, err
	}
	return tag.RowsAffected() > 0, nil
}

func (s *Store) SaveOTPCode(ctx context.Context, clientID, codeHash string) (bool, error) {
	tag, err := s.pool.Exec(ctx, `
		UPDATE otp_requests
		SET code_hash = $2
		WHERE client_message_id = $1 AND status = 'sent' AND client_message_id <> ''
	`, clientID, codeHash)
	if err != nil {
		return false, err
	}
	return tag.RowsAffected() > 0, nil
}
