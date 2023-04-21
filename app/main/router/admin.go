package router

import (
	"github.com/labstack/echo/v4"
	"github.com/ockibagusp/golang-website-example/app/main/controller"
)

func SetAdminRoutes(router *echo.Echo, controller *controller.Controller, jwtAuthMiddleware echo.MiddlewareFunc) (admin *echo.Group) {
	admin = router.Group("/admin", jwtAuthMiddleware)
	admin.GET("/delete-permanently", controller.DeletePermanently).
		Name = "/admin/delete-permanently get"
	admin.GET("/restore/:id", controller.RestoreUser).
		Name = "/admin/restore/:id get"
	// "/admin/delete-permanently/:id" unable
	// "/admin/delete/permanently/:id" can
	admin.GET("/delete/permanently/:id", controller.DeletePermanentlyByID).
		Name = "/admin/delete/permanently/:id get"

	return
}
