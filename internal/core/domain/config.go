package domain

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// CDBConfig holds Cassandra DB related configurations
type CDBConfig struct {
	// Database host, default: 127.0.0.1
	Host string `json:"host,omitempty"`
	// Database port, default: 9042
	Port int `json:"port,omitempty"`
	// Database username. Omit if the DB have no authentication
	Username string `json:"username,omitempty"`
	// Database password. Omit if the DB have no authentication
	Password string `json:"password,omitempty"`
	// CQL statements execution timeout, default: 10 seconds
	Timeout time.Duration `json:"timeout,omitempty"`
	// Database connection timeout, default: 10 seconds
	ConnectTimeout time.Duration `json:"connectTimeout,omitempty"`
	// Number of connections per host, default: 2
	Connections int `json:"connections,omitempty"`
}

// Config contains all configuration for this service
type Config struct {
	CassandraDB CDBConfig `json:"cassandraDB"`
	// Server bind IP default 0.0.0.0
	Host string `json:"host,omitempty"`
	// Server bind port default 8080
	Port string `json:"port,omitempty"`
}

// LoadConfiguration Loads the configuration object from a json file
func LoadConfiguration(file string) Config {
	config := Config{
		CassandraDB: CDBConfig{
			Host:           "127.0.0.1",
			Port:           9042,
			Timeout:        10,
			ConnectTimeout: 10,
			Connections:    2,
		},
		Host: "0.0.0.0",
		Port: "8080",
	}
	configFile, err := os.Open(file)

	if err != nil {
		fmt.Println(err.Error())
	}

	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config
}
