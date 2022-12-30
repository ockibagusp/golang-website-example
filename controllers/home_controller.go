package controllers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ockibagusp/golang-website-example/middleware"
	"github.com/ockibagusp/golang-website-example/models"
	log "github.com/sirupsen/logrus"
)

/*
 * Home "home.html"
 *
 * @target: All
 * @method: GET
 * @route: /
 */
func (controller Controller) Home(c echo.Context) error {
	// Please note the the second parameter "home.html" is the template name and should
	// be equal to one of the keys in the TemplateRegistry array defined in main.go
	session, _ := middleware.GetAuth(c)
	log := log.WithFields(log.Fields{
		"username": session.Values["username"],
		"route":    c.Path(),
	})
	log.Info("START request method GET for home")

	var message string
	if session.Values["username"] != "" {
		var user models.User
		if err := controller.DB.Select("name").Where(
			"username = ?", session.Values["username"],
		).First(&user); err.Error != nil { // why?
			log.Warnf(`session values "username" error: %v`, err.Error)
		}

		message = fmt.Sprintf("%v!", user.Name)
		// or,
		// session.Values["username"].(string) + "!"
	}

	log.Info("END request method GET for home: [+]success")
	return c.Render(http.StatusOK, "home.html", echo.Map{
		"name":          "Home",
		"nav":           "home", // (?)
		"session":       session,
		"flash_success": middleware.GetFlashSuccess(c),
		"msg":           message,
	})
}
