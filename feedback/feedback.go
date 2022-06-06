package feedback

import "github.com/nbutton23/zxcvbn-go/match"

// Feedback represents the feedback for a weak password
type Feedback struct {
	Warning string
	Suggestions []string
}

func GetFeedback(score int, sequence []match.Match) Feedback {
	return Feedback{}
}
