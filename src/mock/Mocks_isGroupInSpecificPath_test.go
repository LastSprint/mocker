package mock

import (
	"testing"
)

// ----- isGroupInSpecificPath Tests -------

func TestIsGroupInSpecificPathTrueForEmptySpecificPath(t *testing.T) {

	// Arrange

	specificPath := ""
	groupURL := "awesome/test/path"

	// Act

	result := isGroupInSpecificPath(specificPath, groupURL)

	// Assert

	if !result {
		t.Fail()
	}
}

func TestIsGroupInSpecificPathFalseForFullyWrongPath(t *testing.T) {

	// Arrange

	specificPath := "bla/bla/bla"
	groupURL := "awesome/test/path"

	// Act

	result := isGroupInSpecificPath(specificPath, groupURL)

	// Assert

	if result {
		t.Fail()
	}
}

func TestIsGroupInSpecificPathFalseForLastWrongComponent(t *testing.T) {

	// Arrange

	specificPath := "awesome/test/wrong_component"
	groupURL := "awesome/test/path"

	// Act

	result := isGroupInSpecificPath(specificPath, groupURL)

	// Assert

	if result {
		t.Fail()
	}
}

func TestIsGroupInSpecificPathFalseForFirstWrongComponent(t *testing.T) {

	// Arrange

	specificPath := "wrong_component/test/path"
	groupURL := "awesome/test/path"

	// Act

	result := isGroupInSpecificPath(specificPath, groupURL)

	// Assert

	if result {
		t.Fail()
	}
}

func TestIsGroupInSpecificPathTrueForSamePathes(t *testing.T) {

	// Arrange

	specificPath := "awesome/test/path"
	groupURL := "awesome/test/path"

	// Act

	result := isGroupInSpecificPath(specificPath, groupURL)

	// Assert

	if !result {
		t.Fail()
	}
}

func TestIsGroupInSpecificPathFalseForLongSpecificPath(t *testing.T) {

	// Arrange

	specificPath := "awesome/test/path/test"
	groupURL := "awesome/test/path"

	// Act

	result := isGroupInSpecificPath(specificPath, groupURL)

	// Assert

	if result {
		t.Fail()
	}
}

func TestIsGroupInSpecificPathSuccessForSubpath(t *testing.T) {

	// Arrange

	specificPath := "awesome/test/path"
	groupURL := "awesome/test/path/test"

	// Act

	result := isGroupInSpecificPath(specificPath, groupURL)

	// Assert

	if !result {
		t.Fail()
	}
}
