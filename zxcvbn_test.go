package zxcvbn

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"os/exec"
	"testing"
	"testing/quick"
	"time"
)

func TestPasswordStrength(t *testing.T) {
	cfg := &quick.Config{Rand: rand.New(rand.NewSource(time.Now().Unix()))}
	err := quick.CheckEqual(GoPasswordStrength, PythonPasswordStrength, cfg)
	if err != nil {
		t.Error(err)
	}
}

func GoPasswordStrength(password string, userInputs []string) float64 {
	return PasswordStrength(password, userInputs).Entropy
}

func PythonPasswordStrength(password string, userInputs []string) float64 {
	cmd := exec.Command("python", append([]string{"-", password}, userInputs...)...)
	cmd.Stdin = bytes.NewBufferString(py)

	o, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("outErr:", err)
	}

	var pmatch pyMatch
	if err := json.Unmarshal(o, &pmatch); err != nil {
		fmt.Println("json:", err)
	}

	return pmatch.Entropy
}

const py = `import zxcvbn
import json
import sys

print json.dumps(zxcvbn.password_strength(sys.argv[1], sys.argv[2:len(sys.argv)]))
`

type pyMatch struct {
	CalcTime         float64 `json:"calc_time"`
	CrackTime        float64 `json:"crack_time"`
	CrackTimeDisplay string  `json:"crack_time_display"`
	Entropy          float64 `json:"entropy"`
	MatchSequence    []struct {
		BaseEntropy      float64 `json:"base_entropy"`
		DictionaryName   string  `json:"dictionary_name"`
		Entropy          float64 `json:"entropy"`
		I                int64   `json:"i"`
		J                int64   `json:"j"`
		L33tEntropy      float64 `json:"l33t_entropy"`
		MatchedWord      string  `json:"matched_word"`
		Pattern          string  `json:"pattern"`
		Rank             int64   `json:"rank"`
		Token            string  `json:"token"`
		UppercaseEntropy float64 `json:"uppercase_entropy"`
	} `json:"match_sequence"`
	Password string  `json:"password"`
	Score    float64 `json:"score"`
}
