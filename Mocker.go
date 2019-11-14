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
	"time"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

var models []mock.RequestModelGroup
var configuration config.Config

const update = "/update_models"

const isNeedProxyHeaderKey = "X-Mocker-Redirect-Is-On"
const redirectHostHeaderKey = "X-Mocker-Redirect-Host"
const redirectURLSchemeHeaderKey = "X-Mocker-Redirect-Scheme"
const specificPathHeaderKey = "X-Mocker-Specific-Path"

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

	specificPath := r.Header.Get(specificPathHeaderKey)

	if isNeedProxy == "true" && scheme != "" && host != "" {

		data, err := proxyRequest(r, host, scheme)

		if err == nil {
			// Если метод проксирования не вернул ошибки, то просто записывает ответ в response и заканчиваем обработку
			w.WriteHeader(http.StatusOK)
			w.Write(data)
			return
		}
	}

	// Проверяем, является ли полученный запрос запросом на обновление моделей
	// Если да, то обновляем модели и выходим.

	if strings.Compare(r.URL.String(), update) == 0 {

		err := startUpdateModels()

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"message": "Cant update models"})
			return
		}
	}

	// Если мы дошли сюда, то нужно найти нужный мок. Сначала ищем нужную группу

	item := mock.FindGroupByURL(&models, r.URL.String(), r.Method)

	if item == nil {

		// Если группа не найдена, то возвращаем ошибку

		logFields := log.Fields{
			"success":      false,
			"err":          "GroupNotFound",
			"requestedUrl": r.URL.String(),
		}

		logAnalytics(logFields, EventKeyGetMock)

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "not found mock for url" + r.URL.String()})
		return
	}

	var next *mock.RequestModel

	// Читаем тело запроса

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err == nil {
		// Если тело запроса считалось, то находим мок в группе, у которого значние request такое же
		next = item.CompareByRequest(body)
	}

	if next == nil {
		// Если не нашлось мока по телу - выдаем просто следующий по списку
		next = item.Next(specificPath)
	}

	if next == nil {

		// Если следующего нет -> группа пуста. Возвращаем ошибку

		logFields := log.Fields{
			"success":      false,
			"err":          "ItemInGroupNotFound",
			"requestedUrl": r.URL.String(),
			"groupUrl":     item.URL,
			"groupMethod":  item.Method,
		}

		logAnalytics(logFields, EventKeyGetMock)

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "not found mock for url" + r.URL.String()})
		return
	}

	// Записываем мок в ответ

	logFields := log.Fields{
		"success":      true,
		"requestedUrl": r.URL.String(),
	}

	logAnalytics(logFields, EventKeyGetMock)

	w.WriteHeader(next.StatusCode)
	json.NewEncoder(w).Encode(next.Response)
}

func proxyRequest(r *http.Request, host string, scheme string) ([]byte, error) {

	// Выполняем проксирование с сохранением файла

	logFields := logrus.Fields{
		"host":       host,
		"scheme":     scheme,
		"url":        r.URL.String(),
		"proxyStart": time.Now().Format(time.RFC3339),
	}

	resp, err := startProxing(r, host, scheme)

	logFields["resp"] = resp
	logFields["proxyEnd"] = time.Now().Format(time.RFC3339)

	defer resp.Body.Close()

	if err != nil {

		// Если проексирование завершилось c ошибкой, то возвращаем ее

		logFields["success"] = false
		logFields["err"] = err
		logAnalyticsProxy(logFields)
		return []byte{}, err
	}

	data, err := ioutil.ReadAll(resp.Body)

	if err == nil {
		// Если проксирование завершилось без ошибок и удалось почитать данные из ответа, то возвращаем их клиенту

		logFields["success"] = true
		logAnalyticsProxy(logFields)
		return data, nil
	}

	// Если произошла ошибка при считывании возвращаем ошибку

	logFields["success"] = false
	logFields["err"] = err

	logAnalyticsProxy(logFields)

	time.Now().Format(time.RFC3339)

	return []byte{}, err
}

func startUpdateModels() error {

	logFields := log.Fields{
		"success":   false,
		"startTime": time.Now().Format(time.RFC3339),
	}

	err := updateModels()

	logFields["endTime"] = time.Now().Format(time.RFC3339)

	if err != nil {

		logFields["success"] = false
		logFields["err"] = err

		logAnalytics(logFields, EventKeyUpdateModels)
		return err
	}

	logFields["proxyEnd"] = time.Now().Format(time.RFC3339)

	logAnalytics(logFields, EventKeyUpdateModels)
	return nil
}
