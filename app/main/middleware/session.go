package middleware

import (
	"net/http"
	"strings"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	modelsTest "github.com/ockibagusp/golang-website-example/app/main/controller/mock/models"
	selectUser "github.com/ockibagusp/golang-website-example/business/user"
	"github.com/ockibagusp/golang-website-example/config"
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

func SessionMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			conf := config.GetAPPConfig()

			var (
				session_gorilla *sessions.Session
				err             error
			)

			// Test: session_test = true
			if conf.SetSessionToFalse() {
				session_gorilla, err = modelsTest.GetAuthSession()
			} else {
				session_gorilla, err = session.Get("session", c)
			}

			if err != nil {
				return c.HTML(http.StatusForbidden, "no session")
			}

			path := c.Request().URL.Path
			role := session_gorilla.Values["role"]
			if role != "admin" && role != "user" {
				role = "anonymous"
			}

			// -> role = "anonymous"
			if strings.Contains(path, "/login") || strings.Contains(path, "/logout") {
				return next(c)
			}

			id := session_gorilla.Values["id"]
			if id == "" {
				id = -1
			}

			username := session_gorilla.Values["username"]
			if username != "" {
				username = "anonymous"
			}

			c.Set("id", id)
			c.Set("username", username)
			c.Set("role", role)

			return next(c)
		}
	}
}

// SetSession: set session from User
func SetSession(user *selectUser.User, c echo.Context) (session_gorilla *sessions.Session, err error) {
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

	session_gorilla.Values["id"] = user.ID
	session_gorilla.Values["username"] = user.Username
	session_gorilla.Values["role"] = user.Role

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

	session_gorilla.Values["id"] = -1
	session_gorilla.Values["username"] = "anonymous"
	session_gorilla.Values["role"] = "anonymous"
	session_gorilla.Save(c.Request(), c.Response())
	return
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
