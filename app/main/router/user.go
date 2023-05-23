package router

import (
	"github.com/labstack/echo/v4"
	"github.com/ockibagusp/golang-website-example/app/main/controller"
)

func SetUserRoutes(router *echo.Echo, controller *controller.Controller, jwtAuthMiddleware echo.MiddlewareFunc) (user *echo.Group) {
	user = router.Group("/users", jwtAuthMiddleware)
	user.GET("", controller.Users).Name = "users"
	user.GET("/add", controller.CreateUser).Name = "user/add get"
	user.POST("/add", controller.CreateUser).Name = "user/add post"
	user.GET("/read/:id", controller.ReadUser).Name = "user/read get"
	user.GET("/view/:id", controller.UpdateUser).Name = "user/view get"
	user.POST("/view/:id", controller.UpdateUser).Name = "user/view post"
	user.GET("/view/:id/password", controller.UpdateUserByPassword).
		Name = "user/view/:id/password get"
	user.POST("/view/:id/password", controller.UpdateUserByPassword).
		Name = "user/view/:id/password post"
	user.GET("/delete/:id", controller.DeleteUser).Name = "user/delete get"

	return
}
