package main

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"log"
	"os"
)

func main() {
	fileName := "./test.xml"
	xmlFile, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}

	decoder := xml.NewTokenDecoder(xml.NewDecoder(xmlFile))

	xmlGraph := NewXMLGraph(decoder)

	fmt.Println(xmlGraph)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		jsonData := scanner.Text()
		fmt.Printf("JSON data: %v\n", jsonData)
		graph := BuildGraphGraphFromXML(xmlGraph.Nodes, xmlGraph.Edges)

		allPathQueries, cheapestQueries, err := unmarshalQueries([]byte(jsonData))
		if err != nil {
			log.Fatalf("Invalid json: %v", err)
		}

		for _, query := range cheapestQueries {
			fmt.Println(CheapestPath(graph, query))
		}

		for _, query := range allPathQueries {
			fmt.Println(AllPaths(graph, query))
		}
	}
	if scanner.Err() != nil {
		// Handle error.
		log.Printf("Error parsing STDIN: %v", err)
	}
}
