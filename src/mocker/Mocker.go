package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mocklistner/config"
	"mocklistner/mock"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var models []mock.RequestModelGroup
var configuration config.Config

func main() {
	log.Println("Start")

	conf, err := config.LoadConfig(os.Args[1])

	if err != nil {
		log.Fatal(err)
	}

	configuration = conf

	updateModels()

	http.HandleFunc("/", handler)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", conf.Port), nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	log.Println(r.URL.String())

	if strings.Compare(r.URL.String(), "/update_models") == 0 {
		err := updateModels()
		if err != nil {

		}
		return
	}

	item := mock.FindGroupByURL(&models, r.URL.String(), r.Method)

	if item == nil {
		log.Println("NOT FOUND GROUP", r.URL)
		return
	}

	next := item.Next()

	if next == nil {
		log.Println("GROUP IS EMPTY", r.URL)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "not found mock for url" + r.URL.String()})
		return
	}

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
