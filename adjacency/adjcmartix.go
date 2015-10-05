package adjacency
import (
	"log"
	"encoding/json"
	"io/ioutil"
//	"fmt"
	"path/filepath"
)


type AdjacencyGraph struct {
	Graph map[string][6]string
	averageDegree float32
	Name string
}


var AdjacencyGph []AdjacencyGraph;
func init(){
	//todo get currentloc so that i don't have to know the whole path
	log.SetFlags(log.Lshortfile)
	AdjacencyGph = append(AdjacencyGph, buildQwerty())
	AdjacencyGph = append(AdjacencyGph, buildDvorak())
	AdjacencyGph = append(AdjacencyGph, buildKeypad())
	AdjacencyGph = append(AdjacencyGph, buildMacKeypad())



}

func buildQwerty() AdjacencyGraph {
	filePath, _ := filepath.Abs("adjacency/Qwerty.json")
	return getAdjancencyGraphFromFile(filePath, "qwerty")
}
func buildDvorak() AdjacencyGraph {
	filePath, _ := filepath.Abs("adjacency/Dvorak.json")
	return getAdjancencyGraphFromFile(filePath, "dvorak")
}
func buildKeypad() AdjacencyGraph {
	filePath, _ := filepath.Abs("adjacency/Keypad.json")
	return getAdjancencyGraphFromFile(filePath, "keypad")
}
func buildMacKeypad() AdjacencyGraph {
	filePath, _ := filepath.Abs("adjacency/MacKeypad.json")
	return getAdjancencyGraphFromFile(filePath, "mac_keypad")
}

func getAdjancencyGraphFromFile(filePath string, name string) AdjacencyGraph {
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
func (adjGrp AdjacencyGraph) CalculateAvgDegree() (float32) {
	if adjGrp.averageDegree != float32(0) {
		return adjGrp.averageDegree
	}
	var avg float32
	var count float32
	for _, value := range adjGrp.Graph {

		for _, char := range value {
			if char != "" || char != " " {
				avg += float32(len(char))
				count++
			}
		}

	}

	adjGrp.averageDegree = avg/count

	return adjGrp.averageDegree
}

