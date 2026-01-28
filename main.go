package main

import (
	"fmt"
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
	protected_controller.Config = cfg

	fmt.Println("init_config(): Configuration loaded successfully.")
}

func main() {
	var err error
	_ = err

	//Gin stuff
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	//Enforce Strict CSP
	r.Use(func(c *gin.Context) {
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; connect-src 'self';")
		c.Next()
	})

	//Bootstrap
	r.Static("/css/bootstrap", "./app/static/vendor/bootstrap/css")
	r.Static("/js/bootstrap", "./app/static/vendor/bootstrap/js")

	//JQuery
	r.Static("/js/jquery", "./app/static/vendor/jquery/js")

	//Select2
	r.Static("/css/select2", "./app/static/vendor/select2/css")
	r.Static("/js/select2", "./app/static/vendor/select2/js")

	//Datatables
	r.Static("/css/datatables", "./app/static/vendor/datatables/css")
	r.Static("/js/datatables", "./app/static/vendor/datatables/js")

	lib_controller.Router(r)
	protected_controller.Router(r)
	public_controller.Router(r)

	fmt.Println("Application running...")

	r.Run(":8080")
}
