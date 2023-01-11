package middleware

import (
	"errors"
	"os"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/ockibagusp/golang-website-example/models"
	modelsTest "github.com/ockibagusp/golang-website-example/tests/models"
)

// base.html -> {{if eq ((index .session.Values "is_auth_type") | tostring) -1 }}ok{{end}}

// GetAuth: get session to authenticated
func GetAuth(c echo.Context) (session_gorilla *sessions.Session, err error) {
	// Test: session_test = true
	if false { // os.Getenv("session_test") == "1" ???
		if modelsTest.UserSelectTest == "" && session_gorilla.IsNew == false {
			session_gorilla = &sessions.Session{
				Values: map[interface{}]interface{}{
					"username":     "",
					"is_auth_type": -1,
				},
			}

			err = errors.New("no session")
			return
		}

		for _, testUser := range modelsTest.UsersTest {
			if modelsTest.UserSelectTest == testUser.Username {
				session_gorilla = &sessions.Session{
					Values: map[interface{}]interface{}{
						"username": testUser.Username,
					},
				}

				if testUser.IsAdmin == 1 {
					session_gorilla.Values["is_auth_type"] = 1 // admin: 1
				} else if testUser.IsAdmin == 0 {
					session_gorilla.Values["is_auth_type"] = 2 // user: 2
				}
			}
		}
	} else {
		if session_gorilla, err = session.Get("session", c); err != nil {
			return
		}
	}

	is_auth_type := session_gorilla.Values["is_auth_type"]
	if IsAdmin(is_auth_type) || IsUser(is_auth_type) {
		return session_gorilla, nil
	}

	if _, ok := session_gorilla.Values["username"]; !ok {
		session_gorilla.Values["username"] = ""
	}
	if _, ok := session_gorilla.Values["is_auth_type"]; !ok {
		session_gorilla.Values["is_auth_type"] = -1
	}

	return
}

// IsAdmin: allows access only to authenticated administrators
func IsAdmin(is_auth_type interface{}) bool {
	return is_auth_type == 1
}

// IsUser: allows access only to authenticated users
func IsUser(is_auth_type interface{}) bool {
	return is_auth_type == 2
}

// GetAdmin: allows access only to authenticated administrators
func GetAdmin(c echo.Context) (session_gorilla *sessions.Session, err error) {
	if session_gorilla, err = session.Get("session", c); err != nil {
		return
	}

	is_auth_type := session_gorilla.Values["is_auth_type"]
	if IsAdmin(is_auth_type) {
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
	// Test: session_test = true
	if os.Getenv("session_test") == "1" {
		for _, testUser := range modelsTest.UsersTest {
			if modelsTest.UserSelectTest == testUser.Username {
				user.Username = testUser.Username
			}
		}
	}

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
	if user.IsAdmin == 1 {
		session_gorilla.Values["is_auth_type"] = 1 // admin: 1
	} else if user.IsAdmin == 0 {
		session_gorilla.Values["is_auth_type"] = 2 // user: 2
	}
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

const session_flash = "flash"

// cookieStoreFlash: new cookie store session for flash
func cookieStoreFlash() *sessions.CookieStore {
	return sessions.NewCookieStore(
		[]byte("secret-session-key"),
	)
}

// SetFlash: set session for flash message
func SetFlash(c echo.Context, name, value string) {
	tx_session_flash := session_flash
	if name == "message" {
		tx_session_flash += "-message"
	} else if name == "error" {
		tx_session_flash += "-error"
	}

	session, _ := cookieStoreFlash().Get(c.Request(), tx_session_flash)

	session.AddFlash(value, name)
	session.Save(c.Request(), c.Response())
}

// GetFlash: get session for flash messages
func GetFlash(c echo.Context, name string) (flashes []string) {
	tx_session_flash := session_flash
	if name == "message" {
		tx_session_flash += "-message"
	} else if name == "error" {
		tx_session_flash += "-error"
	}

	session, _ := cookieStoreFlash().Get(c.Request(), tx_session_flash)

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
