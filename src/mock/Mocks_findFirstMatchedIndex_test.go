package mock

import "testing"

import "fmt"

func TestFindFirstMatchedIndexSuccess(t *testing.T) {
	// Arrange

	path := "/test"
	method := "GET"

	group := RequestModelGroup{
		URL:             path,
		Method:          method,
		iteratorIndexes: map[string]int{"path": 0},
		models: []RequestModel{
			RequestModel{
				FilePath: path,
			},
		},
	}

	// Act

	result := group.findFirstMatchedIndex(path, 0)

	// Assert

	if result != 0 {
		t.Fail()
	}
}

func TestFindFirstMatchedIndexSuccessWithEmptyMap(t *testing.T) {
	// Arrange

	path := "/test"
	method := "GET"

	group := RequestModelGroup{
		URL:             path,
		Method:          method,
		iteratorIndexes: map[string]int{},
		models: []RequestModel{
			RequestModel{
				FilePath: path,
			},
		},
	}

	// Act

	result := group.findFirstMatchedIndex(path, 0)

	// Assert

	if result != 0 {
		t.Fail()
	}
}

func TestFindFirstMatchedIndexWithTwoRequests(t *testing.T) {
	// Arrange

	path := "/test"
	method := "GET"

	group := RequestModelGroup{
		URL:             path,
		Method:          method,
		iteratorIndexes: map[string]int{},
		models: []RequestModel{
			RequestModel{
				FilePath: "/path",
			},
			RequestModel{
				FilePath: path,
			},
		},
	}

	// Act

	result := group.findFirstMatchedIndex(path, 0)

	// Assert

	if result != 1 {
		t.Fail()
	}
}

func TestFindFirstMatchedIndexWhenStartIndexMoreThenLen(t *testing.T) {
	// Arrange

	path := "/test"
	method := "GET"

	group := RequestModelGroup{
		URL:             path,
		Method:          method,
		iteratorIndexes: map[string]int{},
		models: []RequestModel{
			RequestModel{
				FilePath: "/path",
			},
			RequestModel{
				FilePath: path,
			},
		},
	}

	// Act

	result := group.findFirstMatchedIndex(path, 3)

	// Assert

	if result != 1 {
		t.Fail()
	}
}

func TestFindFirstMatchedIndexWithStartIndexMoreThenNeededRequest(t *testing.T) {
	// Arrange

	path := "/test"
	method := "GET"

	group := RequestModelGroup{
		URL:             path,
		Method:          method,
		iteratorIndexes: map[string]int{},
		models: []RequestModel{
			RequestModel{
				FilePath: path,
			},
			RequestModel{
				FilePath: "/path",
			},
			RequestModel{
				FilePath: "/tmp" + path,
			},
		},
	}

	// Act

	result := group.findFirstMatchedIndex(path, 1)

	// Assert
	fmt.Println(group.models[2])

	fmt.Print(result)

	if result != 0 {
		t.Error(result)
	}
}

func TestFindFirstMatchedIndexSuccessForWrongPath(t *testing.T) {
	// Arrange

	path := "/test"
	method := "GET"

	group := RequestModelGroup{
		URL:             path,
		Method:          method,
		iteratorIndexes: map[string]int{},
		models: []RequestModel{
			RequestModel{
				FilePath: "/path",
			},
		},
	}

	// Act

	result := group.findFirstMatchedIndex(path, 0)

	// Assert

	if result != -1 {
		t.Fail()
	}
}
