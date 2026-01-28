package protected_controller

import (
	"github.com/Builderhummel/thesis-app/app/Constants/roles"
	middleware "github.com/Builderhummel/thesis-app/app/Middleware"
	"github.com/gin-gonic/gin"
)

func Router(r *gin.Engine) {
	// Protected (User is logged in)
	protected := r.Group("/")
	protected.Use(middleware.JwtAuthMiddleware())
	protected.Use(middleware.RoleMiddleware())

	protected.GET("/", RenderHomepage)

	protected.GET("/all", RenderAllSupervisions)
	protected.StaticFile("/js/custom/all/datatables.js", "./app/Views/templates/protected/all_supervisions/datatables.js")

	protected.GET("/my_supervisions", RenderMySupervisions)
	protected.StaticFile("/js/custom/my_supervisions/datatables.js", "./app/Views/templates/protected/my_supervisions/datatables.js")

	protected.GET("/open_requests", RenderOpenRequests)
	protected.StaticFile("/js/custom/open_requests/datatables.js", "./app/Views/templates/protected/open_requests/datatables.js")

	protected.GET("/add", RenderAddSupervisionRequestForm)
	protected.POST("/add", HandleNewSupervisionRequest)

	protected.GET("/delete", RenderDeleteSupervisionRequestForm)
	protected.POST("/delete", HandleDeleteThesisRequest)

	protected.GET("/view", RenderViewSupervisionRequestForm)
	protected.StaticFile("/js/custom/view_supervision_request/file-management.js", "./app/Views/templates/protected/view_supervision_request/file-management.js")

	protected.GET("/edit", RenderEditSupervisionRequestForm)
	protected.POST("/edit", HandleEditSupervisionRequest)
	protected.StaticFile("/js/custom/edit/select2.js", "./app/Views/templates/protected/edit_supervison_request/select2.js")

	protected.GET("/export", HandleExport)

	// File management routes
	protected.POST("/files/upload", HandleFileUpload)
	protected.GET("/files/download", HandleFileDownload)
	protected.GET("/files/list", HandleFileList)
	protected.DELETE("/files/delete", HandleFileDelete)

	// Admin (User has role "administrator")
	admin_route := protected.Group("/admin")
	admin_route.Use(middleware.RequireRole(roles.RoleAdministrator))

	admin_route.GET("/users", RenderListAllUsers)

	admin_route.GET("/new_user", RenderNewUser)
	admin_route.POST("/new_user", HandlePostNewUser)

	admin_route.GET("/edit_user", RenderEditUser)
	admin_route.POST("/edit_user", HandlePostEditUser)
}
