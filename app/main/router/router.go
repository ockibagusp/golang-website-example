package router

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/ockibagusp/golang-website-example/app/main/controller"
	"github.com/ockibagusp/golang-website-example/app/main/template"
	"github.com/ockibagusp/golang-website-example/config"
)

func RegisterPath(
	e *echo.Echo,
	appConfig *config.Config,
	controller *controller.Controller,
) {
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Why bootstrap.min.css, bootstrap.min.js, jquery.min.js?
	e.Static("/assets", "assets")

	// Instantiate a template registry with an array of template set
	e.Renderer = template.NewTemplates()

	e.GET("/", controller.Users)
	e.GET("/login", controller.Login).Name = "login get"
	e.POST("/login", controller.Login).Name = "login post"
}
