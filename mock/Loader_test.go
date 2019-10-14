package mock

import "testing"

func TestSimpleUrlComparasionSuccess(t *testing.T) {

	// Arrange

	lhs := "temp/path/to"
	rhs := "temp/path/to"

	// Act

	comparasionResult := CompareURLPath(lhs, rhs)

	// Assert

	if comparasionResult {
		return
	}

	t.Error(
		"Arrange: ", []string{lhs, rhs},
		"Act: ", "CompareURLPath",
		"Assert: ", comparasionResult,
		"Awaiting: ", !comparasionResult,
	)
}

func TestSimpleUrlComparasionFailure(t *testing.T) {
	// Arrange

	lhs := "temp/path/to"
	rhs := "to/temp/path"

	// Act

	comparasionResult := CompareURLPath(lhs, rhs)

	// Assert

	if !comparasionResult {
		return
	}

	t.Error(
		"Arrange: ", []string{lhs, rhs},
		"Act: ", "CompareURLPath",
		"Assert: ", comparasionResult,
		"Awaiting: ", !comparasionResult,
	)
}

func TestPathPatternMatchingWorkSuccess(t *testing.T) {
	// Arrange

	lhs := "/request/url/1"
	rhs := "/request/url/{id}"

	// Act

	comparasionResult := CompareURLPath(lhs, rhs)

	// Assert

	if comparasionResult {
		return
	}

	t.Error(
		"Arrange: ", []string{lhs, rhs},
		"Act: ", "CompareURLPath",
		"Assert: ", comparasionResult,
		"Awaiting: ", !comparasionResult,
	)
}

func TestPathWithDifferentNumberFailure(t *testing.T) {
	// Arrange

	lhs := "temp/path/to"
	rhs := "to/temp"

	// Act

	comparasionResult := CompareURLPath(lhs, rhs)

	// Assert

	if !comparasionResult {
		return
	}

	t.Error(
		"Arrange: ", []string{lhs, rhs},
		"Act: ", "CompareURLPath",
		"Assert: ", comparasionResult,
		"Awaiting: ", !comparasionResult,
	)
}

func TestSamePatternInPathAndInParamsMatchSuccess(t *testing.T) {

	// Arrange

	lhs := "temp/{id}/to?foo={foo}&bar={bar}"
	rhs := "temp/{id}/to?foo={foo}&bar={bar}"

	// Act

	comparasionResult := CompareURLPath(lhs, rhs)

	// Assert

	if comparasionResult {
		return
	}

	t.Error(
		"Arrange: ", []string{lhs, rhs},
		"Act: ", "CompareURLPath",
		"Assert: ", comparasionResult,
		"Awaiting: ", !comparasionResult,
	)
}

func TestDifferentPatternInPathAndInParamsMatchFailure(t *testing.T) {

	// Arrange

	lhs := "temp/{id}/to?foo={foo}"
	rhs := "temp/{id}/to?foo={foo}&bar={bar}"

	// Act

	comparasionResult := CompareURLPath(lhs, rhs)

	// Assert

	if !comparasionResult {
		return
	}

	t.Error(
		"Arrange: ", []string{lhs, rhs},
		"Act: ", "CompareURLPath",
		"Assert: ", comparasionResult,
		"Awaiting: ", !comparasionResult,
	)
}
