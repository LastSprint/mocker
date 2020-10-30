package mock

import (
	"encoding/json"
	"reflect"
	"strings"
)

// RequestModel это модель мокового файла
type RequestModel struct {

	// IsDisabled состояние мока. Если `true`, то мок исключается из выдачи.
	// Если IsDisabled == nil, то флаг считается опущенным и мок учавствует в выдаче.
	IsDisabled *bool `json:"isDisabled"`
	// IsOnly указывает на то, что мок, для которого этот флаг `true` становится единственным в выдаче.
	// При этом, если isOnly = true, то `isDisabled` не учитывается.
	// В случае если isOnly == nil, то считается, что флаг опущен.
	// При этом итератор не меняет своей позиции.
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

// RequestModelGroup это модель для группы моковых файлов
type RequestModelGroup struct {
	models          []RequestModel
	URL             string
	Method          string
	iteratorIndexes map[string]int
}

// Next итерирует на следующий элемент в RequestModelGroup
// при этом, он передвигает указатель на следующий мок только для той части моков, которая подходит под path
// Например есть есть два мока `/test/dir` и `/tmp/dir` и в `Next` передали строку `/test`
// То вернется мок `/test/dir` и указатель передвинется на следующий мок с `filePath ~ "/test"`.
// При этом, если затем вызывать у группы `Next` с параметром `/tmp` то вернется мок с `/tmp/dir` и тогда его указатель сдвинется.
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

// FindGroupByURL в группе моков `groups` находит группу, которой соответствуют `url` и `method`
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

// MakeGroups группирует моки, используя метод `FindGroupByURLStruct`
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

// CompareByBodyAndHeaders находит нужный мок по телу запроса и хедерам запроса
//
// Приоритеты:
// 1. Ищет мок у которого одновременно совпадает и тело запроса и хедеры
// 2. Если не выполнилось 1, то ищет мок у которого совпадает тело
// 3. Если не выполнилось 2, то ищет мок у которого совпадают хедеры
func (model *RequestModelGroup) CompareByBodyAndHeaders(body []byte, headers map[string]string) *RequestModel {

	// сначала находим тот у которого и тело и хедеры подходят
	for _, item := range model.models {
		if item.CompareByRequest(body) && item.CompareByHeaders(headers) {
			return &item
		}
	}

	// если мы попали сюда, значит такой не нашелся

	// если не нашелся то возвращаем сначала тот который сматчился по телу

	if val := model.CompareByRequest(body); val != nil {
		return val
	}

	// если не нашелся возвращаем тот который сматчился по заголовкам

	return model.CompareByHeaders(headers)
}

// CompareByRequest вызывает `CompareByRequest` для каждого мока из группы и если находит нужный - возвращает его.
// В противном случае вернется nil.
func (model *RequestModelGroup) CompareByRequest(requestData []byte) *RequestModel {
	for index := 0; index < len(model.models); index++ {
		if model.models[index].CompareByRequest(requestData) {
			return &model.models[index]
		}
	}
	return nil
}

// CompareByHeaders используется для того, чтобы сравнить "входящие" header-ы `headers`
// С "ожидаемыми" (прописаными в моке) заголовках `RequestModel.RequestHeaders`
// Метод самостоятельно нормилизует ключи приводя их к нижнему регистру
// Проверка происходит в "мягком" виде. То есть сравнивается только то что пришедшие хедеры содержаться в моке и значения равны
func (model *RequestModelGroup) CompareByHeaders(headers map[string]string) *RequestModel {
	for _, item := range model.models {
		if item.CompareByHeaders(headers) {
			return &item
		}
	}
	return nil
}

func (model *RequestModel) CompareByHeaders(headers map[string]string) bool {
	// если параметр RequestHeaders был не указан, то мы даже не будем сравнивать
	if len(model.RequestHeaders) == 0 {
		return false
	}

	normalizedHeaders := map[string]string{}

	for key, value := range headers {
		normalizedHeaders[strings.ToLower(key)] = value
	}

	for key, value := range model.RequestHeaders {
		val, ok := normalizedHeaders[strings.ToLower(key)]

		// если такой ключ есть в моке
		// и значение по этому ключу равно значению в запросе
		// значит все ок и пропускаем
		if ok && val == value {
			continue
		}
		return false
	}

	// если нашли мок у которого все ключи совпали, значит это наш парень
	// возврщаем его и заканчиваем поиск
	// в противном случае берем следующий мок и по-новой
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

// findIsOnlyMock находит мок, у которого флаг `IsOnly == true`
// В слуачае, если такого мока нет, то вернется nil.
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
