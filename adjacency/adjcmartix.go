package adjacency
import (
	"log"
	"encoding/json"
	"io/ioutil"
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

