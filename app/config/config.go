package config

import (
	"encoding/json"
	"errors"
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
		if errors.Is(err, os.ErrNotExist) {
			//TODO: Make this more beautiful, implemented this only for docker!!!
			return loadFromEnvvar(), nil
		}
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

func loadFromEnvvar() *Configuration {
	var cfg Configuration
	cfg.DBIP = os.Getenv("DB_IP")
	cfg.DBPort = os.Getenv("DB_PORT")
	cfg.DBUsername = os.Getenv("DB_USERNAME")
	cfg.DBPassword = os.Getenv("DB_PASSWORD")
	cfg.DBName = os.Getenv("DB_NAME")
	cfg.LDAPUrl = os.Getenv("LDAP_URL")
	cfg.LDAPDn = os.Getenv("LDAP_DN")
	cfg.JWTSecret = os.Getenv("JWT_SECRET")

	return &cfg
}
