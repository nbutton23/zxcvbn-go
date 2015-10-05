package main

import (
	"fmt"
	"zxcvbn-go/matching"
)

func main() {
	fmt.Println("Start")
	fmt.Println(matching.SpatialMatch("qw@!andghjandfTandftg"))
}