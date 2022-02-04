package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ockibagusp/golang-website-example/middleware"
	log "github.com/sirupsen/logrus"
)

/*
 * Home "home.html"
 *
 * @target: All
 * @method: GET
 * @route: /
 */
func (Controller) Home(c echo.Context) error {
	// Please note the the second parameter "home.html" is the template name and should
	// be equal to one of the keys in the TemplateRegistry array defined in main.go
	session, _ := middleware.GetAuth(c)
	log := log.WithFields(log.Fields{
		"username": session.Values["username"],
		"route":    c.Path(),
	})
	log.Info("START request method GET for home")
	// ---
	log.Info("END request method GET for home: [+]success")
	return c.Render(http.StatusOK, "home.html", echo.Map{
		"name":    "Home",
		"nav":     "home", // (?)
		"session": session,
		"msg":     "Ocki Bagus Pratama!",
	})
}
