package mock

import (
	"encoding/json"
	"reflect"
	"strings"
)

// RequestModel this is a mock model
type RequestModel struct {

	// IsDisabled the mock state. if is set to `true`, the mock will be skipped in server response.
	// nil value is the same as `false`
	IsDisabled *bool `json:"isDisabled"`

	// IsOnly if set to `true` the mock will be the only one in responses.
	// Other mocks from group won't be processing
	//
	// If `isOnly` = true then `isDisabled` won't considered
	// nil value is the same as `false`
	IsOnly *bool `json:"isOnly"`

	URL        string      `json:"url"`
	Response   interface{} `json:"response"`
	Method     string      `json:"method"`
	StatusCode int         `json:"statusCode"`
	Request    interface{} `json:"request"`

	FilePath string `json:"-"`

	Delay int `json:"responseDelay"`

	ResponseHeaders map[string]string `json:"responseHeaders"`
	RequestHeaders  map[string]string `json:"requestHeaders"`

	IsExcludedFromIteration *bool `json:"isExcludedFromIteration"`
}

// RequestModelGroup this is a model for group of mocks
type RequestModelGroup struct {
	models          []RequestModel
	URL             string
	Method          string
	iteratorIndexes map[string]int
}

// Next iterates at the next element следующий элемент in RequestModelGroup
// and the pointer will be moved only for the part of mocks which is confirmed to `path`
// For example
// There are two mocks: `/test/dir` and `/tmp/dir`
// If `path` is `/test` then will be returned `test/dir` and pointer will be moved to next mock with path `/test/*`
func (model *RequestModelGroup) Next(path string) *RequestModel {

	if mock := model.findIsOnlyMock(); mock != nil {
		return mock
	}

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

		if isGroupInSpecificPath(path, model.models[index].FilePath) && !model.models[index].isDisabled() && !model.models[index].isExcludedFromIteration() {
			return index
		}
	}

	if currentIndex == 0 {
		return -1
	}

	return model.findFirstMatchedIndex(path, 0)
}

// FindGroupByURL tries to find specific group by url and method
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

// MakeGroups create groups from plain array of mocks
func MakeGroups(allMocks []RequestModel) []RequestModelGroup {
	var result []RequestModelGroup

	for _, item := range allMocks {
		group := FindGroupByURL(&result, item.URL, item.Method)

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

// CompareByRequest rules:
// - If `RequestModel.Request` == nil will return false
// - If while marshaling `RequestModel.Request` an error is thrown will return false
// - If bodies are not equal (byte-to-byte) will return false
// ------
// - Parameters:
//	- requestData: binary request body (raw)
func (model *RequestModel) CompareByRequest(requestData []byte) bool {

	if model.Request == nil {
		return false
	}

	modeRequestData, err := json.Marshal(model.Request)

	if err != nil {
		return false
	}

	var bytes interface{}

	err = json.Unmarshal(requestData, &bytes)

	if err != nil {
		return false
	}

	resultReuqestBytes, err := json.Marshal(bytes)

	if err != nil {
		return false
	}

	return reflect.DeepEqual(modeRequestData, resultReuqestBytes)
}

// LookUpByBodyAndHeaders looks up for specific mock by body and headers
//
// Priorities:
// 1. Looks up for mock which has equal Request (to body) and Headers
// 2. If 1 is wrong then will look up for mock with same body only
// 3. If 2 is wrong then will look up for mock with same headers
func (model *RequestModelGroup) LookUpByBodyAndHeaders(body []byte, headers map[string]string) *RequestModel {

	// 1
	for _, item := range model.models {
		if item.CompareByRequest(body) && item.CompareByHeaders(headers) {
			return &item
		}
	}

	// there is no mock with same headers and body at the same time

	// 2
	if val := model.CompareByRequest(body); val != nil {
		return val
	}

	// 3
	return model.CompareByHeaders(headers)
}

// CompareByRequest calls `CompareByRequest` for each mock from the group
// if finds correct mock then will return it. Or will return nil
func (model *RequestModelGroup) CompareByRequest(requestData []byte) *RequestModel {
	for index := 0; index < len(model.models); index++ {
		if model.models[index].CompareByRequest(requestData) {
			return &model.models[index]
		}
	}
	return nil
}

func (model *RequestModelGroup) CompareByHeaders(headers map[string]string) *RequestModel {
	for _, item := range model.models {
		if item.CompareByHeaders(headers) {
			return &item
		}
	}
	return nil
}

func (model *RequestModel) CompareByHeaders(headers map[string]string) bool {

	if len(model.RequestHeaders) == 0 {
		return false
	}

	normalizedHeaders := map[string]string{}

	for key, value := range headers {
		normalizedHeaders[strings.ToLower(key)] = value
	}

	for key, value := range model.RequestHeaders {
		val, ok := normalizedHeaders[strings.ToLower(key)]

		if ok && val == value {
			continue
		}
		return false
	}

	return true
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

func (model *RequestModelGroup) findIsOnlyMock() *RequestModel {
	for _, mock := range model.models {
		if mock.IsOnly != nil && *mock.IsOnly == true {
			return &mock
		}
	}
	return nil
}

func (model RequestModel) isDisabled() bool {
	if model.IsDisabled == nil {
		return false
	}

	if *model.IsDisabled == false {
		return false
	}

	return true
}

func (model RequestModel) isExcludedFromIteration() bool {
	if model.IsExcludedFromIteration == nil {
		return false
	}

	if *model.IsExcludedFromIteration == true {
		return true
	}

	return false
}
