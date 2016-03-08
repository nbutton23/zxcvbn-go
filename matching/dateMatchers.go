package matching

import (
	"github.com/nbutton23/zxcvbn-go/match"
	"strings"
	"regexp"
	"strconv"
	"github.com/nbutton23/zxcvbn-go/entropy"
)



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
func dateSepMatcher(password string) []match.Match {
	dateMatches := dateSepMatchHelper(password)

	var matches []match.Match
	for _, dateMatch := range dateMatches {
		match := match.Match{
			I:dateMatch.I,
			J:dateMatch.J,
			Entropy:entropy.DateEntropy(dateMatch),
			DictionaryName:"date_match",
			Token:dateMatch.Token,
		}

		matches = append(matches, match)
	}

	return matches
}
func dateSepMatchHelper(password string) []match.DateMatch {

	var matches []match.DateMatch

	matcher := regexp.MustCompile(DATE_RX_YEAR_SUFFIX)
	for _, v := range matcher.FindAllString(password, len(password)) {
		splitV := matcher.FindAllStringSubmatch(v, len(v))
		i := strings.Index(password, v)
		j := i + len(v)
		day, _ := strconv.ParseInt(splitV[0][4], 10, 16)
		month, _ := strconv.ParseInt(splitV[0][2], 10, 16)
		year, _ := strconv.ParseInt(splitV[0][6], 10, 16)
		match := match.DateMatch{Day: day, Month: month, Year: year, Separator: splitV[0][5], I: i, J: j, Token:password[i:j]}
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
		match := match.DateMatch{Day: day, Month: month, Year: year, Separator: splitV[0][5], I: i, J: j, Token:password[i:j]}
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
