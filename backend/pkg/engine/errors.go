package engine

import (
	"github.com/kyverno/playground/backend/pkg/engine/models"
)

type PolicyViolationError struct {
	Violations []models.PolicyValidation
}

func (e PolicyViolationError) Error() string {
	return "policy validation failed"
}
