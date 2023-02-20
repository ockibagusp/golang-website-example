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
	e *echo.Echo,
	appConfig *config.Config,
	controller *controller.Controller,
) {
	// Middleware
	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())
	e.Use(middleware.SessionNewCookieStore())

	// Why bootstrap.min.css, bootstrap.min.js, jquery.min.js?
	e.Static("/assets", "assets")

	// Instantiate a template registry with an array of template set
	e.Renderer = template.NewTemplates()

	sessionMiddleware := middleware.SessionMiddleware()

	// public
	public := e.Group("", sessionMiddleware)
	public.GET("/", controller.Home).Name = "home"
	public.GET("/about", controller.About).Name = "about"
	public.GET("/login", controller.Login).Name = "login get"
	public.POST("/login", controller.Login).Name = "login post"

	// user
	user := e.Group("/users", sessionMiddleware)
	user.GET("", controller.Users).Name = "users"
	user.GET("/add", controller.CreateUser).Name = "user/add get"
	user.POST("/add", controller.CreateUser).Name = "user/add post"
}
