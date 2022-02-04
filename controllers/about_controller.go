package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ockibagusp/golang-website-example/middleware"
	log "github.com/sirupsen/logrus"
)

/*
 * About "about.html"
 *
 * @target: All
 * @method: GET
 * @route: /about
 */
func (Controller) About(c echo.Context) error {
	// Please note the the second parameter "about.html" is the template name and should
	// be equal to one of the keys in the TemplateRegistry array defined in main.go
	session, _ := middleware.GetAuth(c)
	log := log.WithFields(log.Fields{
		"username": session.Values["username"],
		"route":    c.Path(),
	})
	log.Info("START request method GET for about")
	// ---
	log.Info("END request method GET for about: [+]success")
	return c.Render(http.StatusOK, "about.html", echo.Map{
		"name":    "About",
		"nav":     "about", // (?)
		"session": session,
		"msg":     "All about Ocki Bagus Pratama!",
	})
}
