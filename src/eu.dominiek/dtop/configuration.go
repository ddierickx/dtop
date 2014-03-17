package main

import (
	"os"
	"errors"
	"fmt"
	"encoding/json"
)

// Holds the application's configuration settings.
type DTopConfiguration struct {
	Users []DTopUser
	Name string
	Description string
	StaticFolder string
	Port int
}

// Defines a dtop user.
type DTopUser struct {
	Username string
	Password string
}

func NewDTopUser(username string, password string) *DTopUser {
	user := new(DTopUser)
	user.Username = username
	user.Password = password
	return user
}

// Constructor for DTopConfiguration
func NewDTopConfiguration(name string, description string, users []DTopUser, staticFolder string, port int) *DTopConfiguration {
	cfg := new(DTopConfiguration)
	cfg.Name = name
	cfg.Description = description
	cfg.Users = users
	cfg.StaticFolder = staticFolder
	cfg.Port = port
	return cfg
}

// Validate the DTopConfiguration instance.
func (cfg *DTopConfiguration) IsValid() (bool, error) {
	if cfg.Name == "" {
		return false, errors.New(fmt.Sprintf("No name defined."))
	}

	if cfg.Description == "" {
		return false, errors.New(fmt.Sprintf("No description defined."))
	}

	if cfg.Port < 1 || cfg.Port > 65536 {
		return false, errors.New(fmt.Sprintf("Invalid port number: %s", cfg.Port))
	}

	if _, err := os.Stat(cfg.StaticFolder); os.IsNotExist(err) {
		return false, errors.New(fmt.Sprintf("Static folder does not exist: %s", cfg.StaticFolder))
	}

	return true, nil
}

// Performs serialization to JSON of a DTopConfiguration.
func SerializeDTopConfigurationToJson(cfg *DTopConfiguration) ([]byte, error) {
	return json.Marshal(cfg)
}

// Performs deserialization from JSON of a DTopConfiguration.
func DeserializeDTopConfigurationFromJson(jsonData []byte) (*DTopConfiguration, error) {
	cfg := new(DTopConfiguration)
	err := json.Unmarshal(jsonData, &cfg)
	return cfg, err
}