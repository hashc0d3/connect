package policy

func EasySMSStatus(state string) (status int, upstream string, ok bool) {
	switch state {
	case "102", "106":
		return StatusDelivered, "delivered", true
	case "103", "104", "105", "107":
		return StatusUndelivered, "undelivered", true
	case "108":
		return StatusReject, "failed", true
	default:
		return 0, "", false
	}
}
