package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	selectTemplate "github.com/ockibagusp/golang-website-example/app/main/template"
)

func init() {
	// Templates: session controller
	templates := selectTemplate.AppendTemplates
	templates["login.html"] = selectTemplate.ParseFileHTMLOnly("views/login.html")
}

func (ctrl *Controller) Login(c echo.Context) error {
	if c.Request().Method == "POST" {
		return c.Redirect(http.StatusFound, "/")
	}

	return c.Render(http.StatusOK, "login.html", echo.Map{
		"csrf":         c.Get("csrf"),
		"flash_error":  []string{},
		"is_html_only": true,
	})
}
