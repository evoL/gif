package config

import (
	"encoding/json"
	"os"
	"os/user"
	"path"
)

type Config struct {
	storePath string
}

var globalConfig *Config

func New(path string) (*Config, error) {
	return loadConfig(path, false)
}

func Load(path string) error {
	var err error
	globalConfig, err = New(path)
	return err
}

func NewDefault() (*Config, error) {
	current, _ := user.Current()
	defaultPath := path.Join(current.HomeDir, ".gifconfig")

	return loadConfig(defaultPath, true)
}

func Default() error {
	var err error
	globalConfig, err = NewDefault()
	return err
}

func loadConfig(configPath string, passIfMissing bool) (*Config, error) {
	file, err := os.Open(configPath)
	if err != nil {
		if passIfMissing && os.IsNotExist(err) {
			return &Config{}, nil
		}
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	cfg := &Config{}

	if err = decoder.Decode(&cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) StorePath() string {
	if c.storePath == "" {
		current, _ := user.Current()
		defaultPath := path.Join(current.HomeDir, ".gif", "store")
		return defaultPath
	}
	return c.storePath
}

func StorePath() string {
	return globalConfig.StorePath()
}
