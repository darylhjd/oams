package to

// Ptr is a helper function to return a pointer to a given value.
func Ptr[T any](t T) *T {
	return &t
}
