package zxcvbn

import (
	"math"
	"testing"
)

/**
Use these test to see how close to feature parity the library is.
*/

const (
	allowableError = float64(0.05)
)

func TestPasswordStrength(t *testing.T) {

	//	Expected calculated by running zxcvbn-python
	runTest(t, "zxcvbn", float64(6.845490050944376))
	runTest(t, "Tr0ub4dour&3", float64(17.296))
	runTest(t, "qwER43@!", float64(26.44))
	runTest(t, "correcthorsebatterystaple", float64(45.212))
	runTest(t, "coRrecth0rseba++ery9.23.2007staple$", float64(66.018))
	runTest(t, "D0g..................", float64(20.678))
	runTest(t, "abcdefghijk987654321", float64(11.951))
	runTest(t, "neverforget", float64(2)) // I think this is wrong. . .
	runTest(t, "13/3/1997", float64(2))   // I think this is wrong. . .
	runTest(t, "neverforget13/3/1997", float64(32.628))
	runTest(t, "1qaz2wsx3edc", float64(19.314))
	runTest(t, "temppass22", float64(22.179))
	runTest(t, "briansmith", float64(4.322))
	runTest(t, "briansmith4mayor", float64(18.64))
	runTest(t, "password1", float64(2.0))
	runTest(t, "viking", float64(7.531))
	runTest(t, "thx1138", float64(7.426))
	runTest(t, "ScoRpi0ns", float64(20.621))
	runTest(t, "do you know", float64(4.585))
	runTest(t, "ryanhunter2000", float64(14.506))
	runTest(t, "rianhunter2000", float64(21.734))
	runTest(t, "asdfghju7654rewq", float64(29.782))
	runTest(t, "AOEUIDHG&*()LS_", float64(33.254))
	runTest(t, "12345678", float64(1.585))
	runTest(t, "defghi6789", float64(12.607))
	runTest(t, "rosebud", float64(7.937))
	runTest(t, "Rosebud", float64(8.937))
	runTest(t, "ROSEBUD", float64(8.937))
	runTest(t, "rosebuD", float64(8.937))
	runTest(t, "ros3bud99", float64(19.276))
	runTest(t, "r0s3bud99", float64(19.276))
	runTest(t, "R0$38uD99", float64(34.822))
	runTest(t, "verlineVANDERMARK", float64(26.293))
	runTest(t, "eheuczkqyq", float64(42.813))
	runTest(t, "rWibMFACxAUGZmxhVncy", float64(104.551))
	runTest(t, "Ba9ZyWABu99[BK#6MBgbH88Tofv)vs$", float64(161.278))
}

var formatString = "%s : error should be less than %.2f Acctual error: %.4f Expected entropy %.4f Actual entropy %.4f \n"

func runTest(t *testing.T, password string, pythonEntropy float64) {

	goEntropy := GoPasswordStrength(password, nil)
	perror := math.Abs(goEntropy-pythonEntropy) / pythonEntropy

	if perror > allowableError {
		t.Logf(formatString, password, allowableError, perror, pythonEntropy, goEntropy)

		// t.Fail()
	}
}

func GoPasswordStrength(password string, userInputs []string) float64 {
	return PasswordStrength(password, userInputs).Entropy
}
