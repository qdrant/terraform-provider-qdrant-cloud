package qdrant

// Generic function to create a pointer to any type.
func newPointer[T any](value T) *T {
	return &value
}

// Generic function to dereference a pointer with a default-value fallback.
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
