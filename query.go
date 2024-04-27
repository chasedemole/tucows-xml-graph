package main

import "encoding/json"

type Queries struct {
	Queries []Query `json:"queries"`
}

type Query struct {
	Paths    *AllPathsQuery     `json:"paths,omitempty"`
	Cheapest *CheapestPathQuery `json:"cheapest,omitempty"`
}

type AllPathsQuery struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type CheapestPathQuery struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

// Function to unmarshal the JSON and separate the queries
func unmarshalQueries(data []byte) ([]AllPathsQuery, []CheapestPathQuery, error) {
	// Unmarshal the data into the Queries struct
	var queries Queries
	err := json.Unmarshal(data, &queries)
	if err != nil {
		return nil, nil, err
	}

	// Separate the queries into respective slices
	var allPathQueries []AllPathsQuery
	var cheapestQueries []CheapestPathQuery

	for _, query := range queries.Queries {
		if query.Paths != nil {
			allPathQueries = append(allPathQueries, *query.Paths)
		}
		if query.Cheapest != nil {
			cheapestQueries = append(cheapestQueries, *query.Cheapest)
		}
	}

	return allPathQueries, cheapestQueries, nil
}
