package tools

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func LoadConfigFile(fileName string, target interface{}) {
	jsonFile, err := os.Open(fileName)
	if err != nil {
		panic("Something wrong with config file. " + err.Error())
	}
	defer jsonFile.Close()
	fileContent, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(fileContent, target)
}
