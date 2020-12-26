package main

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const updateUrlPath = "/update_models"

const isNeedProxyHeaderKey = "X-Mocker-Redirect-Is-On"
const projectIdHeader = "X-Mocker-Project-Id"
const redirectHostHeaderKey = "X-Mocker-Redirect-Host"
const redirectURLSchemeHeaderKey = "X-Mocker-Redirect-Scheme"
const specificPathHeaderKey = "X-Mocker-Specific-Path"

// tryToHandleSpecificFeatures tries to handle some special features like proxying
// if them method can deal something with request it will return true.
//
// About result handling
// if the method returns `true` it means that the method did something and result was written to `w`
// If the method returns `false` it means that the method did nothing with `w` and you can continue processing
func tryToHandleSpecificFeatures(w http.ResponseWriter, r *http.Request) bool {
	host := r.Header.Get(redirectHostHeaderKey)
	scheme := r.Header.Get(redirectURLSchemeHeaderKey)
	isNeedProxy := r.Header.Get(isNeedProxyHeaderKey)
	projectID := r.Header.Get(projectIdHeader)

	if isNeedProxy == "true" && scheme != "" && host != "" {

		resp, err := proxyRequest(r, host, scheme, projectID)
		if err != nil {
			w.WriteHeader(http.StatusGone)
			dt, _ := json.Marshal(err)
			logAnalyticsProxy(log.Fields{"err": err, "resp": resp})
			wrap(log.Fields{"operation": "proxyRequest"}, EventKeyCantWriteResponse, func() error {
				_, err := w.Write(dt)
				return err
			})
			return true
		}
		data, err := ioutil.ReadAll(resp.Body)

		if err == nil {
			// if proxy doesnt return error then just write response
			for key, val := range resp.Header {
				for _, it := range val {
					w.Header().Add(key, it)
				}
			}
			w.WriteHeader(resp.StatusCode)
			wrap(log.Fields{"operation": "proxyRequest"}, EventKeyCantWriteResponse, func() error {
				_, err := w.Write(data)
				return err
			})
			return true
		}
	}

	return false
}

func proxyRequest(r *http.Request, host string, scheme string, projectID string) (*http.Response, error) {

	logFields := log.Fields{
		"host":       host,
		"scheme":     scheme,
		"url":        r.URL.String(),
		"proxyStart": time.Now().Format(time.RFC3339),
	}

	resp, err := startProxing(r, host, scheme, projectID)

	logFields["proxyEnd"] = time.Now().Format(time.RFC3339)

	if err != nil {
		logFields["success"] = false
		logFields["err"] = err
		logAnalyticsProxy(logFields)
		return nil, err
	}

	logFields["success"] = true
	logAnalyticsProxy(logFields)
	return resp, nil
}

// tryToUpdateModels understands is the request is request for updating state
// if it is then the method will updateUrlPath state (look at `FileWorker.go`) and return `true`
// if it isn't then the method will return `false`
func tryToUpdateModels(w http.ResponseWriter, r *http.Request) bool {
	if strings.Compare(r.URL.String(), updateUrlPath) == 0 {

		err := startUpdateModels()

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			wrap(log.Fields{"operation": "updateModels"}, EventKeyCantWriteResponse, func() error {
				return json.NewEncoder(w).Encode(map[string]string{"message": "Can't updateUrlPath state"})
			})
			return true
		}

		wrap(log.Fields{"operation": "updateModels"}, EventKeyCantWriteResponse, func() error {
			return json.NewEncoder(w).Encode(map[string]string{"message": "Success"})
		})
		return true
	}

	return false
}

// startUpdateModels starts reading mocks form storage (FS) and updates state
func startUpdateModels() error {

	logFields := log.Fields{
		"success":   true,
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

	logAnalytics(logFields, EventKeyUpdateModels)
	return nil
}
