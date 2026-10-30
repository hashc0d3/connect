package policy

import (
	"fmt"
	"regexp"
	"strings"
)

var nonDigit = regexp.MustCompile(`\D`)

// NormalizePhone приводит номер к виду EasySMS: только цифры, код страны, без «+».
func NormalizePhone(raw string) (string, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return "", fmt.Errorf("номер пустой")
	}
	if strings.Contains(trimmed, "+") {
		return "", fmt.Errorf("номер передаётся без «+»")
	}
	if nonDigit.MatchString(trimmed) {
		compact := nonDigit.ReplaceAllString(trimmed, "")
		if compact == "" {
			return "", fmt.Errorf("в номере нет цифр")
		}
		trimmed = compact
	}
	if len(trimmed) < 10 {
		return "", fmt.Errorf("номер короче 10 цифр")
	}
	if len(trimmed) > 15 {
		return "", fmt.Errorf("номер длиннее 15 цифр")
	}
	if trimmed[0] == '0' {
		return "", fmt.Errorf("номер должен начинаться с кода страны")
	}
	return trimmed, nil
}
