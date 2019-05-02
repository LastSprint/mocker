package mock

import (
	"log"
	"strings"
)

type RequestModel struct {
	URL        string      `json:"url"`
	Response   interface{} `json:"response"`
	Method     string      `json:"method"`
	StatusCode int         `json:"statusCode"`
}

type RequestModelGroup struct {
	models        []RequestModel
	URL           string
	Method        string
	iteratorIndex int
}

// Next iterate on next element in array of RequestModelGroup
func (model *RequestModelGroup) Next() *RequestModel {
	log.Println("NEXT")
	if model.iteratorIndex == len(model.models) {
		model.iteratorIndex = 0
		return &model.models[0]
	}

	result := &model.models[model.iteratorIndex]

	model.iteratorIndex++

	return result
}

func FindGroupByURL(groups *[]RequestModelGroup, url string, method string) *RequestModelGroup {

	for index := 0; index < len(*groups); index++ {

		if CompareURLPath(url, (*groups)[index].URL) && strings.Compare(method, (*groups)[index].Method) == 0 {
			return &(*groups)[index]
		}
	}
	return nil
}

func FindGroupByURLStruct(groups []RequestModelGroup, url string, method string) *RequestModelGroup {

	for index := 0; index < len(groups); index++ {

		if CompareURLPath(url, groups[index].URL) && strings.Compare(method, groups[index].Method) == 0 {
			return &groups[index]
		}
	}
	return nil
}

func MakeGroups(allMocks []RequestModel) []RequestModelGroup {
	var result []RequestModelGroup

	for _, item := range allMocks {
		group := FindGroupByURLStruct(result, item.URL, item.Method)

		if group == nil {
			newGroup := RequestModelGroup{}
			newGroup.Method = item.Method
			newGroup.URL = item.URL
			newGroup.models = []RequestModel{item}
			result = append(result, newGroup)
		} else {
			group.models = append(group.models, item)
		}
	}

	return result
}
