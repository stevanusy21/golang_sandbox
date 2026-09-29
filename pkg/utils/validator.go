package utils

import (
	"net/mail"
	"strings"
)

func IsEmpty(s any) bool {
	switch v := s.(type) {
	case string:
		return strings.TrimSpace(v) == ""
	case int, int8, int16, int32, int64:
		return v == 0
	case float32, float64:
		return v == 0.0
	case bool:
		return !v
	default:
		return false
	}
}

func MinLength(s string, min int) bool {
	return len(strings.TrimSpace(s)) >= min
}

func MaxLength(s string, max int) bool {
	return len(strings.TrimSpace(s)) <= max
}

func IsEmail(s string) bool {
	_, err := mail.ParseAddress(s)
	return err == nil
}