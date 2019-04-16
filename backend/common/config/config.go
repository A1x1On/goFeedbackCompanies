package config

import(
	"gov/backend/common/helper"
	"github.com/pkg/errors"
	"gov/backend/models"
	"path/filepath"
	"encoding/json"
	"io/ioutil"
)

var config *models.ConfigModel

func Get() *models.ConfigModel{
	jsonFileName  := filepath.Base("config.json")
	jsonData, err := ioutil.ReadFile(jsonFileName)
	if err != nil {
		helper.CheckError(errors.Wrap(err, "error read config file"))
	}

	err = json.Unmarshal(jsonData, &config)
	if err != nil {
		helper.CheckError(errors.Wrap(err, "error parse json config file"))
	}

	return config
}