package scoring
import (
	"zxcvbn-go/match"
	"unicode"
	"fmt"
	"math"
	"sort"
	"regexp"
	"zxcvbn-go/utils/math"
	"zxcvbn-go/matching"
)


const (
	START_UPPER string = `^[A-Z][^A-Z]+$`
	END_UPPER string = `^[^A-Z]+[A-Z]$'`
	ALL_UPPER string = `^[A-Z]+$`

//for a hash function like bcrypt/scrypt/PBKDF2, 10ms per guess is a safe lower bound.
//(usually a guess would take longer -- this assumes fast hardware and a small work factor.)
//adjust for your site accordingly if you use another hash function, possibly by
//several orders of magnitude!
	SINGLE_GUESS float64 = 0.010
	NUM_ATTACKERS float64 = 100 //Cores used to make guesses
	SECONDS_PER_GUESS float64 = SINGLE_GUESS / NUM_ATTACKERS


)
type MinEntropyMatch struct {
	Password         string
	Entropy          float64
	MatchSequence    []match.Match //TODO ?
	CrackTime        float64
	CrackTimeDisplay string
	Score            int
	CalcTime         float64
}

/*
Returns minimum entropy TODO

    Takes a list of overlapping matches, returns the non-overlapping sublist with
    minimum entropy. O(nm) dp alg for length-n password with m candidate matches.
 */
func MinimumEntropyMatchSequence(password string, matches []match.Match) MinEntropyMatch {
	bruteforceCardinality := float64(calcBruteforceCardinality(password))
	upToK := make([]float64, len(password))
	backPointers := make([]match.Match, len(password))

	for k := 0; k < len(password); k++ {
		upToK[k] = get(upToK, k - 1) + math.Log2(bruteforceCardinality)

		for _, match := range matches {
			if match.J != k {
				continue
			}

			i, j := match.I, match.J
			//			see if best entropy up to i-1 + entropy of match is less that current min at j
			upTo := get(upToK, i - 1)
			calculatedEntropy := calcEntropy(match)
			match.Entropy = calculatedEntropy
			candidateEntropy := upTo + calculatedEntropy

			if candidateEntropy < upToK[j] {
				upToK[j] = candidateEntropy
				match.Entropy = candidateEntropy
				backPointers[j] = match
			}
		}
	}

	//		walk backwards and decode the best sequence
	var matchSequence []match.Match
	passwordLen := len(password)
	passwordLen--
	for k := passwordLen; k >= 0; {
		match := backPointers[k]
		if match.Pattern != "" {
			matchSequence = append(matchSequence, match)
			k = match.I - 1

		}  else {
			k--
		}

	}
	sort.Sort(match.Matches(matchSequence))

	makeBruteForecMatch := func(i, j int) match.Match {
		return match.Match{Pattern:"bruteforce",
			I:i,
			J:j,
			Token:password[i:j + 1],
			Entropy:math.Log2(math.Pow(bruteforceCardinality, float64(j - i)))}

	}

	k := 0
	var matchSequenceCopy []match.Match
	for _, match := range matchSequence {
		i, j := match.I, match.J
		if i - k > 0 {
			matchSequenceCopy = append(matchSequenceCopy, makeBruteForecMatch(k, i - 1))
		}
		k = j + 1
		matchSequenceCopy = append(matchSequenceCopy, match)
	}

	if k < len(password) {
		matchSequenceCopy = append(matchSequenceCopy, makeBruteForecMatch(k, len(password) - 1))
	}
	var minEntropy float64
	if len(password) == 0 {
		minEntropy = float64(0)
	} else {
		minEntropy = upToK[len(password) - 1 ]
	}

	crackTime := roundToXDigits(entropyToCrackTime(minEntropy), 3)
	return MinEntropyMatch{Password:password,
		Entropy:roundToXDigits(minEntropy, 3),
		MatchSequence:matchSequenceCopy,
		CrackTime:crackTime,
		CrackTimeDisplay:displayTime(crackTime),
		Score:crackTimeToScore(crackTime)}

}
func get(a []float64, i int) float64 {
	if i < 0 || i >= len(a) {
		return float64(0)
	}

	return a[i]
}
func calcBruteforceCardinality(password string) float64 {
	lower, upper, digits, symbols := float64(0), float64(0), float64(0), float64(0)

	for _, char := range password {
		if unicode.IsLower(char) {
			lower = float64(26)
		} else if unicode.IsDigit(char) {
			digits = float64(10)
		} else if unicode.IsUpper(char) {
			upper = float64(26)
		} else {
			symbols = float64(33)
		}
	}

	cardinality := lower + upper + digits + symbols
	return cardinality
}

func calcEntropy(match match.Match) float64 {
	if match.Entropy > float64(0) {
		return match.Entropy
	}

	var entropy float64
	if match.Pattern == "dictionary" {
		entropy = dictionaryEntropy(match)
	} else if match.Pattern == "spatial" {
		entropy = spatialEntropy(match)
	} else if match.Pattern == "repeat" {
		entropy = repeatEntropy(match)
	}

	match.Entropy = entropy
	//TODO finish implement this. . . this looks to be the meat and potatoes of the calculation
	return match.Entropy
}

