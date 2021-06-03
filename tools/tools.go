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

func StringEquals(first *string, second *string) bool {
	if first == nil && second == nil {
		return true
	}
	if (first == nil && second != nil) || (first != nil && second == nil) {
		return false
	}
	if *first == *second {
		return true
	}
	return false
}

// compares to pointers and returns false if references cannot be compared
func ObjectComparableInDepth(first interface{}, second interface{}) bool {
	if first == nil && second == nil {
		return false
	}
	if (first == nil && second != nil) || (first != nil && second == nil) {
		return false
	}
	return true
}
