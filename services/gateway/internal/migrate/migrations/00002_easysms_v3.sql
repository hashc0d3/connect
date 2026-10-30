-- +goose Up
ALTER TABLE messages ADD COLUMN client_message_id TEXT NOT NULL DEFAULT '';
ALTER TABLE messages ADD COLUMN upstream_message_id TEXT NOT NULL DEFAULT '';
CREATE UNIQUE INDEX messages_client_message_id_key ON messages (client_message_id) WHERE client_message_id <> '';

ALTER TABLE otp_requests ADD COLUMN client_message_id TEXT NOT NULL DEFAULT '';
CREATE INDEX otp_requests_client_message_id_idx ON otp_requests (client_message_id);

-- +goose Down
DROP INDEX IF EXISTS otp_requests_client_message_id_idx;
ALTER TABLE otp_requests DROP COLUMN IF EXISTS client_message_id;
DROP INDEX IF EXISTS messages_client_message_id_key;
ALTER TABLE messages DROP COLUMN IF EXISTS upstream_message_id;
ALTER TABLE messages DROP COLUMN IF EXISTS client_message_id;
