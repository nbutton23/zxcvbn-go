package adjacency

import (
	"github.com/nbutton23/zxcvbn-go/Godeps/_workspace/src/github.com/stretchr/testify/assert"
	"testing"
)

/*
nbutton: Really the value is not as important to me than they don't change, which happened during development.
*/
func TestCalculateDegreeQwert(t *testing.T) {
	avgDegreeQwert := BuildQwerty().CalculateAvgDegree()

	assert.Equal(t, float32(1.531915), avgDegreeQwert, "Avg degree for qwerty should be 1.531915")
}

func TestCalculateDegreeDvorak(t *testing.T) {
	avgDegreeQwert := BuildDvorak().CalculateAvgDegree()

	assert.Equal(t, float32(1.531915), avgDegreeQwert, "Avg degree for dvorak should be 1.531915")
}

func TestCalculateDegreeKeypad(t *testing.T) {
	avgDegreeQwert := BuildKeypad().CalculateAvgDegree()

	assert.Equal(t, float32(0.62222224), avgDegreeQwert, "Avg degree for keypad should be 0.62222224")
}

func TestCalculateDegreeMacKepad(t *testing.T) {
	avgDegreeQwert := BuildMacKeypad().CalculateAvgDegree()

	assert.Equal(t, float32(0.6458333), avgDegreeQwert, "Avg degree for mackeyPad should be 0.6458333")
}
