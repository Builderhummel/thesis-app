package protected_controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Router(r *gin.Engine) {
	//r.LoadHTMLGlob("app/Views/protected/**/*.html")
	r.LoadHTMLGlob("app/Views/**/**/*.html")

	//r.LoadHTMLFiles("app/Views/protected/homepage/index.tmpl")
	//r.LoadHTMLFiles("app/Views/protected/common/header.tmpl")

	r.GET("/", index)
}

func index(c *gin.Context) {
	c.HTML(http.StatusOK, "protected/homepage/index.html", gin.H{})
}

//TODO: find out how namespaces work in gin
//template homepage/index.html, common/header.tmpl
