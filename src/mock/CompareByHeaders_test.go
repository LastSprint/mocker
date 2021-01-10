package mock

import "testing"

func TestCompareByHeadersReturnsNilForEmptyRequestHeader(t *testing.T) {

	// Arrange

	models := []RequestModel{
		{RequestHeaders: nil},
		{RequestHeaders: map[string]string{}},
	}

	group := &RequestModelGroup{
		models:          models,
		URL:             "/tmp",
		Method:          "GET",
		iteratorIndexes: map[string]int{},
	}

	input := map[string]string{"tmp": "tmp"}

	// Act

	result := group.CompareByHeaders(input)

	// Assert

	if result != nil {
		t.Error("If there isn't`RequestHeaders` in mock then nil should be returned, but got", result)
	}
}

func TestCompareByHeadersReturnsRightRequest(t *testing.T) {

	// Arrange

	right := RequestModel{RequestHeaders: map[string]string{"tmp": "tmp"}, URL: "/tmp/1"}

	models := []RequestModel{
		{RequestHeaders: nil},
		{RequestHeaders: map[string]string{}},
		{RequestHeaders: map[string]string{"NoTmp": "NoTmp"}},
		right,
	}

	group := &RequestModelGroup{
		models:          models,
		URL:             "/tmp",
		Method:          "GET",
		iteratorIndexes: map[string]int{},
	}

	input := map[string]string{"tmp": "tmp"}

	// Act

	result := group.CompareByHeaders(input)

	// Assert

	if result == nil {
		t.Error("Expected a mock")
		return
	}

	if result.URL != right.URL {
		t.Error("Expected: ", right, "But got: ", result)
	}
}

func TestCompareByHeadersReturnsNilWhenThereIsNotRightMock(t *testing.T) {
	// Arrange

	models := []RequestModel{
		{RequestHeaders: nil},
		{RequestHeaders: map[string]string{}},
		{RequestHeaders: map[string]string{"NoTmp": "NoTmp"}},
	}

	group := &RequestModelGroup{
		models:          models,
		URL:             "/tmp",
		Method:          "GET",
		iteratorIndexes: map[string]int{},
	}

	input := map[string]string{"tmp": "tmp"}

	// Act

	result := group.CompareByHeaders(input)

	// Assert

	if result != nil {
		t.Error("Expected nil but got: ", result)
		return
	}
}

func TestCompareByHeadersReturnsRightRequestWithDifferentMocks(t *testing.T) {

	// Arrange

	right := RequestModel{RequestHeaders: map[string]string{"tmp": "tmp"}, URL: "/tmp/1"}

	models := []RequestModel{
		{RequestHeaders: nil},
		{RequestHeaders: map[string]string{}},
		{RequestHeaders: map[string]string{"NoTmp": "NoTmp"}},
		{RequestHeaders: map[string]string{"tmp": "tmp", "NoTmp": "NoTmp"}},
		right,
	}

	group := &RequestModelGroup{
		models:          models,
		URL:             "/tmp",
		Method:          "GET",
		iteratorIndexes: map[string]int{},
	}

	input := map[string]string{"tmp": "tmp"}

	// Act

	result := group.CompareByHeaders(input)

	// Assert

	if result == nil {
		t.Error("Expected a mock")
		return
	}

	if result.URL != right.URL {
		t.Error("Expected: ", right, "But got: ", result)
	}
}
