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
	var mapped []R

	for _, el := range input {
		mapped = append(mapped, transform(el))
	}

	return mapped
}

func IsUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}
