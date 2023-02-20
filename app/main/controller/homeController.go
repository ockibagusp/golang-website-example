package controller

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	selectTemplate "github.com/ockibagusp/golang-website-example/app/main/template"
)

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
func (Controller) Home(c echo.Context) error {
	// Please note the the second parameter "home.html" is the template name and should
	// be equal to one of the keys in the TemplateRegistry array defined in main.go
	// ?
	username, _ := c.Get("username").(string)
	role, _ := c.Get("role").(string)

	return c.Render(http.StatusOK, "home.html", echo.Map{
		"name":             "Home",
		"nav":              "home", // (?)
		"session_username": username,
		"session_role":     role,
		"flash_success":    []string{},
		"msg":              fmt.Sprintf("%v!", "Ocki Bagus Pratama"),
	})
}
