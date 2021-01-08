package main

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"mocker/features"
	"mocker/mock"
	"net/http"
)

// lookUpAndWrite searches specific mock by specific parameters
// it will write something to `w` in any case
//
// 1. Search group of mocks which have equal url as was requested.
// 2. In group search specific mock (by parameters (like headers, body)
// 3. If mock was found then will return it
// 4. If wasn't it will return mock from group (for details look at `RequestModelGroup`)
// 5. If group is empty will return error
func lookUpAndWrite(w http.ResponseWriter, r *http.Request) {

	specificPath := r.Header.Get(specificPathHeaderKey)

	// firstly, try to find group of mocks which have the same url as requested

	item := mock.FindGroupByURL(&state, r.URL.String(), r.Method)

	if item == nil {

		// if item is nil it means that mocks group with requested url is not found
		// in other words, it means that there are not mocks with requested url
		// for more details look at implementation `mock.FindGroupByURL`

		logFields := log.Fields{
			"success":            false,
			"err":                "GroupNotFound",
			"specificHeaderPath": specificPath,
			"requestedUrl":       r.URL.String(),
		}

		writeNotFoundErr(w, r, logFields, EventKeyGetMock)
		return
	}

	var next *mock.RequestModel

	// if we found group then we needs to find specific mock in this group

	flatHeaders := map[string]string{}

	for key, value := range r.Header {
		flatValue := ""

		for _, headerItem := range value {
			flatValue += headerItem
		}

		flatHeaders[key] = flatValue
	}

	// read request body

	body, _ := ioutil.ReadAll(r.Body)

	defer func() {
		if err := r.Body.Close(); err != nil {
			logFields := log.Fields{
				"err":                "SystemError",
				"requestedUrl":       r.URL.String(),
				"groupUrl":           item.URL,
				"specificHeaderPath": specificPath,
				"groupMethod":        item.Method,
				"headers":            flatHeaders,
				"request":            string(body),
				"description":        "Can't close request body",
				"errorMessage":       err.Error(),
			}

			logAnalytics(logFields, EventKeyGetMock)
		}
	}()

	next = item.LookUpByBodyAndHeaders(body, flatHeaders)

	if next == nil {
		logFields := log.Fields{
			"err":                "SearchingWasFail",
			"requestedUrl":       r.URL.String(),
			"groupUrl":           item.URL,
			"specificHeaderPath": specificPath,
			"groupMethod":        item.Method,
			"headers":            flatHeaders,
			"request":            string(body),
		}
		logAnalytics(logFields, EventKeyGetMock)
		// if we can't find specific mock (by parameters)
		// then we will just return `next` mock
		next = item.Next(specificPath)
	}

	if next == nil {
		// if at this step next is empty it means that the group is empty
		// end we will return error
		logFields := log.Fields{
			"success":            false,
			"err":                "ItemInGroupNotFound",
			"requestedUrl":       r.URL.String(),
			"groupUrl":           item.URL,
			"specificHeaderPath": specificPath,
			"groupMethod":        item.Method,
		}
		writeNotFoundErr(w, r, logFields, EventKeyGetMock)
		return
	}

	// write mock to response

	logFields := log.Fields{
		"success":            true,
		"specificHeaderPath": specificPath,
		"requestedUrl":       r.URL.String(),
	}

	logAnalytics(logFields, EventKeyGetMock)

	if next.Delay != 0 {
		delayer := features.Throttler{}
		delayer.Throttle(next.Delay)
	}

	for key, value := range next.ResponseHeaders {
		w.Header().Set(key, value)
	}

	w.WriteHeader(next.StatusCode)

	wrap(log.Fields{"operation": "lookUpAndWrite"}, EventKeyCantWriteResponse, func() error {
		return json.NewEncoder(w).Encode(next.Response)
	})
}

func writeNotFoundErr(w http.ResponseWriter, r *http.Request, payload log.Fields, event EventKey) {
	logAnalytics(payload, event)
	w.WriteHeader(http.StatusBadRequest)
	wrap(payload, event, func() error {
		return json.NewEncoder(w).Encode(map[string]string{"message": "not found mock for url" + r.URL.String()})
	})
}
