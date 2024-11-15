package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/Builderhummel/thesis-app/app/Controllers/lib_controller"
	"github.com/Builderhummel/thesis-app/app/Controllers/protected_controller"
	"github.com/Builderhummel/thesis-app/app/Models/db_model"
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
	db_model.Config = cfg
}

func main() {
	//DB stuff
	var dbc db_model.DBController
	_, err := dbc.OpenConnection()
	if err != nil {
		log.Fatalf("could not open connection: %v", err)
	}
	defer dbc.CloseConnection()

	check, err := dbc.CheckIfDatabaseIsInitialized()
	if err != nil {
		log.Fatalf("could not check if database is initialized: %v", err)
	}
	if !check {
		err = dbc.InitDatabase()
		if err != nil {
			log.Fatalf("could not init database: %v", err)
		}
	}

	//Gin stuff
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	//Bootstrap
	r.Static("/css/bootstrap", "./app/static/vendor/bootstrap/css")
	r.Static("/js/bootstrap", "./app/static/vendor/bootstrap/js")

	protected_controller.Router(r)
	lib_controller.Router(r)

	r.Run(":8080")

}
