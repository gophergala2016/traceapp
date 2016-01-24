package node

type Node struct {
	Name    string   `json:"name"`
	Imports []string `json:"imports"`
	Size    int      `json:"size"`
}
