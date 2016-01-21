package entropy

import (
	"github.com/nbutton23/zxcvbn-go/Godeps/_workspace/src/github.com/stretchr/testify/assert"
	"github.com/nbutton23/zxcvbn-go/match"
	"testing"
)

/**
These test are more for make sure things don't change than check that the calculation is correct.
I need to hand calc a few of these for that
*/

func TestDictionaryEntropyCalculation(t *testing.T) {
	match := match.Match{
		Pattern: "dictionary",
		I:       0,
		J:       4,
		Token:   "first",
	}

	entropy := DictionaryEntropy(match, float64(20))
	//20  4.321928094887363

	assert.Equal(t, 4.321928094887363, entropy)
}

func TestSpatialEntropyCalculation(t *testing.T) {
	matchPlain := match.Match{
		Pattern: "spatial",
		I:       0,
		J:       5,
		Token:   "asdfgh",
	}

	entropy := SpatialEntropy(matchPlain, 0, 0)

	assert.Equal(t, 5.832890014164741, entropy)

	matchShift := match.Match{
		Pattern: "spatial",
		I:       0,
		J:       5,
		Token:   "asdFgh",
	}

	entropyShift := SpatialEntropy(matchShift, 0, 1)

	assert.Equal(t, 8.640244936222345, entropyShift)

	matchTurn := match.Match{
		Pattern: "spatial",
		I:       0,
		J:       5,
		Token:   "asdcxz",
	}

	entropyTurn := SpatialEntropy(matchTurn, 2, 0)

	assert.Equal(t, 8.307061631635431, entropyTurn)
}

func TestRepeatMatchEntropyCalculation(t *testing.T) {
	matchRepeat := match.Match{
		Pattern: "repeat",
		I:       0,
		J:       4,
		Token:   "aaaaa",
	}

	entropy := RepeatEntropy(matchRepeat)
	//20  4.321928094887363

	assert.Equal(t, 7.022367813028454, entropy)
}

func TestSequenceCalculation(t *testing.T) {
	matchLower := match.Match{
		Pattern: "sequence",
		I:       0,
		J:       4,
		Token:   "jklmn",
	}

	entropy := SequenceEntropy(matchLower, len("abcdefghijklmnopqrstuvwxyz"), true)

	assert.Equal(t, 7.022367813028454, entropy)

	matchUpper := match.Match{
		Pattern: "sequence",
		I:       0,
		J:       4,
		Token:   "JKLMN",
	}

	entropy = SequenceEntropy(matchUpper, len("abcdefghijklmnopqrstuvwxyz"), true)

	assert.Equal(t, 8.022367813028454, entropy)

	matchUpperDec := match.Match{
		Pattern: "sequence",
		I:       0,
		J:       4,
		Token:   "JKLMN",
	}

	entropy = SequenceEntropy(matchUpperDec, len("abcdefghijklmnopqrstuvwxyz"), false)

	assert.Equal(t, 9.022367813028454, entropy)

	matchDigit := match.Match{
		Pattern: "sequence",
		I:       0,
		J:       4,
		Token:   "34567",
	}

	entropy = SequenceEntropy(matchDigit, 10, true)

	assert.Equal(t, 5.643856189774724, entropy)
}
