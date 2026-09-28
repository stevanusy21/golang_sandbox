package utils

import (
	"net/mail"
	"strings"
)

func IsEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
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