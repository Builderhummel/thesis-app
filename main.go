package main

import (
	"fmt"
	"log"

	dbcontroller "github.com/Builderhummel/thesis-app/app/Controllers/db-controller"
	"github.com/Builderhummel/thesis-app/app/config"
)

func init() {
	init_config()
}

func init_config() {
	var cfg *config.Configuration
	var err error
	cfg, err = config.LoadConfig()
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	//Config Injections here
	dbcontroller.Config = cfg
}

func main() {
	var dbc dbcontroller.DBController
	_, err := dbc.OpenConnection()
	if err != nil {
		log.Fatalf("could not open connection: %v", err)
	}
	defer dbc.CloseConnection()

	//err = dbc.InitDatabase()
	check, err := dbc.CheckIfDatabaseIsInitialized()
	fmt.Println(check)
	if err != nil {
		log.Fatalf("could not init database: %v", err)
	}
}
