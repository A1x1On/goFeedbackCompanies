package config

import(
	"gov/backend/common/helper"
	"github.com/pkg/errors"
	"gov/backend/models"
	"path/filepath"
	"encoding/json"
	"io/ioutil"
)

var Set *models.ConfigModel

func Get(){
	jsonFileName  := filepath.Base("config.json")
	jsonData, err := ioutil.ReadFile(jsonFileName)
	if err != nil {
		helper.CheckError(errors.Wrap(err, "error read config file"))
	}

	err = json.Unmarshal(jsonData, &Set)
	if err != nil {
		helper.CheckError(errors.Wrap(err, "error parse json config file"))
	}
}