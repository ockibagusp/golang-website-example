package tests

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/ockibagusp/golang-website-example/models"
	"github.com/ockibagusp/golang-website-example/types"
)

/*
Setup test sever

TODO: .env debug: {true} or {false}, insyaallah

1. function debug (bool)
@function debug: {true} or {false}

2. os.Setenv("debug", ...)
@debug: {true} or {1}
os.Setenv("debug", "true") or,
os.Setenv("debug", "1")

@debug: {false} or {0}
os.Setenv("debug", "false") or,
os.Setenv("debug", "0")
*/
func setupTestServer(t *testing.T, debug ...bool) (no_auth *httpexpect.Expect) {
	os.Setenv("session_test", "1")
	os.Setenv("debug", "0")

	handler := setupTestHandler()

	server := httptest.NewServer(handler)
	defer server.Close()

	new_config := httpexpect.Config{
		BaseURL: server.URL,
		Client: &http.Client{
			Transport: httpexpect.NewBinder(handler),
			Jar:       httpexpect.NewJar(),
		},
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewCompactPrinter(t),
		},
	}

	if (len(debug) == 1 && debug[0] == true) || (os.Getenv("debug") == "1" || os.Getenv("debug") == "true") {
		new_config.Printers = []httpexpect.Printer{
			httpexpect.NewDebugPrinter(t, true),
		}
	} else if len(debug) > 1 {
		panic("func setupTestServer: (debug [1]: true or false) or no debug")
	}

	no_auth = httpexpect.WithConfig(new_config)
	setupTestSetCookie(no_auth)

	return
}

// Setup test server to set cookie
func setupTestSetCookie(no_auth *httpexpect.Expect) {
	// database: just `users.username` varchar 15
	users := []models.User{
		{
			Username: "admin",
			// Email:    "admin@website.com",
			Password: "$2a$10$XJAj65HZ2c.n1iium4qUEeGarW0PJsqVcedBh.PDGMXdjqfOdN1hW",
			// Name:     "Admin",
			// IsAdmin:  1,
		},
	}

	no_auth.POST("/login").
		WithForm(types.LoginForm{
			Username: users[0].Username,
			Password: users[0].Password,
		}).
		Expect().
		Status(http.StatusOK).
		Cookies().Raw()
}

// Setup test server authentication
// request with cookie session and csrf
//
// @type is_user: 1 admin, 2 sugriwa and 3 subali.
func setupTestServerAuth(e *httpexpect.Expect, is_user int) (auth *httpexpect.Expect) {
	auth = e.Builder(func(request *httpexpect.Request) {
		var session string
		if is_user == 1 {
			// session_admin: Expires=Fri, 06 Jan 2023 15:14:26 GMT
			// username: admin
			session = "MTY3MjQxMjk4M3xEdi1CQkFFQ180SUFBUkFCRUFBQVJ2LUNBQUlHYzNSeWFXNW5EQW9BQ0hWelpYSnVZVzFsQm5OMGNtbHVad3dIQUFWaFpHMXBiZ1p6ZEhKcGJtY01EZ0FNYVhOZllYVjBhRjkwZVhCbEEybHVkQVFDQUFJPXzGyQRp82nJXsgy9r8M36Dkx6TgQaU0DkXmG33b1oMemw=="
		} else if is_user == 2 {
			// session_sugriwa: Expires=Fri, 06 Jan 2023 15:14:26 GMT
			// username: sugriwa
			session = "MTY3MjQxMzE4OHxEdi1CQkFFQ180SUFBUkFCRUFBQVNQLUNBQUlHYzNSeWFXNW5EQW9BQ0hWelpYSnVZVzFsQm5OMGNtbHVad3dKQUFkemRXZHlhWGRoQm5OMGNtbHVad3dPQUF4cGMxOWhkWFJvWDNSNWNHVURhVzUwQkFJQUJBPT18AWZDueSxCE0bnsSgZ9JAhiZa-8BAH8_EGRa8wjDApoI="
		} else if is_user == 3 {
			// session_subali: Expires=Fri, 06 Jan 2023 15:14:26 GMT
			// username: subali
			session = "MTY3MjQxMzI2NnxEdi1CQkFFQ180SUFBUkFCRUFBQVJfLUNBQUlHYzNSeWFXNW5EQW9BQ0hWelpYSnVZVzFsQm5OMGNtbHVad3dJQUFaemRXSmhiR2tHYzNSeWFXNW5EQTRBREdselgyRjFkR2hmZEhsd1pRTnBiblFFQWdBRXzZ_DZu_InZaU3feck1AT0uGJnDETiBLWEe14y3OkiiWA=="
		} else {
			panic("func setupTestServerAuth is type is_user: 1=admin, 2=sugriwa or 3=subali")
		}

		request.WithCookies(map[string]string{
			"session": session,
		})
	})
	return
}

/*
	HTTP(s)
	----
	[+] Request Headers
	Cookie: session=...

	or,

	[+] Request Cookies
	session: ...

	-------
	Cookie:
	MaxAge=0 means no Max-Age attribute specified and the cookie will be
	deleted after the browser session ends.

	sessions.Options{.., MaxAge: 0,..}

	-------
	func. SetSession:

	session_gorilla, err = session.Get("session", ...)
	...
	session_gorilla.Values["username"] = user.Username
	session_gorilla.Values["is_auth_type"] = 2 // admin: 1 and user: 2
	---
	[+] Session:
	"username" = "ockibagusp"
	"is_auth_type" = 2
*/

func TestServer(t *testing.T) {}

func TestMain(m *testing.M) {
	exit := m.Run()
	// why?
	os.Exit(exit)
}
