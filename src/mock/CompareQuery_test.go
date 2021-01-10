package mock

import "testing"

func TestSimpleQueryCoparasionSuccess(t *testing.T) {

	// Arrange

	lhs := "load?data=1"
	rhs := "load?data=1"

	// Act

	result := compareQuery(lhs, rhs)

	// Assert

	if result {
		return
	}

	t.Error(
		"Arrange: ", []string{lhs, rhs},
		"Act: ", "compareQuery",
		"Assert: ", result,
		"Awaiting: ", !result,
	)
}

func TestSimpleDifferentQueryCoparasionFailure(t *testing.T) {

	// Arrange

	lhs := "load?data=1"
	rhs := "make?data=1"

	// Act

	result := compareQuery(lhs, rhs)

	// Assert

	if !result {
		return
	}

	t.Error(
		"Arrange: ", []string{lhs, rhs},
		"Act: ", "compareQuery",
		"Assert: ", result,
		"Awaiting: ", !result,
	)
}

func TestWrongQuryComparasinFailure(t *testing.T) {

	// Arrange

	lhs := "data=1"
	rhs := "make?data=1"

	// Act

	result := compareQuery(lhs, rhs)

	// Assert

	if !result {
		return
	}

	t.Error(
		"Arrange: ", []string{lhs, rhs},
		"Act: ", "compareQuery",
		"Assert: ", result,
		"Awaiting: ", !result,
	)
}
