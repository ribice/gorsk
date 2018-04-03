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
	Server *ServerConfig
	DB     *DBConfig
	JWT    *JWTConfig
}

// DBConfig holds data necessery for database configuration
type DBConfig struct {
	PSN          string
	Log          bool
	CreateSchema bool
}

// ServerConfig holds data necessery for server configuration
type ServerConfig struct {
	Port int
}

// JWTConfig holds data necessery for JWT configuration
type JWTConfig struct {
	Realm            string
	Secret           string
	Duration         int
	RefreshDuration  int
	MaxRefresh       int
	SigningAlgorithm string
}
