package matching
import (
	"strings"
	"github.com/bradfitz/slice"
	"regexp"
	"strconv"
)

var DICTIONARY_MATCHERS []func(password string) []Match

 const (
	 //TODO: Invalid regex for Golang since it has a \2
	 DATE_RX_YEAR_SUFFIX string = `((\d{1,2})(\s|-|\/|\\|_|\.)(\d{1,2})(\s|-|\/|\\|_|\.)(19\d{2}|200\d|201\d|\d{2}))`
	 DATE_RX_YEAR_PREFIX string = `((19\d{2}|200\d|201\d|\d{2})(\s|-|/|\\|_|\.)(\d{1,2})(\s|-|/|\\|_|\.)(\d{1,2}))`
 )
type Match struct {
	Pattern string
	I, J int
	Token string
	MatchedWord string
	Rank int
	DictionaryName string
}

type DateMatch struct {
	Pattern string
	I, J int
	Token string
	Separator string
	Day, Month, Year int64

}

func Omnimatch(password string, userInputs []string) []Match  {
	var rankedUserInputsDir map[string]int

	for i, v := range userInputs {
		rankedUserInputsDir[strings.ToLower(v)] = i+1
	}
	userInputMatcher := buildDictMatcher("user_inputs", rankedUserInputsDir)
	matches := userInputMatcher(password)

	for _, matcher := range DICTIONARY_MATCHERS {
		mtemp := matcher(password)
		for _,v:= range mtemp {
			matches = append(matches, v)
		}
	}
	slice.Sort(matches,func(i, j int)bool{
		//TODO fix this
		return false;
	})
	return matches
}




func buildDictMatcher(dictName string, rankedDict map[string]int) func(password string) []Match {
	return func (password string) []Match{
		matches := dictionaryMatch(password, rankedDict)
		for _, v := range matches {
			v.DictionaryName = dictName
		}
		return matches
	}

}

func dictionaryMatch(password string, rankedDict map[string]int) []Match{
	length := len(password)
	var results []Match
	pwLower := strings.ToLower(password)

	for i :=0; i<length; i++ {
		for j := i; j<length; j++ {
			word := pwLower[i:j+1]
			if val, ok := rankedDict[word]; ok {
				results = append(results, Match{Pattern:"dictionary",
				I:i,
				J:j,
				Token:password[i:j+1],
				MatchedWord:word,
				Rank:val})
			}
		}
	}

	return results
}

func checkDate(day, month, year int64)( bool, int64, int64, int64){
	if (12 <= month && month <= 31) && day <= 12 {
		day, month = month, day
	}

	if day > 31 || month > 12 {
		return false, 0, 0, 0
	}

	if !(1900 <= year && year <=2019) {
		return false, 0, 0, 0
	}

	return true, day, month, year
}

func DateSepMatch(password string) []DateMatch {

	var matches []DateMatch


	matcher := regexp.MustCompile(DATE_RX_YEAR_SUFFIX)
	for _, v := range matcher.FindAllString(password,len(password)) {
		splitV := matcher.FindAllStringSubmatch(v, len(v))
		i := strings.Index(password,v)
		j := i+len(v)
		day, _ := strconv.ParseInt(splitV[0][4],10,16)
		month, _ := strconv.ParseInt(splitV[0][2], 10, 16)
		year, _ := strconv.ParseInt(splitV[0][6], 10, 16)
		match := DateMatch{Day:day, Month:month, Year:year, Separator:splitV[0][5], I:i, J:j }
		matches = append(matches, match)
	}


	matcher = regexp.MustCompile(DATE_RX_YEAR_PREFIX)
	for _, v := range matcher.FindAllString(password,len(password)) {
		splitV := matcher.FindAllStringSubmatch(v, len(v))
		i := strings.Index(password,v)
		j := i+len(v)
		day, _ := strconv.ParseInt(splitV[0][4],10,16)
		month, _ := strconv.ParseInt(splitV[0][6], 10, 16)
		year, _ := strconv.ParseInt(splitV[0][2], 10, 16)
		match := DateMatch{Day:day, Month:month, Year:year, Separator:splitV[0][5], I:i, J:j }
		matches = append(matches, match)
	}

	var out []DateMatch
	for _, match := range matches {
		if valid, day, month, year := checkDate(match.Day, match.Month, match.Year); valid{
			match.Pattern = "date"
			match.Day = day
			match.Month = month
			match.Year = year
			out = append(out, match)
		}
	}
	return out

}