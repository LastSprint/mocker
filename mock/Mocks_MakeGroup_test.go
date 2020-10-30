package mock

import (
	"reflect"
	"testing"
)

func TestMakeGroupWorkSuccessForSameMocks(t *testing.T) {

	// Arrange

	model := RequestModel{
		URL:      "/path/to/point",
		Method:   "GET",
		Response: nil,
	}

	mocks := []RequestModel{model, model}

	// Act

	groups := MakeGroups(mocks)

	// Assert

	if len(groups) != 1 {
		t.Fail()
	}

	for _, item := range groups[0].models {
		if !reflect.DeepEqual(item, model) {
			t.Fail()
		}
	}
}

func TestMakeGroupWorkSuccessForEqualsMocksWithDifferentResponse(t *testing.T) {
	// Arrange

	baseModel := RequestModel{
		URL:    "/path/to/point",
		Method: "GET",
	}

	model1 := baseModel
	model2 := baseModel

	model1.Response = "model"
	model2.Response = "test"

	mocks := []RequestModel{model1, model2}

	// Act

	groups := MakeGroups(mocks)

	// Assert

	if len(groups) != 1 {
		t.Fail()
	}

	for _, item := range groups[0].models {
		if !reflect.DeepEqual(item, model1) && !reflect.DeepEqual(item, model2) {
			t.Fail()
		}
	}
}

func TestMakeGroupWorkSuccessForEqualsMocksWithDifferentRquest(t *testing.T) {
	// Arrange

	baseModel := RequestModel{
		URL:    "/path/to/point",
		Method: "GET",
	}

	model1 := baseModel
	model2 := baseModel

	model1.Request = "model"
	model2.Request = "test"

	mocks := []RequestModel{model1, model2}

	// Act

	groups := MakeGroups(mocks)

	// Assert

	if len(groups) != 1 {
		t.Fail()
	}

	for _, item := range groups[0].models {
		if !reflect.DeepEqual(item, model1) && !reflect.DeepEqual(item, model2) {
			t.Fail()
		}
	}
}

func TestMakeGroupWorkSuccessForDifferentUrl(t *testing.T) {
	// Arrange

	model1 := RequestModel{
		URL:    "/path/to/point",
		Method: "GET",
	}

	model2 := RequestModel{
		URL:    "/path/to/point2",
		Method: "GET",
	}

	mocks := []RequestModel{model1, model2}

	// Act

	groups := MakeGroups(mocks)

	// Assert

	if len(groups) != 2 {
		t.Fail()
	}

	if !reflect.DeepEqual(groups[0].models[0], model1) {
		t.Fail()
	}

	if !reflect.DeepEqual(groups[1].models[0], model2) {
		t.Fail()
	}
}

func TestMakeGroupWorkSuccessForDifferentMethods(t *testing.T) {
	// Arrange

	model1 := RequestModel{
		URL:    "/path/to/point",
		Method: "Post",
	}

	model2 := RequestModel{
		URL:    "/path/to/point",
		Method: "GET",
	}

	mocks := []RequestModel{model1, model2}

	// Act

	groups := MakeGroups(mocks)

	// Assert

	if len(groups) != 2 {
		t.Fail()
	}

	if !reflect.DeepEqual(groups[0].models[0], model1) {
		t.Fail()
	}

	if !reflect.DeepEqual(groups[1].models[0], model2) {
		t.Fail()
	}
}

func TestMakeGroupWorkSuccessForManyDifferentMocks(t *testing.T) {

	// Arrange

	mocks := []RequestModel{
		RequestModel{
			URL:    "/path/to",
			Method: "PUT",
		},
		RequestModel{
			URL:    "/path/to",
			Method: "GET",
		},
		RequestModel{
			URL:    "/path/to/point",
			Method: "GET",
		},
		RequestModel{
			URL:    "/path/to/point",
			Method: "Post",
		},
	}

	// Act

	groups := MakeGroups(mocks)

	// Assert

	if len(groups) != 4 {
		t.Fail()
	}

	for _, group := range groups {
		if len(group.models) != 1 {
			t.Fail()
		}
	}
}

func TestMakeGroupWorkSuccessForManySameMocks(t *testing.T) {

	// Arrange

	model := RequestModel{
		URL:    "/path/to",
		Method: "GET",
	}

	model2 := RequestModel{
		URL:    "/path/to",
		Method: "POST",
	}

	mocks := []RequestModel{model, model2, model2, model, model, model2, model, model2}

	// Act

	groups := MakeGroups(mocks)

	// Assert

	if len(groups) != 2 {
		t.Fail()
	}

	if len(groups[0].models) != 4 {
		t.Fail()
	}

	if len(groups[1].models) != 4 {
		t.Fail()
	}

	for _, item := range groups[0].models {
		if !reflect.DeepEqual(item, model) {
			t.Fail()
		}
	}

	for _, item := range groups[1].models {
		if !reflect.DeepEqual(item, model2) {
			t.Fail()
		}
	}
}
