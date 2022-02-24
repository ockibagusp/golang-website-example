package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ockibagusp/golang-website-example/middleware"
	"github.com/ockibagusp/golang-website-example/models"
	"github.com/ockibagusp/golang-website-example/types"
	log "github.com/sirupsen/logrus"
)

/*
 * Session: Login
 *
 * @target: All
 * @method: GET
 * @route: /login
 */
func (controller *Controller) Login(c echo.Context) error {
	session, err := middleware.GetAuth(c)
	log := log.WithFields(log.Fields{
		"username": session.Values["username"],
		"route":    c.Path(),
	})
	if session.Values["is_auth_type"] != -1 && err == nil {
		log.Info("START request method GET for login")
		log.Warn("to [@route: /] session")
		log.Warn("END request method GET for login: [-]failure")
		return c.Redirect(http.StatusFound, "/")
	}

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

		var user models.User
		// err := controller.DB.Select(...).Where(...).First(...).Error
		if err := controller.DB.Select("username", "password", "is_admin").Where(
			"username = ?", passwordForm.Username,
		).First(&user).Error; err != nil {
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

	log.Info("START request method GET for login")
	log.Info("END request method GET for login")
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
func (controller *Controller) Logout(c echo.Context) error {
	session, _ := middleware.GetAuth(c)
	log := log.WithFields(log.Fields{
		"username": session.Values["username"],
		"route":    c.Path(),
	})
	log.Info("START request method GET for logout")

	if err := middleware.ClearSession(c); err != nil {
		log.Warn("to middleware.ClearSession session not found")
		// err: session not found
		return c.HTML(http.StatusBadRequest, err.Error())
	}

	log.Info("END request method GET for logout")
	return c.Redirect(http.StatusSeeOther, "/")
}
