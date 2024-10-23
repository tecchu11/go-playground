package collection

// SMap maps slice S to slice []T.
func SMap[S ~[]E, E, T any](s S, fn func(E) T) []T {
	result := make([]T, len(s))
	for i, v := range s {
		result[i] = fn(v)
	}
	return result
}
