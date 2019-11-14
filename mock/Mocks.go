package mock

import (
	"encoding/json"
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

	FilePath string `json:"-"`
}

// RequestModelGroup это модель для группы моковых файлов
type RequestModelGroup struct {
	models          []RequestModel
	URL             string
	Method          string
	iteratorIndexes map[string]int
}

// Next iterate on next element in array of RequestModelGroup
func (model *RequestModelGroup) Next(path string) *RequestModel {

	iteratorIndex := model.iteratorIndexes[path]

	found := model.findFirstMatchedIndex(path, iteratorIndex)

	if found == -1 {
		return nil
	}

	model.iteratorIndexes[path] = found + 1

	return &model.models[found]
}

func (model *RequestModelGroup) findFirstMatchedIndex(path string, currentIndex int) int {

	if currentIndex >= len(model.models) {
		currentIndex = 0
	}

	for index := currentIndex; index < len(model.models); index++ {

		if isGroupInSpecificPath(path, model.models[index].FilePath) {
			return index
		}
	}

	if currentIndex == 0 {
		return -1
	}

	return model.findFirstMatchedIndex(path, 0)
}

func FindGroupByURL(groups *[]RequestModelGroup, url string, method string) *RequestModelGroup {

	for index := 0; index < len(*groups); index++ {

		isPathesEqual := CompareURLPath(url, (*groups)[index].URL)
		isMethodsEqual := strings.Compare(method, (*groups)[index].Method) == 0

		if isPathesEqual && isMethodsEqual {
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
			newGroup.iteratorIndexes = map[string]int{}
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

func isGroupInSpecificPath(specificPath, groupURL string) bool {

	if specificPath == "" {
		return true
	}

	specPathSplited := strings.Split(specificPath, "/")
	groupPathSplited := strings.Split(groupURL, "/")

	if len(groupPathSplited) < len(specPathSplited) {
		return false
	}

	for index, item := range specPathSplited {
		if strings.Compare(item, groupPathSplited[index]) != 0 {
			return false
		}
	}

	return true
}
