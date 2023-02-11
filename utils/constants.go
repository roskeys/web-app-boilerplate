package utils

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

var PRODUCTION bool
var SERVER_ADDRESS string
var DOMAIN_NAME string

var DB_USERNAME string
var DB_PASSWORD string
var DB_DATABASE string

var JWT_SECRET string

type Config struct {
	Server *struct {
		Production    *bool   `yaml:"production"`
		ServerAddress *string `yaml:"server_address"`
		DomainName    *string `yaml:"domain_name"`
	} `yaml:"server"`
	DB *struct {
		DbUsername *string `yaml:"db_username"`
		DbPassword *string `yaml:"db_password"`
		DbDatabase *string `yaml:"db_database"`
	} `yaml:"db"`
	JWTSecret *string `yaml:"jwt_secret"`
}

func LoadEnv() bool {
	var config Config
	configFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Fatalln("Failed to load config file: ", err)
	}

	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		log.Fatal("Failed to parse config file: ", err)
	}

	// SERVER
	PRODUCTION = *config.Server.Production
	SERVER_ADDRESS = *config.Server.ServerAddress
	DOMAIN_NAME = *config.Server.DomainName

	// DB
	DB_USERNAME = *config.DB.DbUsername
	DB_PASSWORD = *config.DB.DbPassword
	DB_DATABASE = *config.DB.DbDatabase

	// JWT
	JWT_SECRET = *config.JWTSecret
	return true
}

var _ = LoadEnv()
