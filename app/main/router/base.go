package router

import (
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"

	"golang-website-example/app/main/controller"
	"golang-website-example/app/main/middleware"
	"golang-website-example/app/main/template"
	"golang-website-example/config"
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

	jwtAuthMiddleware := middleware.JwtAuthMiddleware(appConfig.AppJWTAuthSign)

	// public
	SetPublicRoutes(router, controller, jwtAuthMiddleware)
	// user
	SetUserRoutes(router, controller, jwtAuthMiddleware)
	// admin
	SetAdminRoutes(router, controller, jwtAuthMiddleware)

	return
}
