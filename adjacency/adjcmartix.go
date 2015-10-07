package adjacency
import (
	"log"
	"encoding/json"
//	"fmt"
	"github.com/nbutton23/zxcvbn-go/data"
)


type AdjacencyGraph struct {
	Graph map[string][]string
	averageDegree float64
	Name string
}


var AdjacencyGph = make(map[string]AdjacencyGraph);
func init() {
	//todo get currentloc so that i don't have to know the whole path
	log.SetFlags(log.Lshortfile)
	AdjacencyGph["qwerty"] = buildQwerty()
	AdjacencyGph["dvorak"] = buildDvorak()
	AdjacencyGph["keypad"] = buildKeypad()
	AdjacencyGph["macKeypad"] = buildMacKeypad()
}

func buildQwerty() AdjacencyGraph {
	data, err := zxcvbn_data.Asset("data/Qwerty.json")
	if err != nil {
		panic("Can't find asset")
	}
	return GetAdjancencyGraphFromFile(data, "qwerty")
}
func buildDvorak() AdjacencyGraph {
	data, err := zxcvbn_data.Asset("data/Dvorak.json")
	if err != nil {
		panic("Can't find asset")
	}
	return GetAdjancencyGraphFromFile(data, "dvorak")
}
func buildKeypad() AdjacencyGraph {
	data, err := zxcvbn_data.Asset("data/Keypad.json")
	if err != nil {
		panic("Can't find asset")
	}
	return GetAdjancencyGraphFromFile(data, "keypad")
}
func buildMacKeypad() AdjacencyGraph {
	data, err := zxcvbn_data.Asset("data/MacKeypad.json")
	if err != nil {
		panic("Can't find asset")
	}
	return GetAdjancencyGraphFromFile(data, "mac_keypad")
}

func GetAdjancencyGraphFromFile(data []byte, name string) AdjacencyGraph {

	var graph AdjacencyGraph;
	err := json.Unmarshal(data, &graph)
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

