package models

type Flags struct {
	Exceptions Exceptions `json:"exceptions"`
	Cosign     Cosign     `json:"cosign"`
}
