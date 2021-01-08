package mock

import (
	"reflect"
	"testing"
)

func newTrue() *bool {
	b := true
	return &b
}

func TestFindIsOnlyMockReturnsRightModel(t *testing.T) {
	// Arrange

	group := RequestModelGroup{
		models: []RequestModel{
			{},
			{},
			{},
			{},
			{
				IsOnly: newTrue(),
			},
		},
	}

	// Act

	result := group.findIsOnlyMock()

	// Assert

	if result == nil {
		t.Fail()
	}
}

func TestFindIsOnlyMockReturnsFirstIsOnly(t *testing.T) {

	// Arrange

	firstIsOnly := RequestModel{
		URL:    "sdf",
		IsOnly: newTrue(),
	}

	secondIsOnly := RequestModel{
		IsOnly: newTrue(),
	}

	group := RequestModelGroup{
		models: []RequestModel{
			{},
			{},
			firstIsOnly,
			{},
			{},
			secondIsOnly,
		},
	}

	// Act

	result := group.findIsOnlyMock()

	// Assert

	if reflect.DeepEqual(*result, secondIsOnly) {
		t.Fail()
	}

	if !reflect.DeepEqual(*result, firstIsOnly) {
		t.Fail()
	}
}

func TestFindIsOnlyMockReturnsNil(t *testing.T) {
	// Arrange

	group := RequestModelGroup{
		models: []RequestModel{
			{}, {}, {}, {},
		},
	}

	// Act

	result := group.findIsOnlyMock()

	// Assert

	if result != nil {
		t.Fail()
	}
}
