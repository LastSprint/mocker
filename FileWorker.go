package main

import (
	"encoding/json"
	"io/ioutil"
	"mocker/mock"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

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
