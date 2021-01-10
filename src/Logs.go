package main

import (
	"mocker/config"
	"os"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

// EventKey keys for possible events
type EventKey string

const (
	EventKeyProxy             EventKey = "proxy"
	EventKeyGetMock           EventKey = "get_mock"
	EventKeyUpdateModels      EventKey = "update_models"
	EventKeyProxyFileSave     EventKey = "proxy_file_save"
	EventKeyCantWriteResponse EventKey = "cant_write_response"
)

// LogKey key for logging type
type LogKey string

const (
	LogKeyAnalytics LogKey = "analytics"
)

func configureLog(config *config.Config) {

	file, err := os.OpenFile(config.LogsPath, os.O_RDWR|os.O_CREATE, os.ModePerm)

	if err != nil {
		log.WithFields(log.Fields{
			"Action": "Not Found Log",
		}).Warning()
		log.SetOutput(os.Stdout)
		return
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
