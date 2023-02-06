package controller

import (
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	selectTemplate "github.com/ockibagusp/golang-website-example/app/main/template"
)

func init() {
	// Templates: userController
	selectTemplate.AppendTemplates["home.html"] = selectTemplate.ParseFilesBase("views/home.html")
}

/*
 * Home "home.html"
 *
 * @target: All
 * @method: GET
 * @route: /
 */
func (ctrl Controller) Home(c echo.Context) error {
	// Please note the the second parameter "home.html" is the template name and should
	// be equal to one of the keys in the TemplateRegistry array defined in main.go
	// ?
	session := sessions.Session{}

	return c.Render(http.StatusOK, "home.html", echo.Map{
		"name":          "Home",
		"nav":           "home", // (?)
		"session":       session,
		"flash_success": []string{},
		"msg":           fmt.Sprintf("%v!", "Ocki Bagus Pratama"),
	})
}
