package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/ockibagusp/golang-website-example/app/main/middleware"
	selectTemplate "github.com/ockibagusp/golang-website-example/app/main/template"
	"github.com/ockibagusp/golang-website-example/app/main/types"
	"github.com/ockibagusp/golang-website-example/business"
)

func init() {
	// Templates: session controller
	templates := selectTemplate.AppendTemplates
	templates["login.html"] = selectTemplate.ParseFileHTMLOnly("views/login.html")
}

/*
 * Session: Login
 *
 * @target: All
 * @method: GET
 * @route: /login
 */
func (ctrl *Controller) Login(c echo.Context) error {
	if c.Request().Method == "POST" {
		log.Info("START request method POST for login")
		passwordForm := &types.LoginForm{
			Username: c.FormValue("username"),
			Password: c.FormValue("password"),
		}

		err := passwordForm.Validate()
		if err != nil {
			middleware.SetFlashError(c, err.Error())

			log.Warn("for passwordForm.Validate() not nil for login")
			log.Warn("END request method POST for login: [-]failure")
			return c.Render(http.StatusOK, "login.html", echo.Map{
				"csrf":         c.Get("csrf"),
				"flash_error":  middleware.GetFlashError(c),
				"is_html_only": true,
			})
		}

		user, err := ctrl.userService.FirstUserByUsername(
			business.InternalContext{}, passwordForm.Username,
		)
		if err != nil {
			middleware.SetFlashError(c, err.Error())

			log.Warn("for database `username` or `password` not nil for login")
			log.Warn("END request method POST for login: [-]failure")
			return c.Render(http.StatusOK, "login.html", echo.Map{
				"csrf":         c.Get("csrf"),
				"flash_error":  middleware.GetFlashError(c),
				"is_html_only": true,
			})
		}

		// check hash password:
		// match = true
		// match = false
		if !middleware.CheckHashPassword(user.Password, passwordForm.Password) {
			// or, middleware.SetFlashError(c, "username or password not match")
			middleware.SetFlash(c, "error", "username or password not match")

			log.Warn("to check wrong hashed password for login")
			log.Warn("END request method POST for login: [-]failure")
			return c.Render(http.StatusForbidden, "login.html", echo.Map{
				"csrf":         c.Get("csrf"),
				"flash_error":  middleware.GetFlash(c, "error"),
				"is_html_only": true,
			})
		}

		if _, err := middleware.SetSession(user, c); err != nil {
			middleware.SetFlashError(c, err.Error())

			log.Warn("to middleware.SetSession session not found for login")
			log.Warn("END request method POST for login: [-]failure")
			// err: session not found
			return c.HTML(http.StatusForbidden, err.Error())

		}

		log.Info("END request method POST [@route: /]")
		return c.Redirect(http.StatusFound, "/")
	}

	return c.Render(http.StatusOK, "login.html", echo.Map{
		"csrf":         c.Get("csrf"),
		"flash_error":  middleware.GetFlashError(c),
		"is_html_only": true,
	})
}

/*
 * Session: Logout
 *
 * @target: Users
 * @method: GET
 * @route: /logout
 */
func (ctrl *Controller) Logout(c echo.Context) error {
	log.Info("START request method GET for logout")

	if err := middleware.ClearSession(c); err != nil {
		log.Warn("to middleware.ClearSession session not found")
		// err: session not found
		return c.HTML(http.StatusBadRequest, err.Error())
	}

	log.Info("END request method GET for logout")
	return c.Redirect(http.StatusSeeOther, "/")
}
