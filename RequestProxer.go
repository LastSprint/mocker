package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"mocker/mock"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/google/uuid"
	"github.com/mitchellh/hashstructure"
	log "github.com/sirupsen/logrus"
)

var mutex sync.Mutex

func startProxing(r *http.Request, host string, scheme string) (*http.Response, error) {
	newRequest := http.Request{}

	client := http.Client{}

	// Копируем запрос
	newRequest.URL = &url.URL{
		Scheme:     scheme,
		Opaque:     r.URL.Opaque,
		User:       r.URL.User,
		Host:       host,
		Path:       r.URL.Path,
		RawPath:    r.URL.RawPath,
		ForceQuery: r.URL.ForceQuery,
		RawQuery:   r.URL.RawQuery,
		Fragment:   r.URL.Fragment,
	}

	// Копируем хедеры и удаляем служебную информацию

	newRequest.Header = r.Header.Clone()
	newRequest.Header.Del(redirectHostHeaderKey)
	newRequest.Header.Del(redirectURLSchemeHeaderKey)
	newRequest.Header.Del(isNeedProxyHeaderKey)

	newRequest.Body = r.Body
	newRequest.Method = r.Method

	// Выполняем запрос

	resp, err := client.Do(&newRequest)


	if err != nil {
		// Если вернулась ошибка - возвращаем ошибку
		return resp, err
	}

	if resp.StatusCode != 200 {
		// Если http-статус != 200 то возвращаем ответ и возможную ошибку
		return resp, err
	}

	// Ответ нам подходит - считываем тело

	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		// Если считать тело не удалось - возвращаем ответ сервера и ошибку
		return resp, err
	}

	// После того как мы считали тело, то стрим закончился.
	// Нам нужно создать новый стрим с указателем в начале (для того, чтобы можно было это тело записать в ответ клиенту далее)
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(data))

	responseJSON := make(map[string]interface{})

	err = json.Unmarshal(data, &responseJSON)

	if err != nil {
		// Если после десериализации тела ответа сервера в JSON произошла ошибка - возвращаем ошибку и ответ
		return resp, err
	}

	// Если JSON получен, то асинхронно запускаем запись в файл
	go saveNewMock(&newRequest, resp, responseJSON)

	return resp, nil
}

func saveNewMock(req *http.Request, resp *http.Response, responseBody map[string]interface{}) {

	// Кажется, что тут может быть очень неприятная гонка (:

	mutex.Lock()

	defer mutex.Unlock()

	dirPath := getDirPathFromURL(req.URL)

	// Получив путь создаем его - создаст все вложенные папки.
	err := os.MkdirAll(dirPath, os.ModePerm)

	if err != nil {
		// Если при создании папок произошла ошибка - ничего не делаем
		fields := log.Fields{
			"success": false,
			"err":     err,
			"host":    req.URL.Host,
			"url":     req.URL.String(),
		}

		logAnalytics(fields, EventKeyProxyFileSave)
		return
	}

	fileName, err := getFileName(req, responseBody)

	if err != nil {
		// Если при получении имени файла произошла ошибка - ничего не делаем

		fields := log.Fields{
			"success": false,
			"err":     err,
			"host":    req.URL.Host,
			"url":     req.URL.String(),
		}

		logAnalytics(fields, EventKeyProxyFileSave)
		return
	}

	// Получаем итоговый путь до файла
	filePath := filepath.Join(dirPath, fileName)

	mock := mock.RequestModel{}

	mock.URL = req.URL.String()
	mock.Method = req.Method
	mock.StatusCode = 200

	mock.Response = responseBody

	data, err := json.Marshal(mock)

	if err != nil {
		// Если произошла ошибка пр исериализации мока в JSON - ничего не делаем

		fields := log.Fields{
			"success": false,
			"err":     err,
			"host":    req.URL.Host,
			"url":     req.URL.String(),
		}

		logAnalytics(fields, EventKeyProxyFileSave)
		return
	}

	// Записываем файл

	err = ioutil.WriteFile(filePath, data, os.ModePerm)

	if err != nil {
		fields := log.Fields{
			"success": false,
			"err":     err,
			"host":    req.URL.Host,
			"url":     req.URL.String(),
		}

		logAnalytics(fields, EventKeyProxyFileSave)
		return
	}

	fields := log.Fields{
		"success": true,
		"host":    req.URL.Host,
		"url":     req.URL.String(),
	}

	logAnalytics(fields, EventKeyProxyFileSave)
}

func getFileName(request *http.Request, body interface{}) (string, error) {

	hash, err := hashstructure.Hash(body, nil)

	if err != nil {
		return "", err
	}

	return request.Method + "_" + strconv.FormatUint(hash, 16) + ".json", nil
}

func getDirPathFromURL(model *url.URL) string {
	return filepath.Join(configuration.MocksRootDir, model.Host, model.Path)
}

func getUniqName() string {
	val := uuid.New()
	return val.String()
}
