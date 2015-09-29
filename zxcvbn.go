package main

import (
//	"zxcvbn-go/adjacency"
	"fmt"
	"zxcvbn-go/frequency"
)

func main(){
	fmt.Println("Start")
//	fmt.Println(adjacency.AdjacencyGph)

	fmt.Println(len(frequency.FreqLists.Passwords))
}