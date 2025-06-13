package models

type Results struct {
	Mutation          []Response `json:"mutation"`
	ImageVerification []Response `json:"imageVerification"`
	Deletion          []Response `json:"deletion"`
	Validation        []Response `json:"validation"`
	Generation        []Response `json:"generation"`
}
