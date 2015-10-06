package matching
import (
	"strings"
	"regexp"
	"strconv"
	"zxcvbn-go/frequency"
	"path/filepath"
	"zxcvbn-go/adjacency"
	"zxcvbn-go/match"
	"sort"
)

var (
	DICTIONARY_MATCHERS []func(password string) []match.Match
	MATCHERS []func(password string) []match.Match
	ADJACENCY_GRAPHS []adjacency.AdjacencyGraph;
	KEYBOARD_STARTING_POSITIONS int
	KEYBOARD_AVG_DEGREE float64
	KEYPAD_STARTING_POSITIONS int
	KEYPAD_AVG_DEGREE float64
)

const (
	DATE_RX_YEAR_SUFFIX string = `((\d{1,2})(\s|-|\/|\\|_|\.)(\d{1,2})(\s|-|\/|\\|_|\.)(19\d{2}|200\d|201\d|\d{2}))`
	DATE_RX_YEAR_PREFIX string = `((19\d{2}|200\d|201\d|\d{2})(\s|-|/|\\|_|\.)(\d{1,2})(\s|-|/|\\|_|\.)(\d{1,2}))`
	DATE_WITHOUT_SEP_MATCH string = `\d{4,8}`
)


func init() {

}

func Omnimatch(password string, userInputs []string) (matches []match.Match) {

	if DICTIONARY_MATCHERS == nil || ADJACENCY_GRAPHS == nil {
		loadFrequencyList()
	}

	if userInputs != nil {
		userInputMatcher := buildDictMatcher("user_inputs", buildRankedDict(userInputs))
		matches = userInputMatcher(password)
	}

	for _, matcher := range MATCHERS {
		mtemp := matcher(password)
		matches = append(matches, mtemp...)
	}
	sort.Sort(match.Matches(matches))
	return matches
}

func loadFrequencyList() {
	maleFilePath, _ := filepath.Abs("frequency/MaleNames.json")
	femaleFilePath, _ := filepath.Abs("frequency/FemaleNames.json")
	surnameFilePath, _ := filepath.Abs("frequency/Surnames.json")
	englishFilePath, _ := filepath.Abs("frequency/English.json")
	passwordsFilePath, _ := filepath.Abs("frequency/Passwords.json")

	DICTIONARY_MATCHERS = append(DICTIONARY_MATCHERS, buildDictMatcher("MaleNames", buildRankedDict(frequency.GetStringListFromFile(maleFilePath))))
	DICTIONARY_MATCHERS = append(DICTIONARY_MATCHERS, buildDictMatcher("FemaleNames", buildRankedDict(frequency.GetStringListFromFile(femaleFilePath))))
	DICTIONARY_MATCHERS = append(DICTIONARY_MATCHERS, buildDictMatcher("Surnames", buildRankedDict(frequency.GetStringListFromFile(surnameFilePath))))
	DICTIONARY_MATCHERS = append(DICTIONARY_MATCHERS, buildDictMatcher("English", buildRankedDict(frequency.GetStringListFromFile(englishFilePath))))
	DICTIONARY_MATCHERS = append(DICTIONARY_MATCHERS, buildDictMatcher("Passwords", buildRankedDict(frequency.GetStringListFromFile(passwordsFilePath))))

	qwertyfilePath, _ := filepath.Abs("adjacency/Qwerty.json")
	dvorakfilePath, _ := filepath.Abs("adjacency/Dvorak.json")
	keypadfilePath, _ := filepath.Abs("adjacency/Keypad.json")
	macKeypadfilePath, _ := filepath.Abs("adjacency/MacKeypad.json")

	qwertGraph := adjacency.GetAdjancencyGraphFromFile(qwertyfilePath, "qwert")
	keypadGraph := adjacency.GetAdjancencyGraphFromFile(keypadfilePath, "keypad")


	KEYBOARD_AVG_DEGREE = qwertGraph.CalculateAvgDegree()
	KEYBOARD_STARTING_POSITIONS = len(qwertGraph.Graph)
	KEYPAD_AVG_DEGREE = keypadGraph.CalculateAvgDegree()
	KEYPAD_STARTING_POSITIONS = len(keypadGraph.Graph)

	ADJACENCY_GRAPHS = append(ADJACENCY_GRAPHS, qwertGraph)
	ADJACENCY_GRAPHS = append(ADJACENCY_GRAPHS, adjacency.GetAdjancencyGraphFromFile(dvorakfilePath, "dvorak"))
	ADJACENCY_GRAPHS = append(ADJACENCY_GRAPHS, keypadGraph)
	ADJACENCY_GRAPHS = append(ADJACENCY_GRAPHS, adjacency.GetAdjancencyGraphFromFile(macKeypadfilePath, "macKepad"))

	MATCHERS = append(MATCHERS, DICTIONARY_MATCHERS...)
	MATCHERS = append(MATCHERS, SpatialMatch)
}


