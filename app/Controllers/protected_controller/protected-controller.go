package protected_controller

import (
	"net/http"
	"strings"
	"text/template"

	"github.com/gin-gonic/gin"
)

func Router(r *gin.Engine) {
	r.SetFuncMap(template.FuncMap{
		"StringsJoin": strings.Join,
	})

	//r.LoadHTMLGlob("app/Views/protected/**/*.html")
	r.LoadHTMLGlob("app/Views/**/**/*.html")

	//r.LoadHTMLFiles("app/Views/protected/homepage/index.tmpl")
	//r.LoadHTMLFiles("app/Views/protected/common/header.tmpl")

	r.GET("/", func(ctx *gin.Context) {
		tabOpReq := generateTORTestData()
		tabMySupV := generateTMSupVTestData()
		ctx.HTML(http.StatusOK, "protected/homepage/index.html", gin.H{
			"TabOpReq":  tabOpReq,
			"TabMySupV": tabMySupV,
		})
	})
}

/**
func index(c *gin.Context) {
	c.HTML(http.StatusOK, "protected/homepage/index.html", gin.H{})
}
**/

//TODO: find out how namespaces work in gin
//template homepage/index.html, common/header.tmpl
