package frequency
import (
	"io/ioutil"
	"log"
	"encoding/json"
)

type listWrapper struct {
	List []string
}

func GetStringListFromFile(filePath string) []string {
	data, err := ioutil.ReadFile(filePath)

	if err != nil {
		log.Fatal(err)
	}


	var tempList listWrapper;
	err = json.Unmarshal(data, &tempList)
	if err != nil {
		log.Fatal(err)
	}
	return tempList.List
}