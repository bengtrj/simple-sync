package config

import (
	"errors"
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

const (
	knownStateFilePath = "config/.known.yml"

	//DesiredStateFilePath path to the desired state file
	DesiredStateFilePath = "config/desired.yml"
)

//Sync is the main configuration structure
//It holds the Known and Desired states
type Sync struct {
	User         string
	Password     string
	DesiredState *State
	KnownState   *State
}

//State abstracts either a desired or known host state
type State struct {
	Servers []Server `yaml:"servers"`
	Apps    []App    `yaml:"apps"`
}

//Server a simple representation of the hosts
type Server struct {
	IP string `yaml:"ip"`
}

//App represents a service/package that should be present on the remote hosts
type App struct {
	Name     string    `yaml:"name"`
	Packages []Package `yaml:"packages"`
	Files    []File    `yaml:"files"`
}

//Package represents a package that should be present on the remote hosts
type Package struct {
	Name      string `yaml:"name"`
	IsService bool   `yaml:"is-service"`
}

//File represents a file that should be present on the remote hosts
type File struct {
	Path    string `yaml:"path"`
	Content string `yaml:"content"`
	Owner   string `yaml:"owner"`
	Group   string `yaml:"group"`
	Mode    int    `yaml:"mode"`
}

//Load parses the config file into a Sync structure
func Load() (*Sync, error) {

	password := os.Getenv("PASSWORD")
	if password == "" {
		return nil,
			errors.New("Root password not provided. Expecting a environment valiable PASSWORD")
	}

	desired, err := parse(DesiredStateFilePath)
	if err != nil {
		return nil, err
	}

	var known *State
	if fileExists(knownStateFilePath) {
		known, err = parse(knownStateFilePath)
		if err != nil {
			return nil, err
		}
	}

	return &Sync{
		User:         "root",
		Password:     password,
		DesiredState: desired,
		KnownState:   known,
	}, nil
}

func parse(path string) (*State, error) {

	var config State

	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil

}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
