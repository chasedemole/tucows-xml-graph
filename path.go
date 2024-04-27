package main

import (
	"container/list"
	"fmt"
	"math"
)

type Connection struct {
	To        string
	From      string
	CostCents int
}

func BuildGraphGraphFromXML(xmlNodes []XMLNode, xmlEdges []XMLEdge) map[string][]Connection {
	graphNodes := make(map[string][]Connection)
	for _, node := range xmlNodes {
		connections := []Connection{}
		graphNodes[node.Name] = connections
	}

	for _, edge := range xmlEdges {
		connection := Connection{
			To:        edge.To,
			From:      edge.From,
			CostCents: edge.CostCents,
		}
		graphNodes[edge.From] = append(graphNodes[edge.From], connection)
	}

	return graphNodes
}

func dfs(graph map[string][]Connection, query CheapestPathQuery, visited map[string]bool, bestPath *[]string, currentPath []string, node string, minCost *int, currentCost int) {
	if visited[node] {
		return
	}

	currentPath = append(currentPath, node)
	visited[node] = true

	if node == query.End {
		// If the current cost is less than the minimum cost, update the best path
		if currentCost < *minCost {
			*minCost = currentCost
			*bestPath = append([]string(nil), currentPath...)
		}
	} else {
		// Traverse connected nodes
		for _, next := range graph[node] {
			dfs(graph, query, visited, bestPath, currentPath, next.To, minCost, currentCost+next.CostCents)
		}
	}

	// Backtrack
	currentPath = currentPath[:len(currentPath)-1]
	visited[node] = false
}

func AllPaths(graph map[string][]Connection, query AllPathsQuery) [][]string {
	var allPaths [][]string
	queue := list.New()
	queue.PushBack([]string{query.Start})

	for queue.Len() > 0 {
		currentPath := queue.Remove(queue.Front()).([]string)
		lastNode := currentPath[len(currentPath)-1]

		if lastNode == query.End {
			allPaths = append(allPaths, append([]string(nil), currentPath...))
		} else {
			for _, connection := range graph[lastNode] {
				newPath := append([]string(nil), currentPath...)
				newPath = append(newPath, connection.To)
				queue.PushBack(newPath)
			}
		}
	}

	return allPaths
}

func CheapestPath(graph map[string][]Connection, query CheapestPathQuery) ([]string, int) {
	minCost := math.MaxInt
	bestPath := []string{}
	currentPath := []string{}
	visited := make(map[string]bool)

	dfs(graph, query, visited, &bestPath, currentPath, query.Start, &minCost, 0)

	if minCost == math.MaxInt {
		fmt.Printf("No path from %s to %s\n", query.Start, query.End)
		return []string{}, -1
	}
	return bestPath, minCost
}
