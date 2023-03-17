package router

import (
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"

	"github.com/ockibagusp/golang-website-example/app/main/controller"
	"github.com/ockibagusp/golang-website-example/app/main/middleware"
	"github.com/ockibagusp/golang-website-example/app/main/template"
	"github.com/ockibagusp/golang-website-example/config"
)

func RegisterPath(
	appConfig *config.Config,
	controller *controller.Controller,
) (router *echo.Echo) {
	router = echo.New()
	router.Pre(echoMiddleware.RemoveTrailingSlash())

	// Middleware
	router.Use(echoMiddleware.Logger())
	router.Use(echoMiddleware.Recover())
	router.Use(middleware.SessionNewCookieStore())

	// Why bootstrap.min.css, bootstrap.min.js, jquery.min.js?
	router.Static("/assets", "assets")

	// Instantiate a template registry with an array of template set
	router.Renderer = template.NewTemplates()

	sessionMiddleware := middleware.SessionMiddleware()

	// public
	public := router.Group("", sessionMiddleware)
	public.GET("/", controller.Home).Name = "home"
	public.GET("/about", controller.About).Name = "about"
	public.GET("/login", controller.Login).Name = "login get"
	public.POST("/login", controller.Login).Name = "login post"
	public.GET("/logout", controller.Logout).Name = "logout get"

	// user
	user := router.Group("/users", sessionMiddleware)
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

	return
}
