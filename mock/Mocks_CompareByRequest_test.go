package mock

import (
	"encoding/json"
	"testing"
)

func TestCompareByRequestCompliteSuccessForSameJsons(t *testing.T) {

	// Arrange

	jsonRaw := map[string]interface{}{
		"name": "Test",
	}

	group := RequestModelGroup{
		models: []RequestModel{
			RequestModel{
				Request: jsonRaw,
			},
		},
	}

	// Act

	jsonData, _ := json.Marshal(jsonRaw)

	model := group.CompareByRequest(jsonData)

	// Assert

	if model == nil {
		t.Fail()
	}
}

func TestCompareByRequestFailedForDifferentJsons(t *testing.T) {
	// Arrange

	jsonRequest := map[string]interface{}{
		"name": "Test",
	}

	jsonMock := map[string]interface{}{
		"step": 1,
	}

	group := RequestModelGroup{
		models: []RequestModel{
			RequestModel{
				Request: jsonRequest,
			},
		},
	}

	// Act

	jsonData, _ := json.Marshal(jsonMock)

	model := group.CompareByRequest(jsonData)

	// Assert

	if model != nil {
		t.Fail()
	}
}

func TestCompareByRequestFailedForEmptyModelRequest(t *testing.T) {
	// Arrange

	jsonMock := map[string]interface{}{
		"step": 1,
	}

	group := RequestModelGroup{
		models: []RequestModel{
			RequestModel{
				Request: nil,
			},
		},
	}

	// Act

	jsonData, _ := json.Marshal(jsonMock)

	model := group.CompareByRequest(jsonData)

	// Assert

	if model != nil {
		t.Fail()
	}
}

func TestCompareByRequestFailedForNotJsonBytesArray(t *testing.T) {
	// Arrange

	json := map[string]interface{}{
		"step": 1,
	}

	group := RequestModelGroup{
		models: []RequestModel{
			RequestModel{
				Request: json,
			},
		},
	}

	// Act

	data := []byte{1, 1, 1, 1, 1, 1}

	model := group.CompareByRequest(data)

	// Assert

	if model != nil {
		t.Fail()
	}
}

func TestCompareByRequestFailedForNotJsonRequest(t *testing.T) {
	// Arrange

	requestValue := "testString"

	jsonValue := map[string]interface{}{
		"step": 1,
	}

	group := RequestModelGroup{
		models: []RequestModel{
			RequestModel{
				Request: requestValue,
			},
		},
	}

	// Act

	jsonData, _ := json.Marshal(jsonValue)

	model := group.CompareByRequest(jsonData)

	// Assert

	if model != nil {
		t.Fail()
	}
}

// ЗАметили такой баг на одном из реальных кейсов
func TestCompareByRequestFailedWithPhoneAndLongToken(t *testing.T) {
	// Arrange

	requestValue := map[string]interface{}{
		"phone": "71111111112",
		"fcmToken": "31f31cf5-2ec0-459c-8c17-6c4b64c69161-1f840a43-e36b-4163-b69b-ce3db09d5ca2",
	}

	group := RequestModelGroup{
		models: []RequestModel{
			RequestModel{
				Request: requestValue,
			},
		},
	}

	// Act

	httpRequestString := `{"phone":"71111111112","fcmToken":"31f31cf5-2ec0-459c-8c17-6c4b64c69161-1f840a43-e36b-4163-b69b-ce3db09d5ca2"}`

	model := group.CompareByRequest([]byte(httpRequestString))

	// Assert

	if model == nil {
		t.Fail()
	}
}

func TestCompareByRequestFailedWithPhoneAndLongTokenWithDifferentOrder(t *testing.T) {
	// Arrange

	requestValue := map[string]interface{}{
		"phone": "71111111112",
		"fcmToken": "31f31cf5-2ec0-459c-8c17-6c4b64c69161-1f840a43-e36b-4163-b69b-ce3db09d5ca2",
	}

	group := RequestModelGroup{
		models: []RequestModel{
			RequestModel{
				Request: requestValue,
			},
		},
	}

	// Act

	httpRequestString := `{"fcmToken":"31f31cf5-2ec0-459c-8c17-6c4b64c69161-1f840a43-e36b-4163-b69b-ce3db09d5ca2","phone":"71111111112"}`

	model := group.CompareByRequest([]byte(httpRequestString))

	// Assert

	if model == nil {
		t.Fail()
	}
}
