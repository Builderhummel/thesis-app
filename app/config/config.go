package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Configuration struct {
	DBIP       string `json:"db_ip"`
	DBPort     string `json:"db_port"`
	DBUsername string `json:"db_username"`
	DBPassword string `json:"db_password"`
	DBName     string `json:"db_name"`
	LDAPUrl    string `json:"ldap_url"`
	LDAPDn     string `json:"ldap_dn"`
	JWTSecret  string `json:"jwt_secret"`
}

func LoadConfig() (*Configuration, error) {
	var cfg Configuration
	file, err := os.Open("config/config.json")
	if err != nil {
		return nil, fmt.Errorf("could not open config file: %v", err)
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("could not read config file: %v", err)
	}

	if err := json.Unmarshal(bytes, &cfg); err != nil {
		return nil, fmt.Errorf("could not unmarshal config: %v", err)
	}

	return &cfg, nil
}
