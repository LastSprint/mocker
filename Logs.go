package main

import (
	"mocker/config"
	"os"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

// EventKey тип для ключей возможных событий
type EventKey string

const (
	// EventKeyProxy ключ события проксирования запроса
	EventKeyProxy EventKey = "proxy"
	// EventKeyGetMock ключ события получения запроса на чтение мока
	EventKeyGetMock EventKey = "get_mock"
	// EventKeyUpdateModels ключ события на обновление всех моделей моков (чтение файлов из fs)
	EventKeyUpdateModels EventKey = "update_models"
	// EventKeyProxyFileSave ключ для события записи проксированного ответа в моковый файл
	EventKeyProxyFileSave EventKey = "proxy_file_save"
)

// LogKey ключ для типа лога
type LogKey string

const (
	// LogKeyAnalytics ключ для логирования аналитики
	LogKeyAnalytics LogKey = "analytics"
)

func configureLog(config *config.Config) {

	file, err := os.OpenFile(config.LogsPath, os.O_RDWR|os.O_CREATE, os.ModePerm)

	if err != nil {
		log.WithFields(log.Fields{
			"Action": "Not Found Log",
		}).Panic()
	}

	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetOutput(file)
}

func logAnalyticsProxy(payload logrus.Fields) {
	logAnalytics(payload, EventKeyProxy)
}

func logAnalytics(payload logrus.Fields, event EventKey) {
	log.WithFields(log.Fields{
		"key":     LogKeyAnalytics,
		"event":   event,
		"payload": payload,
	}).Info("ANALYTICS")
}

func logEmptyAnalytics(event EventKey) {
	logAnalytics(logrus.Fields{}, event)
}
