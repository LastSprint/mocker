package mock

import "testing"

// TestSimpleUrlComparasionSuccess tests that `path/path/path` equals to `path/path/path`
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

// TestSimpleUrlComparasionFailure test that `temp/path/to` NOT equals to `to/temp/path`
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

// TestPathPatternMatchingWorkSuccess test that `/request/url/1` equals to `request/url/{id}`
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
