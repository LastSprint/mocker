package mock

import (
	"log"
	"strings"
)

func CompareURLPath(lhs string, mockUrl string) bool {

	splitedLhs := strings.Split(lhs, "/")
	splitedRhs := strings.Split(mockUrl, "/")

	log.Println(splitedLhs)
	log.Println(splitedRhs)

	if len(splitedLhs) != len(splitedRhs) {
		return false
	}

	for index := 0; index < len(splitedLhs); index++ {
		lhsItem := splitedLhs[index]
		mockItem := splitedRhs[index]

		if strings.Contains(mockItem, "{") {
			continue
		}

		if strings.Compare(lhsItem, mockItem) != 0 {
			return false
		}
	}

	log.Println(lhs, "EQUAL", mockUrl)

	return true
}
