package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mocker/config"
	"mocker/mock"
	"net/http"
	"os"
	"path/filepath"
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

func configureLog(config *config.Config) {
	fmt.Println(config)
	file, err := os.OpenFile(config.LogsPath, os.O_RDWR|os.O_CREATE, os.ModePerm)

	if err != nil {
		log.WithFields(log.Fields{
			"Action": "Not Found Log",
		}).Panic()
	}

	log.SetFormatter(&logrus.TextFormatter{})
	log.SetOutput(file)
}

func handler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	host := r.Header.Get(redirectHostHeaderKey)
	scheme := r.Header.Get(redirectURLSchemeHeaderKey)
	isNeedProxy := r.Header.Get(isNeedProxyHeaderKey)

	if isNeedProxy == "true" && scheme != "" && host != "" {

		resp, err := startProxing(r, host, scheme)

		if err != nil {
			log.WithFields(log.Fields{
				"key":   "analytics",
				"event": "proxy",
				"payload": logrus.Fields{
					"success": false,
					"err":     err,
					"host":    host,
					scheme:    scheme,
					"url":     r.URL.String(),
					"resp":    resp,
				},
			}).Info("ANALYTICS")
		} else {
			data, err := ioutil.ReadAll(resp.Body)
			if err == nil {
				w.WriteHeader(http.StatusOK)
				w.Write(data)

				log.WithFields(log.Fields{
					"key":   "analytics",
					"event": "proxy",
					"payload": logrus.Fields{
						"success": true,
						"host":    host,
						scheme:    scheme,
						"url":     r.URL.String(),
						"resp":    resp,
					},
				}).Info("ANALYTICS")

				return
			}
			log.WithFields(log.Fields{
				"key":   "analytics",
				"event": "proxy",
				"payload": logrus.Fields{
					"success": false,
					"err":     err,
					"host":    host,
					scheme:    scheme,
					"url":     r.URL.String(),
					"resp":    resp,
				},
			}).Info("ANALYTICS")
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

func readModels() ([]mock.RequestModelGroup, error) {

	allMocks, err := readAllMocks()

	if err != nil {
		return []mock.RequestModelGroup{}, err
	}

	return mock.MakeGroups(allMocks), nil
}

func readAllMocks() ([]mock.RequestModel, error) {
	var models []mock.RequestModel

	err := filepath.Walk(configuration.MocksRootDir, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			log.Println(err)
			return err
		}

		if info.IsDir() {
			return nil
		}

		dat, err := ioutil.ReadFile(path)

		if err != nil {
			return err
		}
		var model mock.RequestModel

		err = json.Unmarshal(dat, &model)

		if err != nil {
			log.Println("CANT PARSE", path, err)
			return nil
		}
		models = append(models, model)
		return nil
	})

	if err != nil {
		return []mock.RequestModel{}, err
	}

	return models, err
}

func updateModels() error {
	newModels, err := readModels()
	if err != nil {
		return err
	}
	models = newModels
	return nil
}
