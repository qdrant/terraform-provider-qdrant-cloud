package qdrant

import (
	"github.com/google/uuid"
)

// Generic function to create a pointer to any type
func newPointer[T any](value T) *T {
	return &value
}

// Generic function to dereference a pointer with a default-value fallback
func derefPointer[T any](ptr *T, defaults ...T) T {
	if ptr != nil {
		return *ptr
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	var empty T
	return empty
}

// uuidArrayAsStringArray converts an array of UUID to an array of strings
func uuidArrayAsStringArray(ptr []uuid.UUID) []string {
	result := []string{}
	for _, uuid := range ptr {
		result = append(result, uuid.String())
	}
	return result
}
