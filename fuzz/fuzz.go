package fuzz

import (
	"github.com/nbutton23/zxcvbn-go"
)

func Fuzz(data []byte) int {
	password := string(data)

	_ = zxcvbn.PasswordStrength(password, nil)
	return 1
}
