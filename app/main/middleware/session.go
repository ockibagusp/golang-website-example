package middleware

import (
	"golang-website-example/config"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

var sessionsCookieStore = config.GetAPPConfig().SessionsCookieStore

func SessionNewCookieStore() echo.MiddlewareFunc {
	return session.Middleware(sessionNewCookieStore())
}

func sessionNewCookieStore() *sessions.CookieStore {
	return sessions.NewCookieStore(
		[]byte(sessionsCookieStore),
	)
}

/////
//	session for flash message.
////

const sessionFlash = "flash"

// SetFlash: set session for flash message
func SetFlash(c echo.Context, name, value string) {
	txSessionFlash := sessionFlash
	if name == "message" {
		txSessionFlash += "-message"
	} else if name == "error" {
		txSessionFlash += "-error"
	}

	session, _ := sessionNewCookieStore().Get(c.Request(), txSessionFlash)

	session.AddFlash(value, name)
	session.Options.MaxAge = 2
	session.Save(c.Request(), c.Response())
}

// GetFlash: get session for flash messages
func GetFlash(c echo.Context, name string) (flashes []string) {
	txSessionFlash := sessionFlash
	if name == "message" {
		txSessionFlash += "-message"
	} else if name == "error" {
		txSessionFlash += "-error"
	}

	session, _ := sessionNewCookieStore().Get(c.Request(), txSessionFlash)
	session.Options.MaxAge = 2

	fls := session.Flashes(name)
	if len(fls) > 0 {
		session.Save(c.Request(), c.Response())
		for _, fl := range fls {
			flashes = append(flashes, fl.(string))
		}

		return flashes
	}
	return nil
}

// SetFlashSuccess: set session for flash message: success
func SetFlashSuccess(c echo.Context, success string) {
	SetFlash(c, "success", success)
}

// GetFlashSuccess: get session for flash messages: []success
func GetFlashSuccess(c echo.Context) []string {
	return GetFlash(c, "success")
}

// SetFlashError: get session for flash message: error
func SetFlashError(c echo.Context, error string) {
	SetFlash(c, "error", error)
}

// GetFlashError: get session for flash message: []error
func GetFlashError(c echo.Context) []string {
	return GetFlash(c, "error")
}
