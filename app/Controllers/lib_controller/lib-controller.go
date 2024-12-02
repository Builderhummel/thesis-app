package lib_controller

import (
	"strings"
	"text/template"

	"github.com/gin-gonic/gin"
)

func Router(r *gin.Engine) {
	r.SetFuncMap(template.FuncMap{
		"StringsJoin": strings.Join,
	})
	r.LoadHTMLGlob("app/Views/templates/**/**/*.html")
}
