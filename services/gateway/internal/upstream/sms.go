package upstream

import "context"

func (c *Client) SendSMS(ctx context.Context, phones []string, text string) error {
	if c.Mock {
		return nil
	}
	receivers := make([]Receiver, 0, len(phones))
	for _, phone := range phones {
		receivers = append(receivers, Receiver{Phone: phone})
	}
	return c.post(ctx, "/api/send-sms", map[string]any{
		"receivers": receivers,
		"text":      text,
	})
}
