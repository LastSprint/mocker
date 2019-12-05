package mock

import (
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
			RequestModel{},
			RequestModel{},
			RequestModel{},
			RequestModel{},
			RequestModel{
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
		IsOnly: newTrue(),
	}

	secondIsOnly := RequestModel{
		IsOnly: newTrue(),
	}

	group := RequestModelGroup{
		models: []RequestModel{
			RequestModel{},
			RequestModel{},
			firstIsOnly,
			RequestModel{},
			RequestModel{},
			secondIsOnly,
		},
	}

	// Act

	result := group.findIsOnlyMock()

	// Assert

	if *result == secondIsOnly {
		t.Fail()
	}

	if *result != firstIsOnly {
		t.Fail()
	}
}

func TestFindIsOnlyMockReturnsNil(t *testing.T) {
	// Arrange

	group := RequestModelGroup{
		models: []RequestModel{
			RequestModel{},
			RequestModel{},
			RequestModel{},
			RequestModel{},
		},
	}

	// Act

	result := group.findIsOnlyMock()

	// Assert

	if result != nil {
		t.Fail()
	}
}
