package mock

import "testing"

// Проверяет что метод `CompareByBody` вернет nil в случае если в группе все моки без хедеров и тела
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

	result := group.CompareByBodyAndHeaders(inputBody, inputHeader)

	// Assert

	if result != nil {
		t.Fatal("Ожидалось что вернется nil, а вернулось:", result)
	}
}

// Проверяет что `CompareByBody` вернет тот мок, у которого совпадают и заголовки и тело
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

	result := group.CompareByBodyAndHeaders(inputBody, inputHeader)

	// Assert

	if result == nil {
		t.Fatal("Ижидалось получить мок с url ", acceptedUrl, "Вместо этого получили nil")
		return
	}

	if result.URL != acceptedUrl {
		t.Fatal("Ожидалось что вернется мок с URL", acceptedUrl, ", а вернулось:", result)
	}
}

// Проверяет что `CompareByBody` вернет тот мок, у которого совпадает тело, если нет того, у которого совпадает все
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

	result := group.CompareByBodyAndHeaders(inputBody, inputHeader)

	// Assert

	if result == nil {
		t.Fatal("Ижидалось получить мок с url ", acceptedUrl, "Вместо этого получили nil")
		return
	}

	if result.URL != acceptedUrl {
		t.Fatal("Ожидалось что вернется мок с URL", acceptedUrl, ", а вернулось:", result)
	}
}

// Проверяет что `CompareByBody` вернет тот мок, у которого совпадут заголовки, если нет того, у которого совпадает тело
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

	result := group.CompareByBodyAndHeaders(inputBody, inputHeader)

	// Assert

	if result == nil {
		t.Fatal("Ижидалось получить мок с url ", acceptedUrl, "Вместо этого получили nil")
		return
	}

	if result.URL != acceptedUrl {
		t.Fatal("Ожидалось что вернется мок с URL", acceptedUrl, ", а вернулось:", result)
	}
}
