package qdrant

// Generic function to create a pointer to any type.
func newPointer[T any](value T) *T {
	return &value
}