func dictionaryEntropy(match match.Match) float64 {
	baseEntropy := math.Log2(match.Rank)
	upperCaseEntropy := extraUpperCaseEntropy(match)
	//TODO: L33t
	return baseEntropy + upperCaseEntropy
}

func spatialEntropy(match match.Match) float64 {
	var s, d float64
	if match.DictionaryName == "qwerty" || match.DictionaryName == "dvorak" {
		s = float64(matching.KEYBOARD_STARTING_POSITIONS)
		d = matching.KEYBOARD_AVG_DEGREE
	} else {
		s = float64(matching.KEYPAD_STARTING_POSITIONS)
		d = matching.KEYPAD_AVG_DEGREE
	}

	possibilities := float64(0)

	lenght := float64(len(match.Token))
	t := match.Turns

	//TODO: Should this be <= or just < ?
	//Estimate the number of possible patterns w/ lenght L or less with t turns or less
	for i := float64(2); i <= lenght + 1; i++ {
		possibleTurns := math.Min(float64(t), i - 1)
		for j := float64(1); j <= possibleTurns + 1; j++ {
			x := zxcvbn_math.NChoseK(i - 1, j - 1) * s * math.Pow(d, j)
			possibilities += x
		}
	}

	entropy := math.Log2(possibilities)

	//add extra entropu for shifted keys. ( % instead of 5 A instead of a)
	//Math is similar to extra entropy for uppercase letters in dictionary matches.

	if S := float64(match.ShiftedCount); S > float64(0) {
		possibilities = float64(0)
		U := lenght - S

		for i := float64(0); i < math.Min(S, U) + 1; i++ {
			possibilities += zxcvbn_math.NChoseK(S + U, i)
		}

		entropy += math.Log2(possibilities)
	}

	return entropy
}
func extraUpperCaseEntropy(match match.Match) float64 {
	word := match.Token

	allLower := true

	for _, char := range word {
		if unicode.IsUpper(char) {
			allLower = false
			break
		}
	}
	if allLower {
		return float64(0)
	}

	//a capitalized word is the most common capitalization scheme,
	//so it only doubles the search space (uncapitalized + capitalized): 1 extra bit of entropy.
	//allcaps and end-capitalized are common enough too, underestimate as 1 extra bit to be safe.

	for _, regex := range []string{START_UPPER, END_UPPER, ALL_UPPER} {
		matcher := regexp.MustCompile(regex)

		if matcher.MatchString(word) {
			return float64(1)
		}
	}
	//Otherwise calculate the number of ways to capitalize U+L uppercase+lowercase letters with U uppercase letters or
	//less. Or, if there's more uppercase than lower (for e.g. PASSwORD), the number of ways to lowercase U+L letters
	//with L lowercase letters or less.

	countUpper, countLower := float64(0), float64(0)
	for _, char := range word {
		if unicode.IsUpper(char) {
			countUpper++
		} else if unicode.IsLower(char) {
			countLower++
		}
	}
	totalLenght := countLower + countUpper
	var possibililities float64

	for i := float64(0); i <= math.Min(countUpper, countLower); i++ {
		possibililities += float64(zxcvbn_math.NChoseK(totalLenght, i))
	}

	if possibililities < 1 {
		return float64(1)
	}

	return float64(math.Log2(possibililities))
}

func repeatEntropy(match match.Match) float64 {
	cardinality := calcBruteforceCardinality(match.Token)
	entropy := math.Log2(cardinality * float64(len(match.Token)))

	return entropy
}

func entropyToCrackTime(entropy float64) float64 {
	crackTime := (0.5 * math.Pow(float64(2), entropy)) * SECONDS_PER_GUESS

	return crackTime
}

func roundToXDigits(number float64, digits int) float64 {
	return zxcvbn_math.Round(number, .5, digits)
}

func displayTime(seconds float64) string {
	formater := "%.1f %s"
	minute := float64(60)
	hour := minute * float64(60)
	day := hour * float64(24)
	month := day * float64(31)
	year := month * float64(12)
	century := year * float64(100)

	if seconds < minute {
		return "instant"
	} else if seconds < hour {
		return fmt.Sprintf(formater, (1 + math.Ceil(seconds / minute)), "minutes")
	} else if seconds < day {
		return fmt.Sprintf(formater, (1 + math.Ceil(seconds / hour)), "hours")
	} else if seconds < month {
		return fmt.Sprintf(formater, (1 + math.Ceil(seconds / day)), "days")
	} else if seconds < year {
		return fmt.Sprintf(formater, (1 + math.Ceil(seconds / month)), "months")
	}else if seconds < century {
		return fmt.Sprintf(formater, (1 + math.Ceil(seconds / century)), "years")
	} else {
		return "centuries"
	}
}

func crackTimeToScore(seconds float64) int {
	if seconds < math.Pow(10, 2) {
		return 0
	} else if seconds < math.Pow(10, 4) {
		return 1
	} else if seconds < math.Pow(10, 6) {
		return 2
	} else if seconds < math.Pow(10, 8) {
		return 3
	}

	return 4
}