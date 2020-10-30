package mock

import "testing"

// проверяет, что метод `CompareByHeaders` вернет nil в случае если в моках нет записанных `RequestHeaders`
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
		t.Error("Если в моке не указан `RequestHeaders` то должен вернуться nil, однако вернулось:", result)
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
		t.Error("Должен был вернуться мок, но ничего не вернулось")
		return
	}

	if result.URL != right.URL {
		t.Error("Ожидалось получить: ", right, "Однако было получено:", result)
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
		t.Error("Должен был вернуться nil, но ничего вернулось:", result)
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
		t.Error("Должен был вернуться мок, но ничего не вернулось")
		return
	}

	if result.URL != right.URL {
		t.Error("Ожидалось получить: ", right, "Однако было получено:", result)
	}
}
