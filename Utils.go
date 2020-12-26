package main

import "github.com/sirupsen/logrus"

// wrap call `fn` and if it returns error then wrap will add
// payload["errDescription"] = err.Error()
// and log this event to analytics
func wrap(payload logrus.Fields, event EventKey, fn func() error) {
	if err := fn(); err != nil {
		payload["errDescription"] = err.Error()
		logAnalytics(payload, event)
	}
}
