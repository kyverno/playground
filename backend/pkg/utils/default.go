package utils

func Fallback[T comparable](value, defaultValue T) T {
	var zero T

	if value == zero {
		return defaultValue
	}
	return value
}
