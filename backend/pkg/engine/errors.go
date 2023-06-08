package engine

type PolicyViolationError struct {
	Violations []PolicyValidation
}

func (e PolicyViolationError) Error() string {
	return "policy validation failed"
}
