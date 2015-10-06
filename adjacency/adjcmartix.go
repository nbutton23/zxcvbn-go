package adjacency
import (
	"log"
	"encoding/json"
	"io/ioutil"
//	"fmt"
	"path/filepath"
)


type AdjacencyGraph struct {
	Graph map[string][]string
	averageDegree float64
	Name string
}




func buildQwerty() AdjacencyGraph {
	filePath, _ := filepath.Abs("Qwerty.json")
	return GetAdjancencyGraphFromFile(filePath, "qwerty")
}
func buildDvorak() AdjacencyGraph {
	filePath, _ := filepath.Abs("Dvorak.json")
	return GetAdjancencyGraphFromFile(filePath, "dvorak")
}
func buildKeypad() AdjacencyGraph {
	filePath, _ := filepath.Abs("Keypad.json")
	return GetAdjancencyGraphFromFile(filePath, "keypad")
}
func buildMacKeypad() AdjacencyGraph {
	filePath, _ := filepath.Abs("MacKeypad.json")
	return GetAdjancencyGraphFromFile(filePath, "mac_keypad")
}

func GetAdjancencyGraphFromFile(filePath string, name string) AdjacencyGraph {
	data, err := ioutil.ReadFile(filePath)

	if err != nil {
		log.Fatal(err)
	}


	var graph AdjacencyGraph;
	err = json.Unmarshal(data, &graph)
	if err != nil {
		log.Fatal(err)
	}
	graph.Name = name
	return graph
}

//on qwerty, 'g' has degree 6, being adjacent to 'ftyhbv'. '\' has degree 1.
//this calculates the average over all keys.
//TODO double check that i ported this correctly scoring.coffee ln 5
func (adjGrp AdjacencyGraph) CalculateAvgDegree() (float64) {
	if adjGrp.averageDegree != float64(0) {
		return adjGrp.averageDegree
	}
	var avg float64
	var count float64
	for _, value := range adjGrp.Graph {

		for _, char := range value {
			if char != "" || char != " " {
				avg += float64(len(char))
				count++
			}
		}

	}

	adjGrp.averageDegree = avg/count

	return adjGrp.averageDegree
}

