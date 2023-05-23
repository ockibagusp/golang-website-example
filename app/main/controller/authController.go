package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ockibagusp/golang-website-example/app/main/helpers"
	"github.com/ockibagusp/golang-website-example/app/main/middleware"
	selectTemplate "github.com/ockibagusp/golang-website-example/app/main/template"
	"github.com/ockibagusp/golang-website-example/app/main/types"
	"github.com/ockibagusp/golang-website-example/business"
	log "github.com/ockibagusp/golang-website-example/logger"
)

var aulogger = log.NewPackage("auth_controller")

func init() {
	// Templates: auth controller
	templates := selectTemplate.AppendTemplates
	templates["login.html"] = selectTemplate.ParseFileHTMLOnly("views/login.html")
}

/*
 * Auth: Login
 *
 * @target: All
 * @method: GET and POST
 * @route: /login
 */
func (ctrl *Controller) Login(c echo.Context) error {
	log := aulogger.Start(c)
	defer log.End()

	trackerID := aulogger.SetTrackerID()
	ic := business.NewInternalContext(trackerID)
	if c.Request().Method == http.MethodPost {
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

		user, validPassword := ctrl.authService.VerifyLogin(ic, passwordForm.Username, passwordForm.Password)
		if !validPassword {
			// or, middleware.SetFlashError(c, "username or password not match")
			middleware.SetFlash(c, "error", "username or password not match")

			log.Warn("for database `username` or `password` not nil for login")
			log.Warn("END request method POST for login: [-]failure")
			return c.Render(http.StatusOK, "login.html", echo.Map{
				"csrf":         c.Get("csrf"),
				"flash_error":  middleware.GetFlashError(c),
				"is_html_only": true,
			})
		}

		if err := middleware.SetCookie(c, user, ctrl.appConfig.AppJWTAuthSign); err != nil {
			// If there is an error in creating the JWT return an internal server error
			return c.JSON(http.StatusInternalServerError, helpers.Response{
				Code:   http.StatusInternalServerError,
				Status: "Internal Server Error",
				Data:   err,
			})
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
 * auth: Logout
 *
 * @target: Users
 * @method: GET
 * @route: /logout
 */
func (ctrl *Controller) Logout(c echo.Context) error {
	log := aulogger.Start(c)
	defer log.End()

	log.Info("START request method GET for logout")

	middleware.ClearCookie(c)

	log.Info("END request method GET for logout")
	return c.Redirect(http.StatusSeeOther, "/")
}
