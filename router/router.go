package router

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/ockibagusp/golang-website-example/controllers"
	"github.com/ockibagusp/golang-website-example/template"
)

// Router init
func New(controllers *controllers.Controller) (router *echo.Echo) {
	// Echo instance
	router = echo.New()

	// Middleware
	router.Use(middleware.Logger())
	router.Use(middleware.Recover())
	// TODO: .env cookie store ?
	router.Use(session.Middleware(sessions.NewCookieStore(
		[]byte("something-very-secret"),
	)))
	router.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		// Optional. Default value "header:X-CSRF-Token".
		// Possible values:
		// - "header:<name>"
		// - "form:<name>"
		// - "query:<name>"
		TokenLookup: "form:X-CSRF-Token",
	}))

	// Instantiate a template registry with an array of template set
	router.Renderer = template.NewTemplates()

	// Why bootstrap.min.css, bootstrap.min.js, jquery.min.js?
	router.Static("/assets", "assets")

	// Router => controllers
	router.GET("/", controllers.Home).Name = "home"
	router.GET("/login", controllers.Login).Name = "login get"
	router.POST("/login", controllers.Login).Name = "login post"
	router.GET("/logout", controllers.Logout).Name = "logout get"
	router.GET("/about", controllers.About).Name = "about"
	router.GET("/users", controllers.Users).Name = "users"
	router.GET("/users/add", controllers.CreateUser).Name = "user/add get"
	router.POST("/users/add", controllers.CreateUser).Name = "user/add post"
	router.GET("/users/read/:id", controllers.ReadUser).Name = "user/read get"
	router.GET("/users/view/:id", controllers.UpdateUser).Name = "user/view get"
	router.POST("/users/view/:id", controllers.UpdateUser).Name = "user/view post"
	router.GET("/users/view/:id/password", controllers.UpdateUserByPassword).
		Name = "user/view/:id/password get"
	router.POST("/users/view/:id/password", controllers.UpdateUserByPassword).
		Name = "user/view/:id/password post"
	router.GET("/users/delete/:id", controllers.DeleteUser).Name = "user/delete get"

	return
}
