package match

type Matches []Match
func (s Matches)Len() int {
	return len(s)
}
func (s Matches)Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s Matches) Less(i, j int) bool {
	if s[i].I < s[j].I {
		return true
	} else if s[i].I == s[j].I {
		return s[i].J < s[j].J
	} else {
		return false
	}
}
type Match struct {
	Pattern        string
	I, J           int
	Token          string
	MatchedWord    string
	Rank           float64
	DictionaryName string
	Turns          int
	ShiftedCount   int
	Entropy        float64
	RepeatedChar	string
}

type DateMatch struct {
	Pattern          string
	I, J             int
	Token            string
	Separator        string
	Day, Month, Year int64

}