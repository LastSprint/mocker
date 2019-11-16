package mock

import (
	"sort"
	"strings"
)

type parameterSlice []string

func (a parameterSlice) Len() int           { return len(a) }
func (a parameterSlice) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a parameterSlice) Less(i, j int) bool { return a[i] < a[j] }

// CompareURLPath сравнивает для URL по специальным правилам.
// При этом доспукается, что lhs URL может быть шаблоном. Тогда mock будет сравниваться с ним как с шаблоном.
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

		// means that param names equal
		if strings.Contains(mockParam[1], "{") {
			continue
		} else if strings.Compare(lhsParam[1], mockParam[1]) == 0 {
			continue
		}
		return false
	}
	return true
}
