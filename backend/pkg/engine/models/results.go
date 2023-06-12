package models

type Results struct {
	Mutation          []Response `json:"mutation"`
	ImageVerification []Response `json:"imageVerification"`
	Validation        []Response `json:"validation"`
	Generation        []Response `json:"generation"`
}
