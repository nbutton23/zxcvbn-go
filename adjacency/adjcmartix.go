package adjacency
import (
	"log"
	"encoding/json"
	"io/ioutil"
//	"fmt"
)


type AdjacencyGraph struct {
	Graph map[string][6]string
}
type AdjacencyGraphs struct {
	Qwerty AdjacencyGraph
	Dvorak AdjacencyGraph
	Keypad AdjacencyGraph
	MacKeypad AdjacencyGraph
}

var AdjacencyGph AdjacencyGraphs;
func init(){
	//todo get currentloc so that i don't have to know the whole path
	log.SetFlags(log.Lshortfile)
	qwerty := buildQwerty()
	dvorak := buildDvorak()
	keyPad := buildKeypad()
	macKeypad := buildMacKeypad()

	AdjacencyGph = AdjacencyGraphs{Qwerty:qwerty, Dvorak:dvorak, Keypad:keyPad, MacKeypad:macKeypad}
}

func buildQwerty() AdjacencyGraph {
	return getAdjancencyGraphFromFile("/Users/nbutton/workspace/src/zxcvbn-go/adjacency/Qwerty.json")
}
func buildDvorak() AdjacencyGraph {
	return getAdjancencyGraphFromFile("/Users/nbutton/workspace/src/zxcvbn-go/adjacency/Dvorak.json")
}
func buildKeypad() AdjacencyGraph {
	return getAdjancencyGraphFromFile("/Users/nbutton/workspace/src/zxcvbn-go/adjacency/Keypad.json")
}
func buildMacKeypad() AdjacencyGraph {
	return getAdjancencyGraphFromFile("/Users/nbutton/workspace/src/zxcvbn-go/adjacency/MacKeypad.json")
}

func getAdjancencyGraphFromFile(filePath string) AdjacencyGraph {
	data, err := ioutil.ReadFile(filePath)

	if err != nil {
		log.Fatal(err)
	}


	var graph AdjacencyGraph;
	err = json.Unmarshal(data, &graph)
	if err != nil {
		log.Fatal(err)
	}

	return graph
}

//on qwerty, 'g' has degree 6, being adjacent to 'ftyhbv'. '\' has degree 1.
//this calculates the average over all keys.
//TODO double check that i ported this correctly scoring.coffee ln 5
func (adjGrp AdjacencyGraph) CalculateAvgDegree() (float32) {
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

	return avg/count
}

