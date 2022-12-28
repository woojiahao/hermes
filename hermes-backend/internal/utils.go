package internal

import "github.com/google/uuid"

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

func Map[T any, R any](input []T, transform func(T) R) []R {
	mapped := make([]R, 0)

	for _, el := range input {
		mapped = append(mapped, transform(el))
	}

	return mapped
}

func FilterMap[T any, R any](input []T, predicate func(T) bool, transform func(T) R) []R {
	return Map(Filter(input, predicate), transform)
}

func ForEach[T any](input []T, fn func(T)) {
	for _, el := range input {
		fn(el)
	}
}

func Flatten[T any](input [][]T) []T {
	var flattened []T

	for _, l := range input {
		flattened = append(flattened, l...)
	}

	return flattened
}

func IsUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

func ThisOrThat[T any](this, that T, condition bool) T {
	if condition {
		return this
	} else {
		return that
	}
}
