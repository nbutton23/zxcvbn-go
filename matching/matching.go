package matching

import (
	"github.com/nbutton23/zxcvbn-go/adjacency"
	"github.com/nbutton23/zxcvbn-go/frequency"
	"github.com/nbutton23/zxcvbn-go/match"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"github.com/nbutton23/zxcvbn-go/entropy"
)

var (
	DICTIONARY_MATCHERS         []func(password string) []match.Match
	MATCHERS                    []func(password string) []match.Match
	ADJACENCY_GRAPHS            []adjacency.AdjacencyGraph
	KEYBOARD_STARTING_POSITIONS int
	KEYBOARD_AVG_DEGREE         float64
	KEYPAD_STARTING_POSITIONS   int
	KEYPAD_AVG_DEGREE           float64
	L33T_TABLE                  adjacency.AdjacencyGraph

	SEQUENCES map[string]string
)

const (
	DATE_RX_YEAR_SUFFIX    string = `((\d{1,2})(\s|-|\/|\\|_|\.)(\d{1,2})(\s|-|\/|\\|_|\.)(19\d{2}|200\d|201\d|\d{2}))`
	DATE_RX_YEAR_PREFIX    string = `((19\d{2}|200\d|201\d|\d{2})(\s|-|/|\\|_|\.)(\d{1,2})(\s|-|/|\\|_|\.)(\d{1,2}))`
	DATE_WITHOUT_SEP_MATCH string = `\d{4,8}`
)

func init() {
	loadFrequencyList()
}

