package controller

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ockibagusp/golang-website-example/app/main/middleware"
	selectTemplate "github.com/ockibagusp/golang-website-example/app/main/template"
	"github.com/ockibagusp/golang-website-example/business"
	log "github.com/ockibagusp/golang-website-example/logger"
)

var hlogger = log.NewLogger()

func init() {
	// Templates: homeController
	selectTemplate.AppendTemplates["home.html"] = selectTemplate.ParseFilesBase("views/home.html")
}

/*
 * Home "home.html"
 *
 * @target: All
 * @method: GET
 * @route: /
 */
func (ctrl *Controller) Home(c echo.Context) error {
	// Please note the the second parameter "home.html" is the template name and should
	// be equal to one of the keys in the TemplateRegistry array defined in main.go
	// ?
	log := hlogger.Start(c)
	defer log.End()

	log.Info("START request method GET for home")

	uid, _ := c.Get("uid").(uint)
	username, _ := c.Get("username").(string)
	role, _ := c.Get("role").(string)

	var message string
	if uid != 0 {
		user, err := ctrl.userService.FindByID(business.InternalContext{}, uid)
		if err != nil {
			log.Warnf(`session values "username" error: %v`, err)
		}

		message = fmt.Sprintf("%v!", user.Name)
	}

	log.Info("END request method GET for home: [+]success")
	return c.Render(http.StatusOK, "home.html", echo.Map{
		"name":             "Home",
		"nav":              "home", // (?)
		"session_username": username,
		"session_role":     role,
		"flash_success":    middleware.GetFlashSuccess(c),
		"message":          message,
	})
}