func buildDictMatcher(dictName string, rankedDict map[string]int) func(password string) []match.Match {
	return func(password string) []match.Match {
		matches := dictionaryMatch(password, dictName, rankedDict)
		for _, v := range matches {
			v.DictionaryName = dictName
		}
		return matches
	}

}

func dictionaryMatch(password string, dictionaryName string, rankedDict map[string]int) []match.Match {
	length := len(password)
	var results []match.Match
	pwLower := strings.ToLower(password)

	for i := 0; i < length; i++ {
		for j := i; j < length; j++ {
			word := pwLower[i:j + 1]
			if val, ok := rankedDict[word]; ok {
				results = append(results, match.Match{Pattern:"dictionary",
					DictionaryName:dictionaryName,
					I:i,
					J:j,
					Token:password[i:j + 1],
					MatchedWord:word,
					Rank:float64(val)})
			}
		}
	}

	return results
}

func buildRankedDict(unrankedList []string) map[string]int {

	result := make(map[string]int)

	for i, v := range unrankedList {
		result[strings.ToLower(v)] = i + 1
	}

	return result
}

func checkDate(day, month, year int64) (bool, int64, int64, int64) {
	if (12 <= month && month <= 31) && day <= 12 {
		day, month = month, day
	}

	if day > 31 || month > 12 {
		return false, 0, 0, 0
	}

	if !(1900 <= year && year <= 2019) {
		return false, 0, 0, 0
	}

	return true, day, month, year
}

func DateSepMatch(password string) []match.DateMatch {

	var matches []match.DateMatch

	matcher := regexp.MustCompile(DATE_RX_YEAR_SUFFIX)
	for _, v := range matcher.FindAllString(password, len(password)) {
		splitV := matcher.FindAllStringSubmatch(v, len(v))
		i := strings.Index(password, v)
		j := i + len(v)
		day, _ := strconv.ParseInt(splitV[0][4], 10, 16)
		month, _ := strconv.ParseInt(splitV[0][2], 10, 16)
		year, _ := strconv.ParseInt(splitV[0][6], 10, 16)
		match := match.DateMatch{Day:day, Month:month, Year:year, Separator:splitV[0][5], I:i, J:j }
		matches = append(matches, match)
	}


	matcher = regexp.MustCompile(DATE_RX_YEAR_PREFIX)
	for _, v := range matcher.FindAllString(password, len(password)) {
		splitV := matcher.FindAllStringSubmatch(v, len(v))
		i := strings.Index(password, v)
		j := i + len(v)
		day, _ := strconv.ParseInt(splitV[0][4], 10, 16)
		month, _ := strconv.ParseInt(splitV[0][6], 10, 16)
		year, _ := strconv.ParseInt(splitV[0][2], 10, 16)
		match := match.DateMatch{Day:day, Month:month, Year:year, Separator:splitV[0][5], I:i, J:j }
		matches = append(matches, match)
	}

	var out []match.DateMatch
	for _, match := range matches {
		if valid, day, month, year := checkDate(match.Day, match.Month, match.Year); valid {
			match.Pattern = "date"
			match.Day = day
			match.Month = month
			match.Year = year
			out = append(out, match)
		}
	}
	return out

}
type DateMatchCandidate struct {
	DayMonth string
	Year     string
	I, J     int
}
//TODO I think Im doing this wrong.
func dateWithoutSepMatch(password string) (matches []match.DateMatch) {
	matcher := regexp.MustCompile(DATE_WITHOUT_SEP_MATCH)
	for _, v := range matcher.FindAllString(password, len(password)) {
		i := strings.Index(password, v)
		j := i + len(v)
		length := len(v)
		lastIndex := length - 1
		var candidatesRoundOne []DateMatchCandidate

		if length <= 6 {
			//2-digit year prefix
			candidatesRoundOne = append(candidatesRoundOne, buildDateMatchCandidate(v[2:], v[0:2], i, j))

			//2-digityear suffix
			candidatesRoundOne = append(candidatesRoundOne, buildDateMatchCandidate(v[0:lastIndex - 2], v[lastIndex - 2:], i, j))
		}
		if length >= 6 {
			//4-digit year prefix
			candidatesRoundOne = append(candidatesRoundOne, buildDateMatchCandidate(v[4:], v[0:4], i, j))

			//4-digit year sufix
			candidatesRoundOne = append(candidatesRoundOne, buildDateMatchCandidate(v[0:lastIndex - 4], v[lastIndex - 4:], i, j))
		}

		var candidatesRoundTwo []match.DateMatch
		for _, c := range candidatesRoundOne {
			if len(c.DayMonth) == 2 {
				candidatesRoundTwo = append(candidatesRoundTwo, buildDateMatchCandidateTwo(c.DayMonth[0], c.DayMonth[1], c.Year, c.I, c.J))
			}
		}
	}

	return matches
}

