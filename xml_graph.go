package main

import (
	"encoding/xml"
	"io"
	"log"
)

type XMLGraph struct {
	ID      string    `xml:"id"`
	Name    string    `xml:"name"`
	Nodes   []XMLNode `xml:"nodes>node"`
	Edges   []XMLEdge `xml:"edges>node"`
	NodeIDs map[string]struct{}
}

type XMLNode struct {
	ID   string `xml:"id"`
	Name string `xml:"name"`
}

type XMLEdge struct {
	ID        string `xml:"id"`
	From      string `xml:"from"`
	To        string `xml:"to"`
	CostCents int    `xml:"cost"`
}

func (g *XMLGraph) ProcessNodes(decoder *xml.Decoder) {
	for {
		token, err := decoder.Token()
		if err != nil {
			// TODO: EOF shouldn't end on <nodes>
			if err == io.EOF {
				return
			}
			log.Printf("Error fetching token: %v", err)
			return
		}

		var node XMLNode
		switch element := token.(type) {
		case xml.StartElement:
			if element.Name.Local == "node" {
				if err := decoder.DecodeElement(&node, &element); err != nil {
					log.Printf("error unmarshalling node into struct: %v", err)
				}

				if !node.Valid() {
					log.Println("invalid node")
					break
				}

				if _, ok := g.NodeIDs[node.ID]; ok {
					log.Fatalf("invalid graph, duplicate node IDs: %s", node.ID)
				}

				g.NodeIDs[node.ID] = struct{}{}
				g.Nodes = append(g.Nodes, node)
			} else {
				if err := decoder.Skip(); err != nil {
					log.Panicf("invalid XML, %s missing closing tag", element.Name.Local)
				}
			}
		case xml.EndElement:
			if element.Name.Local == "nodes" {
				return
			}
		}
	}
}

func (g *XMLGraph) ProcessEdges(decoder *xml.Decoder) {
	for {
		token, err := decoder.Token()
		if err != nil {
			// TODO: EOF shouldn't end on <nodes>
			if err == io.EOF {
				return
			}
			log.Printf("Error fetching token: %v", err)
			return
		}

		var edge XMLEdge
		switch element := token.(type) {
		case xml.StartElement:
			if element.Name.Local == "node" {
				if err := decoder.DecodeElement(&edge, &element); err != nil {
					log.Printf("error unmarshalling node into struct: %v", err)
				}

				if !edge.Valid() {
					log.Println("invalid edge")
					break
				}

				g.Edges = append(g.Edges, edge)
			} else {
				if err := decoder.Skip(); err != nil {
					log.Panicf("invalid XML, %s missing closing tag", element.Name.Local)
				}
			}
		case xml.EndElement:
			if element.Name.Local == "edges" {
				return
			}
		}
	}
}

func (g XMLGraph) Valid() bool {
	if g.ID == "" {
		return false
	} else if g.Name == "" {
		return false
	} else if len(g.Nodes) == 0 {
		return false
	}

	return true
}

func (e XMLEdge) Valid() bool {
	if e.To == "" {
		return false
	} else if e.From == "" {
		return false
	} else if e.CostCents < 0 {
		return false
	}

	return true
}

func (n XMLNode) Valid() bool {
	if n.ID == "" {
		return false
	} else if n.Name == "" {
		return false
	}

	return true
}

func NewXMLGraph(decoder *xml.Decoder) *XMLGraph {
	graph := &XMLGraph{
		Nodes:   []XMLNode{},
		Edges:   []XMLEdge{},
		NodeIDs: make(map[string]struct{}),
	}

	for {
		token, err := decoder.Token()
		if err != nil {
			// TODO: EOF shouldn't end on <nodes>
			if err == io.EOF {
				return graph
			}
			log.Printf("Error fetching token: %v", err)
			return nil
		}

		if token == nil {
			break
		}

		switch element := token.(type) {
		case xml.StartElement:
			if element.Name.Local == "graph" {
				log.Println("found graph\n")
			} else if element.Name.Local == "id" {
				decoder.DecodeElement(&graph.ID, &element)
			} else if element.Name.Local == "name" {
				decoder.DecodeElement(&graph.Name, &element)
			} else if element.Name.Local == "nodes" {
				graph.ProcessNodes(decoder)
			} else if element.Name.Local == "edges" {
				graph.ProcessEdges(decoder)
			}
		}
	}

	if !graph.Valid() {
		return nil
	}

	return graph
}
