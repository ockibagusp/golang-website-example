package controller

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/ockibagusp/golang-website-example/app/main/middleware"
	selectTemplate "github.com/ockibagusp/golang-website-example/app/main/template"
	"github.com/ockibagusp/golang-website-example/app/main/types"
	"github.com/ockibagusp/golang-website-example/business"
	"github.com/ockibagusp/golang-website-example/business/auth"
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
 * @method: GET
 * @route: /login
 */
func (ctrl *Controller) Login(c echo.Context) error {
	log := aulogger.Start(c)
	defer log.End()

	trackerID := aulogger.SetTrackerID()
	ic := business.NewInternalContext(trackerID)
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

		// Declare the expiration time of the token
		// here, we have kept it as 24 hour
		expirationTime := time.Now().Add(24 * time.Hour)

		// Create claims with multiple fields populated
		claims := auth.JwtClaims{
			UserID:   user.ID,
			Username: user.Username,
			Role:     user.Role,
			RegisteredClaims: jwt.RegisteredClaims{
				// A usual scenario is to set the expiration time relative to the current time
				ExpiresAt: jwt.NewNumericDate(expirationTime),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				NotBefore: jwt.NewNumericDate(time.Now()),
			},
		}

		// Declare the token with the algorithm used for signing, and the claims
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		// Create the JWT string
		tokenString, err := token.SignedString([]byte(ctrl.appConfig.AppJWTAuthSign))
		if err != nil {
			// If there is an error in creating the JWT return an internal server error
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"message": "server error",
			})
		}

		// Finally, we set the client cookie for "token" as the JWT we just generated
		// we also set an expiry time which is the same as the token itself
		cookie := new(http.Cookie)
		cookie.Name = "token"
		cookie.Value = tokenString
		cookie.Expires = expirationTime
		c.SetCookie(cookie)

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

	if err := middleware.ClearSession(c); err != nil {
		log.Warn("to middleware.Clearauth auth not found")
		// err: auth not found
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": err.Error(),
		})
	}

	log.Info("END request method GET for logout")
	return c.Redirect(http.StatusSeeOther, "/")
}
