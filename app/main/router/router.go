package router

import (
	"github.com/labstack/echo/v4"
	"github.com/ockibagusp/golang-website-example/app/main/controller"
	"github.com/ockibagusp/golang-website-example/app/main/template"
	"github.com/ockibagusp/golang-website-example/config"
)

func RegisterPath(
	e *echo.Echo,
	appConfig *config.Config,
	controller *controller.Controller,
) {
	// Why bootstrap.min.css, bootstrap.min.js, jquery.min.js?
	e.Static("/assets", "assets")

	// Instantiate a template registry with an array of template set
	e.Renderer = template.NewTemplates()

	e.GET("/", controller.Home).Name = "home"
	e.GET("/about", controller.About).Name = "about"
	e.GET("/login", controller.Login).Name = "login get"
	e.POST("/login", controller.Login).Name = "login post"
	e.GET("/users", controller.Users).Name = "users"
	e.GET("/users/add", controller.CreateUser).Name = "user/add get"
	e.POST("/users/add", controller.CreateUser).Name = "user/add post"
}
