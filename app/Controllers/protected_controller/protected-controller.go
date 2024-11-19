package protected_controller

import (
	"net/http"
	"strings"
	"text/template"

	view_protected_homepage "github.com/Builderhummel/thesis-app/app/Views/handler/protected/homepage"
	"github.com/gin-gonic/gin"
)

func Router(r *gin.Engine) {
	r.SetFuncMap(template.FuncMap{
		"StringsJoin": strings.Join,
	})

	//r.LoadHTMLGlob("app/Views/protected/**/*.html")
	r.LoadHTMLGlob("app/Views/templates/**/**/*.html")

	//r.LoadHTMLFiles("app/Views/protected/homepage/index.tmpl")
	//r.LoadHTMLFiles("app/Views/protected/common/header.tmpl")

	tor := view_protected_homepage.NewTableOpenRequests()
	tor.AddRow("BA", "Max Mustermann", "Lorem Ipsum", "1630000000", "SoSe22", "request", "mailto:", "#", "#")
	tor.AddRow("BA", "Anna Schmidt", "Understanding Data Science", "1622505600", "WiSe22/23", "request", "mailto:", "#", "#")
	tor.AddRow("MA", "John Doe", "AI for Supply Chain Optimization", "1622505600", "WiSe22/23", "request", "mailto:", "#", "#")
	tor.AddRow("PA", "Jane Doe", "Sustainability in Cloud Computing", "1622505600", "WiSe22/23", "request", "mailto:", "#", "#")

	tmsv := view_protected_homepage.NewTableMySupervisions()
	tmsv.AddRow("BA", "Max Mustermann", "Analysis of Big Data Techniques", "1650470400", []string{"Prof. Müller", "Dr. Klein"}, "WiSe23/24", "working", "bg-success", "mailto:", "#", "#")
	tmsv.AddRow("MA", "Anna Schmidt", "Sustainability in Cloud Computing", "1653264000", []string{"Dr. Fischer", "Prof. Zhang"}, "SoSe24", "completed", "bg-secondary", "mailto:", "#", "#")
	tmsv.AddRow("PA", "John Doe", "AI for Supply Chain Optimization", "1653264000", []string{"Dr. Fischer", "Prof. Zhang"}, "SoSe24", "completed", "bg-secondary", "mailto:", "#", "#")

	r.GET("/", func(ctx *gin.Context) {
		tabOpReq := tor
		tabMySupV := tmsv
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
