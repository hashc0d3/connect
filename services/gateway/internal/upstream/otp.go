package upstream

import (
	"context"
	"crypto/rand"
	"fmt"
)

func (c *Client) SendFlashCall(ctx context.Context, clientID, phone string) (Accepted, string, error) {
	if c.Mock {
		code, err := flashCode()
		if err != nil {
			return Accepted{}, "", err
		}
		return Accepted{MessageID: "mock-" + clientID, ClientMessageID: clientID}, code, nil
	}
	accepted, err := c.SendMessage(ctx, Outbound{
		ClientMessageID: clientID,
		Messenger:       "flash_call",
		Recipient:       phone,
	})
	return accepted, "", err
}

func flashCode() (string, error) {
	var buf [3]byte
	if _, err := rand.Read(buf[:]); err != nil {
		return "", err
	}
	n := int(buf[0])<<16 | int(buf[1])<<8 | int(buf[2])
	return fmt.Sprintf("%06d", n%1000000), nil
}
