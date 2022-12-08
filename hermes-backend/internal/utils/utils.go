package utils

func Contains[T comparable](arr []T, target T) bool {
	for _, el := range arr {
		if el == target {
			return true
		}
	}

	return false
}

func Filter[T any](arr []T, predicate func(T) bool) []T {
	var filtered []T

	for _, el := range arr {
		if predicate(el) {
			filtered = append(filtered, el)
		}
	}

	return filtered
}

func HasDuplicates[T comparable](arr []T) bool {
	var seen []T

	for _, el := range arr {
		if Contains(seen, el) {
			return true
		}

		seen = append(seen, el)
	}

	return false
}
