package upstream

import (
	"context"
	"encoding/json"
	"fmt"
)

type Outbound struct {
	ClientMessageID string
	Messenger       string
	Originator      string
	Recipient       string
	Text            string
}

type Accepted struct {
	MessageID       string
	ClientMessageID string
}

func (c *Client) SendMessage(ctx context.Context, msg Outbound) (Accepted, error) {
	if msg.Messenger == "" {
		msg.Messenger = "sms"
	}
	if msg.ClientMessageID == "" {
		id, err := ClientMessageID()
		if err != nil {
			return Accepted{}, err
		}
		msg.ClientMessageID = id
	}
	if c.Mock {
		return Accepted{MessageID: "mock-" + msg.ClientMessageID, ClientMessageID: msg.ClientMessageID}, nil
	}
	message := map[string]string{
		"client_message_id": msg.ClientMessageID,
		"messenger_type":    msg.Messenger,
		"recipient":         msg.Recipient,
	}
	if msg.Messenger != "flash_call" {
		message["originator"] = msg.Originator
		message["text"] = msg.Text
	}
	status, raw, err := c.post(ctx, "/api/v3/easysms/"+c.ConnectID+"/SendMessage", c.envelope(map[string]any{"message": message}))
	if err != nil {
		return Accepted{}, err
	}
	if err := explain(status, raw); err != nil {
		return Accepted{}, err
	}
	var parsed struct {
		Message struct {
			MessageID       string `json:"message_id"`
			ClientMessageID string `json:"client_message_id"`
		} `json:"message"`
	}
	if json.Unmarshal(raw, &parsed) != nil || parsed.Message.MessageID == "" {
		return Accepted{}, fmt.Errorf("EasySMS вернул ответ без message_id")
	}
	id := parsed.Message.ClientMessageID
	if id == "" {
		id = msg.ClientMessageID
	}
	return Accepted{MessageID: parsed.Message.MessageID, ClientMessageID: id}, nil
}
