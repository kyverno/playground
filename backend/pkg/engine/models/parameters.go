package models

type Parameters struct {
	Kubernetes Kubernetes             `json:"kubernetes"`
	Context    Context                `json:"context"`
	Variables  map[string]interface{} `json:"variables"`
	Flags      Flags                  `json:"flags"`
}
