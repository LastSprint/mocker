package mock

import (
	"sort"
	"strings"
)

// CompareURLPath compares urls via specific rules
// and lhs my be an URL template.
// this can compare:
// https://host.dom/path/to/{userId}?param1=2&param2={param2}
// with
// https://host.dom/path/to/1?param1=2&param2=2
// for more details look at tests
func CompareURLPath(lhs, mock string) bool {

	splitedLHS := strings.Split(lhs, "/")
	splitedRHS := strings.Split(mock, "/")

	if len(splitedLHS) != len(splitedRHS) {
		return false
	}

	for index := 0; index < len(splitedLHS); index++ {
		lhsItem := splitedLHS[index]
		mockItem := splitedRHS[index]

		if strings.Contains(lhsItem, "?") {
			return compareQuery(lhsItem, mockItem)
		}

		if strings.Contains(mockItem, "{") {
			continue
		}

		if strings.Compare(lhsItem, mockItem) != 0 {
			return false
		}
	}
	return true
}

func compareQuery(lhs, mock string) bool {
	lhsQuery := strings.Split(lhs, "?")
	mockQuery := strings.Split(mock, "?")

	if len(mockQuery) != 2 || len(lhsQuery) != 2 {
		return false
	}

	if strings.Compare(lhsQuery[0], mockQuery[0]) != 0 {
		return false
	}

	return compareURLQueryPart(lhsQuery[1], mockQuery[1])
}

func compareURLQueryPart(lhs string, rhs string) bool {

	// pattern = paramName={paramValue}&...

	splitedLHS := strings.Split(lhs, "&")
	splitedRHS := strings.Split(rhs, "&")

	sort.Strings(splitedLHS)
	sort.Strings(splitedRHS)

	if len(splitedLHS) != len(splitedRHS) {
		return false
	}

	for index := 0; index < len(splitedLHS); index++ {
		lhsItem := splitedLHS[index]
		mockItem := splitedRHS[index]

		lhsParam := strings.Split(lhsItem, "=")
		mockParam := strings.Split(mockItem, "=")

		if strings.Compare(lhsParam[0], mockParam[0]) != 0 {
			return false
		}

		if len(mockParam) < 2 {
			return false
		}

		// means that param names equal
		if strings.Contains(mockParam[1], "{") {
			continue
		}

		if len(lhsParam) < 2 {
			return false
		}

		if strings.Compare(lhsParam[1], mockParam[1]) == 0 {
			continue
		}
		return false
	}
	return true
}
