package mock

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

// RequestModel это модель мокового файла
type RequestModel struct {
	URL        string      `json:"url"`
	Response   interface{} `json:"response"`
	Method     string      `json:"method"`
	StatusCode int         `json:"statusCode"`
	Request    interface{} `json:"request"`
}

// RequestModelGroup это модель для группы моковых файлов
type RequestModelGroup struct {
	models        []RequestModel
	URL           string
	Method        string
	iteratorIndex int
}

// Next iterate on next element in array of RequestModelGroup
func (model *RequestModelGroup) Next() *RequestModel {

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

// CompareByRequest работает следующим образом:
// - Если `RequestModel.Request` == nil -> false
// - Если при маршалинге `RequestModel.Request` произошла ошибка -> false
// - Если байтовое представление данных не одинаково -> false
// ------
// - Parameters:
//	- requestData: "сырое" бинарное представление тела запроса.
func (model *RequestModel) CompareByRequest(requestData []byte) bool {

	if model.Request == nil {
		return false
	}

	modeRequestData, err := json.Marshal(model.Request)

	var bytes interface{}

	json.Unmarshal(requestData, &bytes)
	fmt.Println(model.Request)
	fmt.Println(bytes)
	if err != nil {
		return false
	}

	return reflect.DeepEqual(modeRequestData, requestData)
}

func (group *RequestModelGroup) CompareByRequest(requestData []byte) *RequestModel {
	for index := 0; index < len(group.models); index++ {

		if group.models[index].CompareByRequest(requestData) {
			return &group.models[index]
		}
	}
	return nil
}
