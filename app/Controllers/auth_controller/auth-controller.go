package auth_controller

import (
	"github.com/Builderhummel/thesis-app/app/Models/db_model"
	"github.com/Builderhummel/thesis-app/app/config"
	"github.com/gin-gonic/gin"
)

var Config *config.Configuration

type AuthUser struct {
	UID   string
	Name  string
	Email string
}

func Login(c *gin.Context) {
	auser := AuthUser{}

	userid := c.PostForm("userid")
	password := c.PostForm("password")

	println("Username:" + userid)

	err := auser.LDAP_authenticate(userid, password)
	if err != nil {
		c.JSON(401, gin.H{"error": "Invalid credentials"}) //TODO: Add login page with error message
		return
	}

	//Check DB and update data on DB
	userAuthorized, err := auser.CheckUserAuthorized()
	if err != nil {
		c.JSON(500, gin.H{"error": "Error checking user"}) //TODO: Add login page with error message
		return
	}
	if !userAuthorized {
		c.JSON(401, gin.H{"error": "User not authorized"}) //TODO: Add login page with error message
		return
	}

	//Update user data in DB
	err = db_model.UpdateUser(auser.UID, auser.Name, auser.Email)
	if err != nil {
		//TODO: What now?!? -> data could not be updated in db
	}

	//Create JWT
	jwt_token, err := GenerateToken(auser.UID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error generating token"}) //TODO: Add login page with error message
		return
	}

	//Return JWT as cookie
	c.SetCookie("token", jwt_token, int(3*60*60), "/", "", false, true) //TODO: Set cookie only true to HTTPS, fix domain

	//Redirect to protected homepage
	c.Redirect(302, "/") //TODO: Is this true?
}

func (auser *AuthUser) CheckUserAuthorized() (bool, error) {
	userExists, err := db_model.VerifyLoginUser(auser.UID)
	if err != nil {
		return false, err
	}
	if !userExists {
		return false, nil
	}

	userActive, err := db_model.CheckUserActive(auser.UID)
	if err != nil {
		return false, err
	}
	if !userActive {
		return false, nil
	}

	return true, nil
}
