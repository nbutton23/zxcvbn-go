package matching

import (
	"github.com/nbutton23/zxcvbn-go/Godeps/_workspace/src/github.com/stretchr/testify/assert"
	"github.com/nbutton23/zxcvbn-go/match"
	"strings"
	"testing"
)

//DateSepMatch("1991-09-11jibjab11.9.1991")
//[{date 16 25  . 9 11 1991} {date 0 10  - 9 11 1991}]

func TestDateSepMatch(t *testing.T) {
	matches := DateSepMatch("1991-09-11jibjab11.9.1991")

	assert.Len(t, matches, 2, "Length should be 2")

	for _, match := range matches {
		if match.Separator == "." {
			assert.Equal(t, 16, match.I)
			assert.Equal(t, 25, match.J)
			assert.Equal(t, int64(9), match.Day)
			assert.Equal(t, int64(11), match.Month)
			assert.Equal(t, int64(1991), match.Year)
		} else {
			assert.Equal(t, 0, match.I)
			assert.Equal(t, 10, match.J)
			assert.Equal(t, int64(9), match.Day)
			assert.Equal(t, int64(11), match.Month)
			assert.Equal(t, int64(1991), match.Year)
		}
	}

}

func TestRepeatMatch(t *testing.T) {
	//aaaBbBb
	matches := RepeatMatch("aaabBbB")

	assert.Len(t, matches, 2, "Lenght should be 2")

	for _, match := range matches {
		if strings.ToLower(match.DictionaryName) == "b" {
			assert.Equal(t, 3, match.I)
			assert.Equal(t, 6, match.J)
			assert.Equal(t, "bBbB", match.Token)
			assert.NotZero(t, match.Entropy, "Entropy should be set")
		} else {
			assert.Equal(t, 0, match.I)
			assert.Equal(t, 2, match.J)
			assert.Equal(t, "aaa", match.Token)
			assert.NotZero(t, match.Entropy, "Entropy should be set")

		}
	}
}

func TestSequenceMatch(t *testing.T) {
	//abcdjibjacLMNOPjibjac1234  => abcd LMNOP 1234

	matches := SequenceMatch("abcdjibjacLMNOPjibjac1234")
	assert.Len(t, matches, 3, "Lenght should be 2")

	for _, match := range matches {
		if match.DictionaryName == "lower" {
			assert.Equal(t, 0, match.I)
			assert.Equal(t, 3, match.J)
			assert.Equal(t, "abcd", match.Token)
			assert.NotZero(t, match.Entropy, "Entropy should be set")
		} else if match.DictionaryName == "upper" {
			assert.Equal(t, 10, match.I)
			assert.Equal(t, 14, match.J)
			assert.Equal(t, "LMNOP", match.Token)
			assert.NotZero(t, match.Entropy, "Entropy should be set")
		} else if match.DictionaryName == "digits" {
			assert.Equal(t, 21, match.I)
			assert.Equal(t, 24, match.J)
			assert.Equal(t, "1234", match.Token)
			assert.NotZero(t, match.Entropy, "Entropy should be set")
		} else {
			assert.True(t, false, "Unknow dictionary")
		}
	}
}

func TestSpatialMatchQwerty(t *testing.T) {
	matches := SpatialMatch("qwerty")
	assert.Len(t, matches, 1, "Lenght should be 1")
	assert.NotZero(t, matches[0].Entropy, "Entropy should be set")

	matches = SpatialMatch("asdf")
	assert.Len(t, matches, 1, "Lenght should be 1")
	assert.NotZero(t, matches[0].Entropy, "Entropy should be set")

}

func TestSpatialMatchDvorak(t *testing.T) {
	matches := SpatialMatch("aoeuidhtns")
	assert.Len(t, matches, 1, "Lenght should be 1")
	assert.NotZero(t, matches[0].Entropy, "Entropy should be set")

}

func TestDictionaryMatch(t *testing.T) {
	var matches []match.Match
	for _, dicMatcher := range DICTIONARY_MATCHERS {
		matchesTemp := dicMatcher("first")
		matches = append(matches, matchesTemp...)
	}

	assert.Len(t, matches, 4, "Lenght should be 4")
	for _, match := range matches {
		assert.NotZero(t, match.Entropy, "Entropy should be set")

	}

}

func TestDateWithoutSepMatch(t *testing.T) {
	matches := dateWithoutSepMatch("11091991")
	assert.Len(t, matches, 1, "Lenght should be 1")
}

//l33t
func TestLeetSubTable(t *testing.T){
	subs := relevantL33tSubtable("password")
	assert.Len(t, subs, 0, "password should produce no leet subs")

	subs = relevantL33tSubtable("p4ssw0rd")
	assert.Len(t, subs, 2, "p4ssw0rd should produce 2 subs")
}
