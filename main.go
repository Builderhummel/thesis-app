package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/Builderhummel/thesis-app/app/Controllers/auth_controller"
	"github.com/Builderhummel/thesis-app/app/Controllers/lib_controller"
	"github.com/Builderhummel/thesis-app/app/Controllers/protected_controller"
	"github.com/Builderhummel/thesis-app/app/Controllers/public_controller"
	"github.com/Builderhummel/thesis-app/app/Models/db_model"
	"github.com/Builderhummel/thesis-app/app/config"
)

func init() {
	init_config()
	db_model.Init()
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
	auth_controller.Config = cfg
}

func main() {
	var err error
	_ = err

	//Gin stuff
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	//Bootstrap
	r.Static("/css/bootstrap", "./app/static/vendor/bootstrap/css")
	r.Static("/js/bootstrap", "./app/static/vendor/bootstrap/js")

	//JQuery
	r.Static("/js/jquery", "./app/static/vendor/jquery/js")

	//Select2
	r.Static("/css/select2", "./app/static/vendor/select2/css")
	r.Static("/js/select2", "./app/static/vendor/select2/js")

	lib_controller.Router(r)
	protected_controller.Router(r)
	public_controller.Router(r)

	r.Run(":8080")

}