func buildDateMatchCandidate(dayMonth, year string, i, j int) DateMatchCandidate {
	return DateMatchCandidate{DayMonth: dayMonth, Year:year, I:i, J:j}
}

func buildDateMatchCandidateTwo(day, month byte, year string, i, j int) match.DateMatch {
	sDay := string(day)
	sMonth := string(month)
	intDay, _ := strconv.ParseInt(sDay, 10, 16)
	intMonth, _ := strconv.ParseInt(sMonth, 10, 16)
	intYear, _ := strconv.ParseInt(year, 10, 16)

	return match.DateMatch{Day:intDay, Month:intMonth, Year:intYear, I:i, J:j}
}


func SpatialMatch(password string) (matches []match.Match) {
	for _, graph := range ADJACENCY_GRAPHS {
		matches = append(matches, spatialMatchHelper(password, graph)...)
	}
	return matches
}

func spatialMatchHelper(password string, graph adjacency.AdjacencyGraph) (matches []match.Match) {
	for i := 0; i < len(password) - 1; {
		j := i + 1
		lastDirection := -99 //and int that it should never be!
		turns := 0
		shiftedCount := 0

		for ;; {
			prevChar := password[j - 1]
			found := false
			foundDirection := -1
			curDirection := -1
			adjacents := graph.Graph[string(prevChar)]
			//				Consider growing pattern by one character if j hasn't gone over the edge
			if j < len(password) {
				curChar := password[j]
				for _, adj := range adjacents {
					curDirection += 1

					if strings.Index(adj, string(curChar)) != -1 {
						found = true
						foundDirection = curDirection

						if strings.Index(adj, string(curChar)) == 1 {
							//								index 1 in the adjacency means the key is shifted, 0 means unshifted: A vs a, % vs 5, etc.
							//								for example, 'q' is adjacent to the entry '2@'. @ is shifted w/ index 1, 2 is unshifted.

							shiftedCount += 1
						}

						if lastDirection != foundDirection {
							//								adding a turn is correct even in the initial case when last_direction is null:
							//								every spatial pattern starts with a turn.
							turns += 1
							lastDirection = foundDirection
						}
						break
					}
				}
			}

			//				if the current pattern continued, extend j and try to grow again
			if found {
				j += 1
			} else {
				//					otherwise push the pattern discovered so far, if any...

				//					don't consider length 1 or 2 chains.
				if j - i > 2 {
					matches = append(matches, match.Match{Pattern:"spatial", I:i, J:j - 1, Token:password[i:j], DictionaryName:graph.Name, Turns:turns, ShiftedCount:shiftedCount })
				}
				//					. . . and then start a new search from the rest of the password
				i = j
				break
			}
		}

	}
	return matches
}
