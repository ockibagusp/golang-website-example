package middleware

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/ockibagusp/golang-website-example/models"
)

// base.html -> {{if eq ((index .session.Values "is_auth_type") | tostring) -1 }}ok{{end}}

// GetAuth: get session from authentication
func GetAuth(c echo.Context) (session_gorilla *sessions.Session, err error) {
	if session_gorilla, err = session.Get("session", c); err != nil {
		return
	}

	if _, ok := session_gorilla.Values["username"]; !ok {
		session_gorilla.Values["username"] = ""
	}
	if _, ok := session_gorilla.Values["is_auth_type"]; !ok {
		session_gorilla.Values["is_auth_type"] = -1
	}

	return
}

// SetSession: set session from User
func SetSession(user models.User, c echo.Context) (session_gorilla *sessions.Session, err error) {
	session_gorilla, err = session.Get("session", c)
	if err != nil {
		return
	}

	session_gorilla.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7 days expired
		HttpOnly: true,
		Secure:   true,
	}

	session_gorilla.Values["username"] = user.Username
	// TODO: user.IsAuthType
	session_gorilla.Values["is_auth_type"] = 2 // TODO: admin: 1 and user: 2
	session_gorilla.Save(c.Request(), c.Response())

	return
}

// ClearSession: delete session from User
func ClearSession(c echo.Context) (err error) {
	var session_gorilla *sessions.Session
	if session_gorilla, err = session.Get("session", c); err != nil {
		return
	}

	session_gorilla.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
	}

	session_gorilla.Values["username"] = ""
	session_gorilla.Values["is_auth_type"] = -1
	session_gorilla.Save(c.Request(), c.Response())
	return
}

// RefreshSession: refresh session from User
func RefreshSession(user models.User, c echo.Context) (session_gorilla *sessions.Session, err error) {
	return
}

/////
//	session for flash message.
////

const sessionFlash = "flash"

// cookieStoreFlash: new cookie store session for flash
func cookieStoreFlash() *sessions.CookieStore {
	return sessions.NewCookieStore(
		[]byte("secret-session-key"),
	)
}

// SetFlash: set session for flash message
func SetFlash(c echo.Context, name, value string) {
	txSessionFlash := sessionFlash
	if name == "message" {
		txSessionFlash += "-message"
	} else if name == "error" {
		txSessionFlash += "-error"
	}

	session, _ := cookieStoreFlash().Get(c.Request(), txSessionFlash)

	session.AddFlash(value, name)
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

	session, _ := cookieStoreFlash().Get(c.Request(), txSessionFlash)

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
