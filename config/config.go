package config

import (
	"encoding/json"
	"os"
	"os/user"
	"path"
)

type Config struct {
	StorePath string
	Db        struct {
		Driver     string
		DataSource string
	}
	Upload struct {
		Provider    string
		Credentials map[string]string
	}
}

var Global *Config

func New(path string) (*Config, error) {
	return loadConfig(path, false)
}

func Load(path string) error {
	var err error
	Global, err = New(path)
	return err
}

func NewDefault() (*Config, error) {
	current, _ := user.Current()
	defaultPath := path.Join(current.HomeDir, ".gifconfig")

	return loadConfig(defaultPath, true)
}

func Default() error {
	var err error
	Global, err = NewDefault()
	return err
}

func defaultConfig() (config *Config) {
	config = &Config{}
	current, _ := user.Current()

	config.StorePath = path.Join(current.HomeDir, ".gif", "store")
	config.Db.Driver = "sqlite3"
	config.Db.DataSource = path.Join(current.HomeDir, ".gif", "gif.db")
	config.Upload.Provider = "imgur"

	return
}

func loadConfig(configPath string, passIfMissing bool) (*Config, error) {
	file, err := os.Open(configPath)
	if err != nil {
		if passIfMissing && os.IsNotExist(err) {
			return defaultConfig(), nil
		}
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	cfg := defaultConfig()

	if err = decoder.Decode(&cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
