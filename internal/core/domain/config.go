package domain

import (
	"encoding/json"
	"os"
	"time"

	"github.com/rs/zerolog/log"
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

// MDBConfig holds Mongo DB related configurations
type MDBConfig struct {
	// Database host, default: 127.0.0.1
	Host string `json:"host,omitempty"`
	// Database port, default: 27017
	Port int `json:"port,omitempty"`
	// Database name
	DB string `json:"db,omitempty"`
	// Database username. Omit if the DB have no authentication
	Username string `json:"username,omitempty"`
	// Database password. Omit if the DB have no authentication
	Password string `json:"password,omitempty"`
	// Statements execution timeout, default: 10 seconds
	Timeout time.Duration `json:"timeout,omitempty"`
	// Database connection timeout, default: 10 seconds
	ConnectTimeout time.Duration `json:"connectTimeout,omitempty"`
	// MaxPoolSize tells how many connections are established with the database
	MaxPoolSize int `json:"maxPoolSize,omitempty"`
}

type Pagination struct {
	// Default page size if no specified
	PageSize int `json:"pageSize,omitempty"`
	// Max page size the user can ask
	MaxPageSize int `json:"maxPageSize,omitempty"`
}

type KeyList struct {
	// For auth related encryptions should be a 32 bits hex encoded string
	Auth string `json:"auth,omitempty"`
}

// Config contains all configuration for this service
type Config struct {
	CassandraDB   CDBConfig `json:"cassandraDB"`
	MongoDBConfig MDBConfig `json:"mongoDB"`
	// Server bind IP default 0.0.0.0
	Host string `json:"host,omitempty"`
	// Server bind port default 8080
	Port string `json:"port,omitempty"`
	// Default pagination settings
	Pagination Pagination `json:"pagination,omitempty"`
	// Encryption and security keys
	Keys KeyList `json:"keys,omitempty"`
}

// DefaultConfig returns an instance of Config with the default values
func DefaultConfig() Config {
	return Config{
		CassandraDB: CDBConfig{
			Host:           "127.0.0.1",
			Port:           9042,
			Timeout:        10,
			ConnectTimeout: 10,
			Connections:    2,
		},
		MongoDBConfig: MDBConfig{
			Host:           "127.0.0.1",
			Port:           27017,
			Timeout:        10,
			ConnectTimeout: 10,
			MaxPoolSize:    50,
		},
		Host: "0.0.0.0",
		Port: "8080",
		Pagination: Pagination{
			PageSize:    10,
			MaxPageSize: 100,
		},
	}
}

// LoadConfiguration Loads the configuration object from a json file
func LoadConfiguration(file string) Config {
	config := DefaultConfig()
	configFile, err := os.Open(file)

	if err != nil {
		log.Warn().Err(err).Msg("Can't load config file. Default values will be used instead")
	}

	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config
}
