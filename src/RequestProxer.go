package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
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
			}).DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	// just copy request
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

	// copy headers and return system information

	newRequest.Header = r.Header.Clone()
	newRequest.Header.Del(redirectHostHeaderKey)
	newRequest.Header.Del(redirectURLSchemeHeaderKey)
	newRequest.Header.Del(isNeedProxyHeaderKey)

	newRequest.Body = r.Body
	newRequest.Method = r.Method

	// do request

	rd, _ := ioutil.ReadAll(newRequest.Body)

	newRequest.Body = ioutil.NopCloser(bytes.NewBuffer(rd))

	resp, err := client.Do(&newRequest)

	fmt.Println(string(rd))
	fmt.Println(r.Header)

	if err != nil {
		// Если вернулась ошибка - возвращаем ошибку
		return resp, err
	}

	if resp.StatusCode != 200 {
		// Если http-статус != 200 то возвращаем ответ и возможную ошибку
		return resp, err
	}

	var data []byte

	data, err = ioutil.ReadAll(resp.Body)

	if err != nil {
		// Если считать тело не удалось - возвращаем ответ сервера и ошибку
		return resp, err
	}

	resp.Body = ioutil.NopCloser(bytes.NewBuffer(data))

	var responseJSON interface{}

	err = json.Unmarshal(data, &responseJSON)

	if err != nil {
		// Если после десериализации тела ответа сервера в JSON произошла ошибка - возвращаем ошибку и ответ
		return resp, err
	}

	// if we get the response the we will write it to new JSON file
	go saveNewMock(&newRequest, responseJSON, projectID)

	return resp, nil
}

func saveNewMock(req *http.Request, responseBody interface{}, projectID string) {

	mutex.Lock()

	defer mutex.Unlock()

	// create path from URL and project
	dirPath := getDirPathFromURL(req.URL, projectID)

	// create folders from dirPath
	err := os.MkdirAll(dirPath, os.ModePerm)

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

	fileName, err := getFileName(req, responseBody)

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

	filePath := filepath.Join(dirPath, fileName)

	mockModel := mock.RequestModel{}

	mockUrl := req.URL.Path

	if len(req.URL.RawQuery) != 0 {
		mockUrl += "?" + req.URL.RawQuery
	}

	mockModel.URL = mockUrl
	mockModel.Method = req.Method
	mockModel.StatusCode = 200

	mockModel.Response = responseBody

	data, err := json.Marshal(mockModel)

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

func getDirPathFromURL(model *url.URL, project string) string {
	return filepath.Join(configuration.MocksRootDir, project, model.Host, model.Path)
}

func getUniqName() string {
	val := uuid.New()
	return val.String()
}