func Omnimatch(password string, userInputs []string) (matches []match.Match) {

	//Can I run into the issue where nil is not equal to nil?
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

	for n, list := range frequency.FrequencyLists {
		DICTIONARY_MATCHERS = append(DICTIONARY_MATCHERS, buildDictMatcher(n, buildRankedDict(list.List)))
	}

	KEYBOARD_AVG_DEGREE = adjacency.AdjacencyGph["querty"].CalculateAvgDegree()
	KEYBOARD_STARTING_POSITIONS = len(adjacency.AdjacencyGph["querty"].Graph)
	KEYPAD_AVG_DEGREE = adjacency.AdjacencyGph["keypad"].CalculateAvgDegree()
	KEYPAD_STARTING_POSITIONS = len(adjacency.AdjacencyGph["keypad"].Graph)
	L33T_TABLE = adjacency.AdjacencyGph["l33t"]

	ADJACENCY_GRAPHS = append(ADJACENCY_GRAPHS, adjacency.AdjacencyGph["qwerty"])
	ADJACENCY_GRAPHS = append(ADJACENCY_GRAPHS, adjacency.AdjacencyGph["dvorak"])
	ADJACENCY_GRAPHS = append(ADJACENCY_GRAPHS, adjacency.AdjacencyGph["keypad"])
	ADJACENCY_GRAPHS = append(ADJACENCY_GRAPHS, adjacency.AdjacencyGph["macKeypad"])

	//l33tFilePath, _ := filepath.Abs("adjacency/L33t.json")
	//L33T_TABLE = adjacency.GetAdjancencyGraphFromFile(l33tFilePath, "l33t")

	SEQUENCES = make(map[string]string)
	SEQUENCES["lower"] = "abcdefghijklmnopqrstuvwxyz"
	SEQUENCES["upper"] = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	SEQUENCES["digits"] = "0123456789"

	MATCHERS = append(MATCHERS, DICTIONARY_MATCHERS...)
	MATCHERS = append(MATCHERS, SpatialMatch)
	MATCHERS = append(MATCHERS, RepeatMatch)
	MATCHERS = append(MATCHERS, SequenceMatch)
	MATCHERS = append(MATCHERS, l33tMatch)

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
			word := pwLower[i : j+1]
			if val, ok := rankedDict[word]; ok {
				matchDic := match.Match{Pattern: "dictionary",
					DictionaryName: dictionaryName,
					I:              i,
					J:              j,
					Token:          password[i : j+1],
				}
				matchDic.Entropy = entropy.DictionaryEntropy(matchDic, float64(val))

				results = append(results, matchDic)
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
		match := match.DateMatch{Day: day, Month: month, Year: year, Separator: splitV[0][5], I: i, J: j}
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
		match := match.DateMatch{Day: day, Month: month, Year: year, Separator: splitV[0][5], I: i, J: j}
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
			candidatesRoundOne = append(candidatesRoundOne, buildDateMatchCandidate(v[0:lastIndex-2], v[lastIndex-2:], i, j))
		}
		if length >= 6 {
			//4-digit year prefix
			candidatesRoundOne = append(candidatesRoundOne, buildDateMatchCandidate(v[4:], v[0:4], i, j))

			//4-digit year sufix
			candidatesRoundOne = append(candidatesRoundOne, buildDateMatchCandidate(v[0:lastIndex-4], v[lastIndex-4:], i, j))
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
	return DateMatchCandidate{DayMonth: dayMonth, Year: year, I: i, J: j}
}

func buildDateMatchCandidateTwo(day, month byte, year string, i, j int) match.DateMatch {
	sDay := string(day)
	sMonth := string(month)
	intDay, _ := strconv.ParseInt(sDay, 10, 16)
	intMonth, _ := strconv.ParseInt(sMonth, 10, 16)
	intYear, _ := strconv.ParseInt(year, 10, 16)

	return match.DateMatch{Day: intDay, Month: intMonth, Year: intYear, I: i, J: j}
}

func SpatialMatch(password string) (matches []match.Match) {
	for _, graph := range ADJACENCY_GRAPHS {
		if graph.Graph != nil {
			matches = append(matches, spatialMatchHelper(password, graph)...)
		}
	}
	return matches
}

func spatialMatchHelper(password string, graph adjacency.AdjacencyGraph) (matches []match.Match) {

	for i := 0; i < len(password)-1; {
		j := i + 1
		lastDirection := -99 //an int that it should never be!
		turns := 0
		shiftedCount := 0

		for {
			prevChar := password[j-1]
			found := false
			foundDirection := -1
			curDirection := -1
			//My graphs seem to be wrong. . . and where the hell is qwerty
			adjacents := graph.Graph[string(prevChar)]
			//Consider growing pattern by one character if j hasn't gone over the edge
			if j < len(password) {
				curChar := password[j]
				for _, adj := range adjacents {
					curDirection += 1

					if strings.Index(adj, string(curChar)) != -1 {
						found = true
						foundDirection = curDirection

						if strings.Index(adj, string(curChar)) == 1 {
							//index 1 in the adjacency means the key is shifted, 0 means unshifted: A vs a, % vs 5, etc.
							//for example, 'q' is adjacent to the entry '2@'. @ is shifted w/ index 1, 2 is unshifted.
							shiftedCount += 1
						}

						if lastDirection != foundDirection {
							//adding a turn is correct even in the initial case when last_direction is null:
							//every spatial pattern starts with a turn.
							turns += 1
							lastDirection = foundDirection
						}
						break
					}
				}
			}

			//if the current pattern continued, extend j and try to grow again
			if found {
				j += 1
			} else {
				//otherwise push the pattern discovered so far, if any...
				//don't consider length 1 or 2 chains.
				if j-i > 2 {
					matchSpc := match.Match{Pattern: "spatial", I: i, J: j - 1, Token: password[i:j], DictionaryName: graph.Name}
					matchSpc.Entropy = entropy.SpatialEntropy(matchSpc, turns, shiftedCount)
					matches = append(matches, matchSpc)
				}
				//. . . and then start a new search from the rest of the password
				i = j
				break
			}
		}

	}
	return matches
}


func l33tMatch(password string) []match.Match {

	subsitutions := relevantL33tSubtable(password)

	permutations := getAllPermutationsOfLeetSubstitutions(password, subsitutions)

	var matches []match.Match

	for _, permutation := range permutations {
		for _, mather := range DICTIONARY_MATCHERS {
			matches = append(matches,mather(permutation)...)
		}
	}

	for _, match := range matches {
		match.Entropy += entropy.ExtraLeetEntropy(match, password)
		match.DictionaryName = match.DictionaryName + "_3117"
	}

	return matches
}

func getAllPermutationsOfLeetSubstitutions(password string, substitutionsMap map[string][]string) []string {

	var permutations []string

	for index, char := range password {
			for value, splice := range substitutionsMap {
				for _, sub := range splice {
					if string(char) == sub {
						var permutation string
						permutation = password[:index]+value+password[index+1:]

						permutations = append(permutations, permutation)
						if index < len(permutation) {
							tempPermutations := getAllPermutationsOfLeetSubstitutions(permutation[index + 1:], substitutionsMap)
							for _, temp := range tempPermutations {
								permutations = append(permutations, permutation[:index + 1] + temp)
							}

						}
					}
				}
			}
	}

	return permutations
}

func relevantL33tSubtable(password string) map[string][]string {
	relevantSubs := make(map[string][]string)
	for key, values := range L33T_TABLE.Graph {
		for _, value := range values {
			if strings.Contains(password, value) {
				relevantSubs[key] = append(relevantSubs[key], value)
			}
		}
	}
	return relevantSubs
}

func RepeatMatch(password string) []match.Match {
	var matches []match.Match

	//Loop through password. if current == prev currentStreak++ else if currentStreak > 2 {buildMatch; currentStreak = 1} prev = current
	var current, prev string
	currentStreak := 1
	var i int
	var char rune
	for i, char = range password {
		current = string(char)
		if i == 0 {
			prev = current
			continue
		}

		if strings.ToLower(current) == strings.ToLower(prev) {
			currentStreak++

		} else if currentStreak > 2 {
			iPos := i - currentStreak
			jPos := i - 1
			matchRepeat := match.Match{
				Pattern:        "repeat",
				I:              iPos,
				J:              jPos,
				Token:          password[iPos : jPos+1],
				DictionaryName: prev}
			matchRepeat.Entropy = entropy.RepeatEntropy(matchRepeat)
			matches = append(matches, matchRepeat)
			currentStreak = 1
		} else {
			currentStreak = 1
		}

		prev = current
	}

	if currentStreak > 2 {
		iPos := i - currentStreak + 1
		jPos := i
		matchRepeat := match.Match{
			Pattern:        "repeat",
			I:              iPos,
			J:              jPos,
			Token:          password[iPos : jPos+1],
			DictionaryName: prev}
		matchRepeat.Entropy = entropy.RepeatEntropy(matchRepeat)
		matches = append(matches, matchRepeat)
	}
	return matches
}

func SequenceMatch(password string) []match.Match {
	var matches []match.Match
	for i := 0; i < len(password); {
		j := i + 1
		var seq string
		var seqName string
		seqDirection := 0
		for seqCandidateName, seqCandidate := range SEQUENCES {
			iN := strings.Index(seqCandidate, string(password[i]))
			var jN int
			if j < len(password) {
				jN = strings.Index(seqCandidate, string(password[j]))
			} else {
				jN = -1
			}

			if iN > -1 && jN > -1 {
				direction := jN - iN
				if direction == 1 || direction == -1 {
					seq = seqCandidate
					seqName = seqCandidateName
					seqDirection = direction
					break
				}
			}

		}

		if seq != "" {
			for {
				var prevN, curN int
				if j < len(password) {
					prevChar, curChar := password[j-1], password[j]
					prevN, curN = strings.Index(seq, string(prevChar)), strings.Index(seq, string(curChar))
				}

				if j == len(password) || curN-prevN != seqDirection {
					if j-i > 2 {
						matchSequence := match.Match{
							Pattern:        "sequence",
							I:              i,
							J:              j - 1,
							Token:          password[i:j],
							DictionaryName: seqName,
						}

						matchSequence.Entropy = entropy.SequenceEntropy(matchSequence, len(seq), (seqDirection == 1))
						matches = append(matches, matchSequence)
					}
					break
				} else {
					j += 1
				}

			}
		}
		i = j
	}
	return matches
}
