package models

// GraphResponse contains graph data for visualization
type GraphResponse struct {
	Nodes []GraphNode `json:"nodes"`
	Edges []GraphEdge `json:"edges"`
}

// GraphNode represents a principal in the graph
type GraphNode struct {
	ID       string                 `json:"id"`
	Label    string                 `json:"label"`
	Type     string                 `json:"type"` // "root", "intermediate", "leaf"
	Metadata map[string]interface{} `json:"metadata"`
}

// GraphEdge represents a delegation relationship
type GraphEdge struct {
	Source     string         `json:"source"`
	Target     string         `json:"target"`
	Capability CapabilityInfo `json:"capability"`
	Valid      bool           `json:"valid"`
	Label      string         `json:"label"`
}