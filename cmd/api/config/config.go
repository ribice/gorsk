package config

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"runtime"

	yaml "gopkg.in/yaml.v2"
)

// Load returns Configuration struct
func Load(env string) (*Configuration, error) {
	_, filePath, _, _ := runtime.Caller(0)
	configFile := filePath[:len(filePath)-9]
	bytes, err := ioutil.ReadFile(configFile +
		"files" + string(filepath.Separator) + "config." + env + ".yaml")
	if err != nil {
		return nil, fmt.Errorf("error reading config file, %s", err)
	}
	var cfg = new(Configuration)
	if err := yaml.Unmarshal(bytes, cfg); err != nil {
		return nil, fmt.Errorf("unable to decode into struct, %v", err)
	}
	return cfg, nil
}

// Configuration holds data necessery for configuring application
type Configuration struct {
	Server *Server
	DB     *Database
	JWT    *JWT
}

// Database holds data necessery for database configuration
type Database struct {
	PSN          string
	Log          bool
	CreateSchema bool
}

// Server holds data necessery for server configuration
type Server struct {
	Port int
}

// JWT holds data necessery for JWT configuration
type JWT struct {
	Realm            string
	Secret           string
	Duration         int
	RefreshDuration  int
	MaxRefresh       int
	SigningAlgorithm string
}
