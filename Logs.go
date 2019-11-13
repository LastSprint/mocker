package main

import (
	"mocker/config"
	"os"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

type EventKey string

const (
	EventKeyProxy EventKey = "proxy"
)

type LogKey string

const (
	LogKeyAnalytics LogKey = "analytics"
)

func configureLog(config *config.Config) {

	file, err := os.OpenFile(config.LogsPath, os.O_RDWR|os.O_CREATE, os.ModePerm)

	if err != nil {
		log.WithFields(log.Fields{
			"Action": "Not Found Log",
		}).Panic()
	}

	log.SetFormatter(&logrus.TextFormatter{})
	log.SetOutput(file)
}

func logAnalytics(payload logrus.Fields) {
	log.WithFields(log.Fields{
		"key":     LogKeyAnalytics,
		"event":   EventKeyProxy,
		"payload": payload,
	}).Info("ANALYTICS")
}
