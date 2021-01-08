package mock

import (
	"testing"
)

func TestFindGroupByURLWorksSuccess(t *testing.T) {
	// Arrange

	url := "/test"
	method := "GET"

	groups := []RequestModelGroup{
		{
			URL:    url,
			Method: method,
		},
	}

	// Act

	group := FindGroupByURL(&groups, url, method)

	// Assert

	if group == nil {
		t.Fail()
	}
}

func TestFindGroupByURLWorksSuccessWithDifferentMethod(t *testing.T) {
	// Arrange

	url := "/test"
	method := "GET"

	groups := []RequestModelGroup{
		{
			URL:    url,
			Method: "POST",
		},
	}

	// Act

	group := FindGroupByURL(&groups, url, method)

	// Assert

	if group != nil {
		t.Fail()
	}
}

func TestFindGroupByURLWorksSuccessWithDifferentURL(t *testing.T) {
	// Arrange

	url := "/test"
	method := "GET"

	groups := []RequestModelGroup{
		{
			URL:    "/test/path",
			Method: method,
		},
	}

	// Act

	group := FindGroupByURL(&groups, url, method)

	// Assert

	if group != nil {
		t.Fail()
	}
}

func TestFindGroupByURLWorksSuccessWithDifferentURLAndMethod(t *testing.T) {
	// Arrange

	url := "/test"
	method := "GET"

	groups := []RequestModelGroup{
		{
			URL:    "/test/path",
			Method: "POST",
		},
	}

	// Act

	group := FindGroupByURL(&groups, url, method)

	// Assert

	if group != nil {
		t.Fail()
	}
}
