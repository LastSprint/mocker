package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mocker/config"
	"mocker/mock"
	"net/http"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

var models []mock.RequestModelGroup
var configuration config.Config

const update = "/update_models"

const isNeedProxyHeaderKey = "X-Mocker-Redirect-Is-On"
const redirectHostHeaderKey = "X-Mocker-Redirect-Host"
const redirectURLSchemeHeaderKey = "X-Mocker-Redirect-Scheme"

func main() {
	conf, err := config.LoadConfig(os.Args[1])

	if err != nil {
		log.Fatal(err)
	}

	configuration = conf

	configureLog(&conf)

	updateModels()

	http.HandleFunc("/", handler)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", conf.Port), nil))
}

func handler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	// Достаем служебные параметры, чтобы понять, что нужно выполнить проксирование

	host := r.Header.Get(redirectHostHeaderKey)
	scheme := r.Header.Get(redirectURLSchemeHeaderKey)
	isNeedProxy := r.Header.Get(isNeedProxyHeaderKey)

	if isNeedProxy == "true" && scheme != "" && host != "" {

		data, err := proxyRequest(r, host, scheme)

		if err == nil {
			// Если метод проксирования не вернул ошибки, то просто записывает ответ в response и заканчиваем обработку
			w.WriteHeader(http.StatusOK)
			w.Write(data)
			return
		}
	}

	var fields = log.Fields{}
	fields["Request URL"] = r.URL
	fields["Request Method"] = r.Method

	if strings.Compare(r.URL.String(), update) == 0 {
		err := updateModels()

		log.WithFields(log.Fields{
			"key":   "analytics",
			"event": "update_models",
		}).Info("ANALYTICS")

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"message": "Cant update models"})
			log.WithFields(fields).Warn("Can't update models!")
		}
		return
	}

	item := mock.FindGroupByURL(&models, r.URL.String(), r.Method)

	if item == nil {

		log.WithFields(log.Fields{
			"key":   "analytics",
			"event": "get_mock",
			"payload": logrus.Fields{
				"success": false,
				"err":     "Cant Find Group",
				"url":     r.URL.String(),
			},
		}).Info("ANALYTICS")

		log.WithFields(fields).Warn("Not found any group")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "not found mock for url" + r.URL.String()})
		return
	}

	var next *mock.RequestModel

	body, err := ioutil.ReadAll(r.Body)

	if err == nil {
		next = item.CompareByRequest(body)
	}

	if next == nil {
		next = item.Next()
	}

	if next == nil {

		log.WithFields(log.Fields{
			"key":   "analytics",
			"event": "get_mock",
			"payload": logrus.Fields{
				"success": false,
				"err":     "Not found mock",
				"url":     r.URL.String(),
			},
		}).Info("ANALYTICS")

		fields["Group URL"] = item.URL
		fields["Group Method"] = item.Method
		log.WithFields(fields).Warn("Not found any group")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "not found mock for url" + r.URL.String()})
		return
	}

	fields["Response URL"] = next.URL
	fields["Response Method"] = next.Method
	fields["Status code"] = next.StatusCode

	log.WithFields(fields).Info("Was Sended")

	log.WithFields(log.Fields{
		"key":   "analytics",
		"event": "get_mock",
		"payload": logrus.Fields{
			"success": true,
			"url":     r.URL.String(),
		},
	}).Info("ANALYTICS")

	w.WriteHeader(next.StatusCode)
	json.NewEncoder(w).Encode(next.Response)
}

func proxyRequest(r *http.Request, host string, scheme string) ([]byte, error) {

	// Выполняем проксирование с сохранением файла

	resp, err := startProxing(r, host, scheme)

	logFields := logrus.Fields{
		"host":   host,
		"scheme": scheme,
		"url":    r.URL.String(),
		"resp":   resp,
	}

	if err != nil {

		// Если проексирование завершилось c ошибкой, то возвращаем ее

		logFields["success"] = false
		logFields["err"] = err
		logAnalytics(logFields)
		return []byte{}, err
	}

	data, err := ioutil.ReadAll(resp.Body)

	if err == nil {

		// Если проксирование завершилось без ошибок и удалось почитать данные из ответа, то возвращаем их клиенту

		logFields["success"] = true
		logAnalytics(logFields)
		return data, nil
	}

	// Если произошла ошибка при считывании возвращаем ошибку

	logFields["success"] = false
	logFields["err"] = err

	logAnalytics(logFields)
	return []byte{}, err
}
