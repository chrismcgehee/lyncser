package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path"

	yaml "gopkg.in/yaml.v3"
)

const (
	// Holds state that helps determine whether a file should be uploaded or downloaded.
	stateFilePath    = "~/.config/lyncser/state.json"
	// Contains global configuration used across all machines associated with this user.
	globalConfigPath = "~/.config/lyncser/globalConfig.yaml"
	// Contains configuration specific to this machine.
	localConfigPath  = "~/.config/lyncser/localConfig.yaml"
)

type GlobalConfig struct {
	// Specifies which files should be synced for machines associated with each tag. The key in this map is the tag
	// name. The value is the list of files/directories that should be synced for that tag.
	TagPaths map[string][]string `yaml:"paths"`
}

type LocalConfig struct {
	// Specifies with tags this machine should be associated with.
	Tags []string `yaml:"tags"`
}

type StateData struct {
	// Key is file path. Value is the state data associated with that file.
	FileStateData map[string]*FileStateData
}

type FileStateData struct {
	// The last time this file has been uploaded/downloaded from the cloud.
	LastCloudUpdate string
}

// getGlobalConfig reads and parses the global config file. If it does not exist, it will create it.
func getGlobalConfig() GlobalConfig {
	fullConfigPath := realPath(globalConfigPath)
	data, err := ioutil.ReadFile(fullConfigPath)
	if errors.Is(err, os.ErrNotExist) {
		configDir := path.Dir(fullConfigPath)
		os.MkdirAll(configDir, 0700)
		data = []byte("files:\n  all:\n    # - ~/.bashrc\n")
		err = os.WriteFile(fullConfigPath, data, 0644)
		panicError(err)
	} else {
		panicError(err)
	}
	var config GlobalConfig
	err = yaml.Unmarshal(data, &config)
	panicError(err)
	return config
}

// getLocalConfig reads and parses the local config file. If it does not exist, it will create it.
func getLocalConfig() LocalConfig {
	fullConfigPath := realPath(localConfigPath)
	data, err := ioutil.ReadFile(fullConfigPath)
	if errors.Is(err, os.ErrNotExist) {
		configDir := path.Dir(fullConfigPath)
		os.MkdirAll(configDir, 0700)
		data = []byte("tags:\n  - all\n")
		err = os.WriteFile(fullConfigPath, data, 0644)
		panicError(err)
	} else {
		panicError(err)
	}
	var config LocalConfig
	err = yaml.Unmarshal(data, &config)
	panicError(err)
	return config
}

// getStateData reads and parses the state data file. If that file does not exist yet, this method will return
// a newly initialized struct.
func getStateData() StateData {
	var stateData StateData
	data, err := ioutil.ReadFile(realPath(stateFilePath))
	if errors.Is(err, os.ErrNotExist) {
		stateData = StateData{
			FileStateData: map[string]*FileStateData{},
		}
	} else {
		panicError(err)
		err = json.Unmarshal(data, &stateData)
		panicError(err)
	}
	return stateData
}

// saveStateData will save the state data to disk.
func saveStateData(stateData StateData) {
	data, err := json.MarshalIndent(stateData, "", " ")
	panicError(err)
	err = ioutil.WriteFile(realPath(stateFilePath), data, 0644)
	panicError(err)
}
