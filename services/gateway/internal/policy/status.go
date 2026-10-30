package policy

const (
	StatusQueue       = 0
	StatusDelivered   = 1
	StatusUndelivered = 2
	StatusSent        = 3
	StatusWait        = 4
	StatusReject      = 6
	StatusModeration  = 8
)

func ExtendStatus(code int) string {
	switch code {
	case StatusQueue:
		return "queue"
	case StatusDelivered:
		return "delivery"
	case StatusUndelivered:
		return "undelivered"
	case StatusSent:
		return "sent"
	case StatusWait:
		return "wait"
	case StatusReject:
		return "reject"
	case StatusModeration:
		return "moderation"
	default:
		return "wait"
	}
}
