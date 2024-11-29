package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/Builderhummel/thesis-app/app/Controllers/auth_controller"
	"github.com/Builderhummel/thesis-app/app/Controllers/lib_controller"
	"github.com/Builderhummel/thesis-app/app/Controllers/protected_controller"
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
	/**
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
	*/

	/*
		//LDAP TEST
		var auser auth_controller.AuthUser
		err = auser.LDAP_authenticate(os.Getenv("LUSERNAME"), os.Getenv("LPASSWORD"))
		if err != nil {
			if ldapErr, ok := err.(*ldap.Error); ok {
				switch ldapErr.ResultCode {
				case ldap.LDAPResultInvalidCredentials:
					println("Hi")
					log.Fatalf("invalid credentials: %v", ldapErr)
				}
			}
			log.Fatalf("could not authenticate: %v", err)
		}
		log.Printf("UID: %v", auser.UID)
		log.Printf("Name: %v", auser.Name)
		log.Printf("Email: %v", auser.Email)

		os.Exit(0)
		//LDAP TEST
	*/

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
