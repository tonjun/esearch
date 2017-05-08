package esearch

import (
	"encoding/json"
)

// Result represents a search result response
type Result struct {
	Took int64 `json:"took"`
	Hits *Hits `json:"hits"`
}

// Hits is the list of search hits
type Hits struct {
	Total    int64    `json:"total"`
	MaxScore *float64 `json:"max_score"`
	Hits     []*Hit   `json:"hits"`
}

// Hit is a single hit
type Hit struct {
	Index  string           `json:"_index"`
	Type   string           `json:"_type"`
	ID     string           `json:"_id"`
	Score  *float64         `json:"_score"`
	Source *json.RawMessage `json:"_source"`
}
