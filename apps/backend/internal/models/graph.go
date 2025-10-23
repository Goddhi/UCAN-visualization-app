package models

import "time"

// GraphResponse contains graph data for visualization
type GraphResponse struct {
	Nodes []GraphNode `json:"nodes"`
	Edges []GraphEdge `json:"edges"`
	Chain ChainInfo   `json:"chain"`
}

// GraphNode represents a principal in the graph
type GraphNode struct {
	ID       string                 `json:"id"`
	Label    string                 `json:"label"`
	Type     string                 `json:"type"` // "root", "intermediate", "leaf"
	Level    int                    `json:"level"`
	Metadata map[string]interface{} `json:"metadata"`
}

// GraphEdge represents a delegation relationship
type GraphEdge struct {
	Source     string                 `json:"source"`
	Target     string                 `json:"target"`
	Capability CapabilityInfo         `json:"capability"`
	Valid      bool                   `json:"valid"`
	Label      string                 `json:"label"`
	Level      int                    `json:"level"`
	Type       string                 `json:"type"` // delegation, invocation, proof
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}
// InvocationGraphResponse contains graph data for invocation visualization
type InvocationGraphResponse struct {
	Nodes        []GraphNode          `json:"nodes"`
	Edges        []GraphEdge          `json:"edges"`
	Chain        ChainInfo            `json:"chain"`
	Invocation   *InvocationResponse  `json:"invocation"`
	IsInvocation bool                 `json:"isInvocation"`
}



// PrincipalInfo represents a DID in the chain
type PrincipalInfo struct {
	DID   string   `json:"did"`
	Role  string   `json:"role"` // "root", "intermediate", "leaf"
	Level int      `json:"level"`
	CIDs  []string `json:"cids"`
}

// TimelineEvent represents a temporal event in the delegation chain
type TimelineEvent struct {
	Type      string    `json:"type"`      // "issued", "expires"
	Time      time.Time `json:"time"`
	CID       string    `json:"cid"`
	Level     int       `json:"level"`
	Principal string    `json:"principal"`
}