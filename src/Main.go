// `Main.go` server entry point
// contains `/` handler (`rootHandler`) which handle all requests
//
// How does it works:
//
// 1. Server is started
// 2. Server reads and applies config. Look at package `config` struct `Config`
// 3. Server configures logging. Look at `configureLog`
// 4. Server reads all mocks files. Look at `updateModels`.
//	  After reading server creates complicated data structure for comfortable searching for specific mock
//    This structure will be called `state`. The state is stored in variable `state`
//	  For more information look at package `mock` file `Mocks.go`
// 5. Server starts listening on specific port for any request. Look at `rootHandler`
//    Then each request will be handled via something like state machine which is described in `rootHandler`

package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"mocker/config"
	"mocker/mock"
	"net/http"
)

// state is the all mocks which may be requested form `Mocker`
// WARNING: This is a "global state"
var state []mock.RequestModelGroup

// configuration just a variable for storing configs
// you can change it in runtime, but you need to be careful
var configuration config.Config

func main() {
	conf, err := config.LoadConfig()

	if err != nil {
		log.Fatal(err)
	}

	configuration = conf

	configureLog(&conf)

	if err := updateModels(); err != nil {
		logFields := log.Fields{
			"errorMessage": err.Error(),
			"description":  "Can't updateUrlPath mocks in server startup",
		}

		logAnalytics(logFields, EventKeyUpdateModels)
	}

	http.HandleFunc("/", rootHandler)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", conf.Port), nil))
}

// rootHandler do:
// 1. Run specific features
// 2. Update models
// 3. Search specific mock for request
//
// if 1 is success then handler will return result.
// if it isn't success then handler will move to 2
//
// if 2 is success then handler will return result.
// if it isn't success then handler will move to 3
func rootHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	// firstly trying to handle special features
	// if method returns true then we need to stop handling and return result
	// which was written by `tryToHandleSpecificFeatures` to the client
	if tryToHandleSpecificFeatures(w, r) {
		return
	}

	// then the request is "state updater" then we will call to updateUrlPath state
	// look at `tryToUpdateModels` for more details
	if tryToUpdateModels(w, r) {
		return
	}

	// if we are there then we need to look up mock for the request
	lookUpAndWrite(w, r)
}
