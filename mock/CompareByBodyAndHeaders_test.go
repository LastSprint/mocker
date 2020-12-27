package mock

import "testing"

// Checks that `CompareByBody` will return nil if all macks dont have body and headers
func TestCompareByBodyAndHeadersReturnsNilForEmptyInput(t *testing.T) {
	// Arrange

	inputBody := []byte("123")
	inputHeader := map[string]string{"123": "123"}

	mocks := []RequestModel{
		{},
		{},
	}

	group := RequestModelGroup{
		models:          mocks,
		URL:             "/tmp",
		Method:          "GET",
		iteratorIndexes: map[string]int{},
	}

	// Act

	result := group.LookUpByBodyAndHeaders(inputBody, inputHeader)

	// Assert

	if result != nil {
		t.Fatal("Expected nil but got:", result)
	}
}

// Checks that `CompareByBody` will return mock witch same body and header
func TestCompareByBodyAndHeadersReturnsEntireMatchedItemFirstly(t *testing.T) {
	// Arrange

	acceptedUrl := "/tmp/1"
	jsonRaw := `{"key":"value"}`
	jsonMap := map[string]string{"key": "value"}
	inputBody := []byte(jsonRaw)
	inputHeader := map[string]string{"123": "123"}

	mocks := []RequestModel{
		{RequestHeaders: inputHeader},
		{Request: jsonMap},
		{Request: jsonMap, RequestHeaders: inputHeader, URL: acceptedUrl},
	}

	group := RequestModelGroup{
		models:          mocks,
		URL:             "/tmp",
		Method:          "GET",
		iteratorIndexes: map[string]int{},
	}

	// Act

	result := group.LookUpByBodyAndHeaders(inputBody, inputHeader)

	// Assert

	if result == nil {
		t.Fatal("Ижидалось получить мок с url ", acceptedUrl, "Вместо этого получили nil")
		return
	}

	if result.URL != acceptedUrl {
		t.Fatal("Ожидалось что вернется мок с URL", acceptedUrl, ", а вернулось:", result)
	}
}

// checks that `CompareByBody` will return mock with same body if there isn't mock with same body and headers
func TestCompareByBodyAndHeadersReturnsBodyMatchedItemSecondly(t *testing.T) {
	// Arrange

	acceptedUrl := "/tmp/1"
	jsonRaw := `{"key":"value"}`
	jsonMap := map[string]string{"key": "value"}
	inputBody := []byte(jsonRaw)
	inputHeader := map[string]string{"123": "123"}

	mocks := []RequestModel{
		{RequestHeaders: inputHeader},
		{Request: map[string]string{"key": "value1"}, RequestHeaders: inputHeader},
		{Request: jsonMap, URL: acceptedUrl},
	}

	group := RequestModelGroup{
		models:          mocks,
		URL:             "/tmp",
		Method:          "GET",
		iteratorIndexes: map[string]int{},
	}

	// Act

	result := group.LookUpByBodyAndHeaders(inputBody, inputHeader)

	// Assert

	if result == nil {
		t.Fatal("Ижидалось получить мок с url ", acceptedUrl, "Вместо этого получили nil")
		return
	}

	if result.URL != acceptedUrl {
		t.Fatal("Ожидалось что вернется мок с URL", acceptedUrl, ", а вернулось:", result)
	}
}

// checks that `CompareByBody` will return mock which has the same headers if there isn't mock which have the same body
func TestCompareByBodyAndHeadersReturnsHadersMatchedItemThirdly(t *testing.T) {
	// Arrange

	acceptedUrl := "/tmp/1"
	jsonRaw := `{"key":"value"}`
	jsonWrongMap := map[string]string{"key": "value1"}
	inputBody := []byte(jsonRaw)
	inputHeader := map[string]string{"123": "123"}

	mocks := []RequestModel{
		{Request: jsonWrongMap, RequestHeaders: map[string]string{"123": "1234"}},
		{Request: jsonWrongMap},
		{RequestHeaders: inputHeader, URL: acceptedUrl},
	}

	group := RequestModelGroup{
		models:          mocks,
		URL:             "/tmp",
		Method:          "GET",
		iteratorIndexes: map[string]int{},
	}

	// Act

	result := group.LookUpByBodyAndHeaders(inputBody, inputHeader)

	// Assert

	if result == nil {
		t.Fatal("Ижидалось получить мок с url ", acceptedUrl, "Вместо этого получили nil")
		return
	}

	if result.URL != acceptedUrl {
		t.Fatal("Ожидалось что вернется мок с URL", acceptedUrl, ", а вернулось:", result)
	}
}
