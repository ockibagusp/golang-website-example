package router

import (
	"os"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/ockibagusp/golang-website-example/controllers"
	"github.com/ockibagusp/golang-website-example/template"
	"github.com/sirupsen/logrus"
)

// Router init
func New(controllers *controllers.Controller) (router *echo.Echo) {
	// Echo instance
	router = echo.New()

	// Middleware
	router.Use(middleware.Logger())
	router.Use(middleware.Recover())

	/*
		Insyaallah, TODO: .env session_test: {1} or {0}
				   session_test: {true} or {false}

		1. os.Setenv("session_test", ...)
		@session_test: {true} or {1}
		os.Setenv("session_test", "true") or,
		os.Setenv("session_test", "1")

		@session_test: {false} or {0}
		os.Setenv("session_test", "false") or,
		os.Setenv("session_test", "0")
	*/
	logrus.Println("Setenv: session_test = session")
	router.Use(session.Middleware(sessions.NewCookieStore(
		[]byte("something-very-secret"),
	)))

	// PROD
	if os.Getenv("session_test") == "0" || os.Getenv("session_test") == "false" {
		logrus.Println(`Setenv: "session_test" != "0" || "session_test" != "false"`)
		router.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
			// Optional. Default value "header:X-CSRF-Token".
			// Possible values:
			// - "header:<name>"
			// - "form:<name>"
			// - "query:<name>"
			TokenLookup: "form:X-CSRF-Token",
		}))
	}

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

	// admin
	router.GET("/admin/delete-permanently", controllers.DeletePermanently).
		Name = "/admin/delete-permanently get"
	router.GET("/admin/restore/:id", controllers.RestoreUser).
		Name = "/admin/restore/:id get"
	// "/admin/delete-permanently/:id" can not
	// "/admin/delete/permanently/:id" can
	router.GET("/admin/delete/permanently/:id", controllers.DeletePermanentlyByID).
		Name = "/admin/delete/permanently/:id get"

	return
}
