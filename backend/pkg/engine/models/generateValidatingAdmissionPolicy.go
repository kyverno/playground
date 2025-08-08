package models

type GenerateValidatingAdmissionPolicy struct {
	Enabled bool `json:"enabled"`
}

type GenerateMutatingAdmissionPolicy struct {
	Enabled bool `json:"enabled"`
}
