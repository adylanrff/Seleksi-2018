package main

import (
	"fmt"

	. "./foodcrawler"
)

const DirPath = "data/"

func main() {
	recordNumber := 75
	filename := DirPath + "output.json"
	foodList := Crawl(filename, "banana", recordNumber)

	fmt.Printf("Parsed %v data \n", len(foodList))
}
