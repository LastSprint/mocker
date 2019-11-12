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
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

var mutex sync.Mutex

func startProxing(r *http.Request, host string, scheme string) (*http.Response, error) {
	client := http.Client{}

	newRequest := http.Request{}
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
	newRequest.Header = r.Header.Clone()
	newRequest.Header.Del(redirectHostHeaderKey)
	newRequest.Header.Del(redirectURLSchemeHeaderKey)
	newRequest.Header.Del(isNeedProxyHeaderKey)

	newRequest.Body = r.Body
	newRequest.Method = r.Method

	resp, err := client.Do(&newRequest)

	if err != nil {
		return resp, err
	}

	if resp.StatusCode != 200 {
		return resp, err
	}

	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return resp, err
	}

	resp.Body = ioutil.NopCloser(bytes.NewBuffer(data))

	responseJSON := make(map[string]interface{})

	err = json.Unmarshal(data, &responseJSON)

	if err != nil {
		return resp, err
	}

	go saveNewMock(&newRequest, resp, responseJSON)

	defer resp.Body.Close()

	return resp, nil
}

func saveNewMock(req *http.Request, resp *http.Response, responseBody map[string]interface{}) {

	mutex.Lock()

	defer mutex.Unlock()

	dirPath := getDirPathFromURL(req.URL)

	err := os.MkdirAll(dirPath, os.ModePerm)

	if err != nil {
		log.WithFields(log.Fields{
			"key":   "analytics",
			"event": "proxy_file_save",
			"payload": logrus.Fields{
				"success": false,
				"err":     err,
				"host":    req.URL.Host,
				"url":     req.URL.String(),
				"resp":    resp,
			},
		}).Info("ANALYTICS")
		return
	}

	fileName, err := getFileName(req, responseBody)

	if err != nil {
		log.WithFields(log.Fields{
			"key":   "analytics",
			"event": "proxy_file_save",
			"payload": logrus.Fields{
				"success": false,
				"err":     err,
				"host":    req.URL.Host,
				"url":     req.URL.String(),
				"resp":    resp,
			},
		}).Info("ANALYTICS")
		return
	}

	filePath := filepath.Join(dirPath, fileName)

	mock := mock.RequestModel{}

	mock.URL = req.URL.String()
	mock.Method = req.Method
	mock.StatusCode = 200

	mock.Response = responseBody

	data, err := json.Marshal(mock)

	if err != nil {
		log.WithFields(log.Fields{
			"key":   "analytics",
			"event": "proxy_file_save",
			"payload": logrus.Fields{
				"success": false,
				"err":     err,
				"host":    req.URL.Host,
				"url":     req.URL.String(),
				"resp":    resp,
			},
		}).Info("ANALYTICS")
		return
	}

	err = ioutil.WriteFile(filePath, data, os.ModePerm)

	if err != nil {
		log.WithFields(log.Fields{
			"key":   "analytics",
			"event": "proxy_file_save",
			"payload": logrus.Fields{
				"success": false,
				"err":     err,
				"host":    req.URL.Host,
				"url":     req.URL.String(),
				"resp":    resp,
			},
		}).Info("ANALYTICS")
	} else {
		log.WithFields(log.Fields{
			"key":   "analytics",
			"event": "proxy_file_save",
			"payload": logrus.Fields{
				"success": true,
				"err":     err,
				"host":    req.URL.Host,
				"url":     req.URL.String(),
				"resp":    resp,
			},
		}).Info("ANALYTICS")
	}
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
