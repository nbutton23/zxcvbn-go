package frequency
import (
	"io/ioutil"
	"log"
	"encoding/json"
)

type TempList struct {
	List []string
}
type FrequencyLists struct {
	MaleNames []string
	FemaleNames []string
	Surnames []string
	Passwords []string
	English []string
}

var FreqLists FrequencyLists
func init() {
	maleNames := GetStringListFromFile("/Users/nbutton/workspace/src/zxcvbn-go/frequency/MaleNames.json")
	femaleNames := GetStringListFromFile("/Users/nbutton/workspace/src/zxcvbn-go/frequency/FemaleNames.json")
	surnames := GetStringListFromFile("/Users/nbutton/workspace/src/zxcvbn-go/frequency/Surnames.json")
	passwords := GetStringListFromFile("/Users/nbutton/workspace/src/zxcvbn-go/frequency/Passwords.json")
	english := GetStringListFromFile("/Users/nbutton/workspace/src/zxcvbn-go/frequency/English.json")




	FreqLists = FrequencyLists{MaleNames:maleNames, FemaleNames:femaleNames, Surnames:surnames, Passwords:passwords, English:english}
}

func GetStringListFromFile(filePath string) []string {
	data, err := ioutil.ReadFile(filePath)

	if err != nil {
		log.Fatal(err)
	}


	var templist TempList;
	err = json.Unmarshal(data, &templist)
	if err != nil {
		log.Fatal(err)
	}
	return templist.List
}