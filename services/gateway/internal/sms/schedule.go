package sms

import (
	"crypto/rand"
	"time"
)

func spread(n int, start, end time.Time) []time.Time {
	out := make([]time.Time, n)
	if n == 0 {
		return out
	}
	if end.IsZero() || !end.After(start) || n == 1 {
		for i := range out {
			out[i] = start
		}
		return out
	}
	span := end.Sub(start)
	for i := 0; i < n; i++ {
		out[i] = start.Add(time.Duration(int64(span) * int64(i) / int64(n-1)))
	}
	return out
}

func chunks(list []string, size int) [][]string {
	var out [][]string
	for len(list) > 0 {
		if len(list) < size {
			size = len(list)
		}
		out = append(out, list[:size])
		list = list[size:]
	}
	return out
}

func randomCode(n int) (string, error) {
	const alphabet = "abcdefghijklmnopqrstuvwxyz0123456789"
	buf := make([]byte, n)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	for i := range buf {
		buf[i] = alphabet[int(buf[i])%len(alphabet)]
	}
	return string(buf), nil
}
