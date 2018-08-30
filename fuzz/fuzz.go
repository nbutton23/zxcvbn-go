package fuzz

import (
	"github.com/nbutton23/zxcvbn-go"
)

// Fuzz is used to run https://github.com/dvyukov/go-fuzz
func Fuzz(data []byte) int {
	password := string(data)

	_ = zxcvbn.PasswordStrength(password, nil)
	return 1
}
