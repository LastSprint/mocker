package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"mocker/mock"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/mitchellh/hashstructure"
	log "github.com/sirupsen/logrus"
)

var mutex sync.Mutex

func startProxing(r *http.Request, host string, scheme string, projectID string) (*http.Response, error) {
	newRequest := http.Request{}

	client := http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			TLSClientConfig: &tls.Config{
				// See comment above.
				// UNSAFE!
				// DON'T USE IN PRODUCTION!
				InsecureSkipVerify: true,
			},
		},
	}

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
	go saveNewMock(&newRequest, resp, responseJSON, projectID)

	return resp, nil
}

func saveNewMock(req *http.Request, resp *http.Response, responseBody map[string]interface{}, projectID string) {

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
	filePath = filepath.Join(projectID, filePath)

	mock := mock.RequestModel{}

	mockUrl := "/" + projectID + req.URL.Path

	if len(req.URL.RawQuery) != 0 {
		mockUrl += "?" + req.URL.RawQuery
	}

	mock.URL = mockUrl
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

	var prettyJSON bytes.Buffer
	_ = json.Indent(&prettyJSON, data, "", "\t")

	err = ioutil.WriteFile(filePath, prettyJSON.Bytes(), os.ModePerm)

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
