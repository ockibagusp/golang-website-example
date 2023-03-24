package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	selectTemplate "github.com/ockibagusp/golang-website-example/app/main/template"
)

func init() {
	// Templates: aboutController
	selectTemplate.AppendTemplates["about.html"] = selectTemplate.ParseFilesBase("views/about.html")
}

/*
 * About "about.html"
 *
 * @target: All
 * @method: GET
 * @route: /about
 */
func (ctrl *Controller) About(c echo.Context) error {
	// Please note the the second parameter "about.html" is the template name and should
	// be equal to one of the keys in the TemplateRegistry array defined in main.go
	ctrl.logger.Info("START request method GET for about")
	username, _ := c.Get("username").(string)
	role, _ := c.Get("role").(string)

	ctrl.logger.Info("END request method GET for about: [+]success")
	return c.Render(http.StatusOK, "about.html", echo.Map{
		"name":             "About",
		"nav":              "about", // (?)
		"session_username": username,
		"session_role":     role,
		"message":          "All about Ocki Bagus Pratama!",
	})
}
