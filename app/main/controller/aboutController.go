package controller

import (
	"net/http"

	"github.com/gorilla/sessions"
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
func (Controller) About(c echo.Context) error {
	// Please note the the second parameter "about.html" is the template name and should
	// be equal to one of the keys in the TemplateRegistry array defined in main.go
	session := sessions.Session{}

	return c.Render(http.StatusOK, "about.html", echo.Map{
		"name":    "About",
		"nav":     "about", // (?)
		"session": session,
		"msg":     "All about Ocki Bagus Pratama!",
	})
}
