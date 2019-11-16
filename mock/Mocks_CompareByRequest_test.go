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
