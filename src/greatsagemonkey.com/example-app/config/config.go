package config

import (
	"encoding/json"
	"log"
	"os"
)

// NewConfig reads a config file in json format and extracts a map of string/string pairs
func NewConfig(filename string, configName string, version string, runtimeStage string) (map[string]string, error) {
	config := map[string]string{}

	configFile, err := os.Open(filename)
	if err != nil {
		return config, err
	}
	log.Println("Reading configuration from ", filename)

	// map is "configName" - "version" - "stage" - "key" - "value"
	allConfigs := map[string]map[string]map[string]map[string]string{}

	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&allConfigs)
	if err != nil {
		return config, err
	}

	config = allConfigs[configName][version][runtimeStage]

	return config, nil
}
