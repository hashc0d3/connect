package sms

import "context"

func (s *Service) Release(ctx context.Context, id int64) error {
	msg, err := s.Store.GetMessageAny(ctx, id)
	if err != nil {
		return err
	}
	if msg == nil || msg.Upstream != "held" {
		return Fail(404, "Сообщение не ждёт модерации")
	}
	if err := s.Store.ChargeAndRelease(ctx, *msg); err != nil {
		if err.Error() == "недостаточно средств" {
			return Fail(402, "Недостаточно средств, чтобы выпустить сообщение")
		}
		return err
	}
	return nil
}
