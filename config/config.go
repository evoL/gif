package config

import (
	"encoding/json"
	"os"
	"os/user"
	"path"
)

type Config struct {
	storePath string
	db        struct {
		driver     string
		dataSource string
	}
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

func (c *Config) DbDriver() string {
	if c.db.driver == "" {
		return "sqlite3"
	}
	return c.db.driver
}

func (c *Config) DbDataSource() string {
	if c.db.dataSource == "" {
		current, _ := user.Current()
		defaultPath := path.Join(current.HomeDir, ".gif", "gif.db")
		return defaultPath
	}
	return c.db.dataSource
}

func StorePath() string {
	return globalConfig.StorePath()
}

func DbDriver() string {
	return globalConfig.DbDriver()
}

func DbDataSource() string {
	return globalConfig.DbDataSource()
}
