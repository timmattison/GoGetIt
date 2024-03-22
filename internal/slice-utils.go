package internal

import (
	"slices"
)

// Difference returns a new slice containing the elements of source that are not in toRemove.
func Difference(source []string, toRemove []string) []string {
	var diff []string

	for _, value := range source {
		if slices.Contains(toRemove, value) {
			continue
		}

		diff = append(diff, value)
	}

	return diff
}

// Dedupe returns a new slice containing the unique elements of input.
func Dedupe(input []string) []string {
	var output []string

	for _, href := range input {
		if slices.Contains(output, href) {
			continue
		}

		output = append(output, href)
	}

	return output
}
