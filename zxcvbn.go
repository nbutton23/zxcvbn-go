package main

import (
	"zxcvbn-go/adjacency"
	"fmt"
	"zxcvbn-go/utils/math"
	"zxcvbn-go/matching"
)

func main() {
	fmt.Println("Start")
	fmt.Println(adjacency.AdjacencyGph.Qwerty.CalculateAvgDegree())
	fmt.Println(adjacency.AdjacencyGph.Dvorak.CalculateAvgDegree())
	fmt.Println(adjacency.AdjacencyGph.Keypad.CalculateAvgDegree())
	fmt.Println(adjacency.AdjacencyGph.MacKeypad.CalculateAvgDegree())


	fmt.Println(math.NChoseK(100, 2))

	fmt.Println(matching.DateSepMatch("1991-09-11jibjab11.9.1991"))
}