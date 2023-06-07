package errors

import (
	"github.com/kyverno/playground/backend/pkg/engine"
)

type Error interface {
	Error() string
	Reason() interface{}
}

type genericError struct {
	error  string
	reason interface{}
}

func (e genericError) Error() string {
	return e.error
}

func (e genericError) Reason() interface{} {
	return e.reason
}

func PolicyValidations(errs []engine.PolicyValidation) Error {
	return genericError{
		error:  "policy validation failed",
		reason: errs,
	}
}
