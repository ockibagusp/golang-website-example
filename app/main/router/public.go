package router

import (
	"golang-website-example/app/main/controller"

	"github.com/labstack/echo/v4"
)

func SetPublicRoutes(router *echo.Echo, controller *controller.Controller, jwtAuthMiddleware echo.MiddlewareFunc) (public *echo.Group) {
	public = router.Group("", jwtAuthMiddleware)
	public.GET("/", controller.Home).Name = "home"
	public.GET("/about", controller.About).Name = "about"
	public.GET("/login", controller.Login).Name = "login get"
	public.POST("/login", controller.Login).Name = "login post"
	public.GET("/logout", controller.Logout).Name = "logout get"

	return
}
