package models

type ImageData struct {
	Image         string      `json:"image"`
	ResolvedImage string      `json:"resolvedImage"`
	Registry      string      `json:"registry"`
	Repository    string      `json:"repository"`
	Identifier    string      `json:"identifier"`
	Manifest      interface{} `json:"manifest"`
	ConfigData    interface{} `json:"configData"`
}
