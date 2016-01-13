package matching

import (
	"github.com/nbutton23/zxcvbn-go/Godeps/_workspace/src/github.com/stretchr/testify/assert"
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
