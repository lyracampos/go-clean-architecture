package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type (
	Config struct {
		App      App
		API      API
		Database Database
	}

	App struct {
		Name string
	}

	API struct {
		Host string
		Port int
	}

	Database struct {
		ConnectionString   string
		MaxOpenConnections int
		MaxIdleConnections int
	}
)

func NewConfig(configFilePath string) (*Config, error) {
	config, err := readConfig(configFilePath)
	if err != nil {
		return &Config{}, fmt.Errorf("failed to get viper config: %w", err)
	}

	return &Config{
		App: App{
			Name: config.GetString("app.name"),
		},
		API: API{
			Host: config.GetString("api.http.host"),
			Port: config.GetInt("api.http.port"),
		},
		Database: Database{
			ConnectionString:   config.GetString("database.connectionString"),
			MaxOpenConnections: config.GetInt("database.maxOpenConnections"),
			MaxIdleConnections: config.GetInt("database.maxIdleConnections"),
		},
	}, nil
}

func readConfig(configFilePath string) (*viper.Viper, error) {
	config := viper.New()
	config.SetEnvPrefix("clean_architecture")
	config.SetConfigType("yaml")
	config.SetConfigFile(configFilePath)
	config.AutomaticEnv()

	if err := config.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	return config, nil
}
