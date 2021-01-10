package mock

import "testing"

func TestSameQueryCompareSuccess(t *testing.T) {
	// Arrange

	lhs := "foo=1&bar=2&baz=3"
	rhs := "foo=1&bar=2&baz=3"

	// Act

	result := compareURLQueryPart(lhs, rhs)

	// Assert

	if result {
		return
	}

	t.Error(
		"Arrange: ", []string{lhs, rhs},
		"Act: ", "compareURLQueryPart",
		"Assert: ", result,
		"Awaiting: ", !result,
	)
}

func TestDifferentQueryCompareFailure(t *testing.T) {
	// Arrange

	lhs := "foo=1&bar=2&baz=3"
	rhs := "foo=1&bar=2&baz=4"

	// Act

	result := compareURLQueryPart(lhs, rhs)

	// Assert

	if !result {
		return
	}

	t.Error(
		"Arrange: ", []string{lhs, rhs},
		"Act: ", "compareURLQueryPart",
		"Assert: ", result,
		"Awaiting: ", !result,
	)
}

func TestSameQueryWithParamsCompareSuccess(t *testing.T) {
	// Arrange

	lhs := "foo=1"
	rhs := "foo={foo_value}"

	// Act

	result := compareURLQueryPart(lhs, rhs)

	// Assert

	if result {
		return
	}

	t.Error(
		"Arrange: ", []string{lhs, rhs},
		"Act: ", "compareURLQueryPart",
		"Assert: ", result,
		"Awaiting: ", !result,
	)
}

func TestSameQueryWithMiddleParamCompareSuccess(t *testing.T) {
	// Arrange

	lhs := "foo=1&bar=2&baz=3"
	rhs := "foo=1&bar={bar_value}&baz=3"

	// Act

	result := compareURLQueryPart(lhs, rhs)

	// Assert

	if result {
		return
	}

	t.Error(
		"Arrange: ", []string{lhs, rhs},
		"Act: ", "compareURLQueryPart",
		"Assert: ", result,
		"Awaiting: ", !result,
	)
}

func TestDifferentQueryWithMiddleParamCompareFailure(t *testing.T) {
	// Arrange

	lhs := "foo=1&bar=2&baz=4"
	rhs := "foo=1&bar={bar_value}&baz=3"

	// Act

	result := compareURLQueryPart(lhs, rhs)

	// Assert

	if !result {
		return
	}

	t.Error(
		"Arrange: ", []string{lhs, rhs},
		"Act: ", "compareURLQueryPart",
		"Assert: ", result,
		"Awaiting: ", !result,
	)
}

func TestSameQuesryWithDifferenetOrderCompareSuccess(t *testing.T) {
	// Arrange

	lhs := "foo=1&bar=2&baz=4"
	rhs := "foo=1&bar=1&baz=3"

	// Act

	result := compareURLQueryPart(lhs, rhs)

	// Assert

	if !result {
		return
	}

	t.Error(
		"Arrange: ", []string{lhs, rhs},
		"Act: ", "compareURLQueryPart",
		"Assert: ", result,
		"Awaiting: ", !result,
	)
}

func TestDifferentNumberOfParamsComparasionFailure(t *testing.T) {
	// Arrange

	lhs := "foo=1&bar=2"
	rhs := "foo=1&bar=1&baz=3"

	// Act

	result := compareURLQueryPart(lhs, rhs)

	// Assert

	if !result {
		return
	}

	t.Error(
		"Arrange: ", []string{lhs, rhs},
		"Act: ", "compareURLQueryPart",
		"Assert: ", result,
		"Awaiting: ", !result,
	)
}
