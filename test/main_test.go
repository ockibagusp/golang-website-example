package test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"testing"

	"github.com/gavv/httpexpect/v2"
)

/*
	Setup test sever

	TODO: .env debug: {true} or {false}

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
func setupTestServer(t *testing.T, debug ...bool) (noAuth *httpexpect.Expect) {
	os.Setenv("debug", "0")

	handler := setupTestHandler()

	server := httptest.NewServer(handler)
	defer server.Close()

	newConfig := httpexpect.Config{
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
		newConfig.Printers = []httpexpect.Printer{
			httpexpect.NewDebugPrinter(t, true),
		}
	} else if len(debug) > 1 {
		panic("func setupTestServer: (debug [1]: true or false) or no debug")
	}

	noAuth = httpexpect.WithConfig(newConfig)

	setupTestSetCookieCSRF(noAuth)

	return
}

// Setup test server to set cookie
func setupTestSetCookie(noAuth *httpexpect.Expect) {
	// TODO: set cookie to user and CSRF
}

// Setup test server to set cookie CSRF-Token
func setupTestSetCookieCSRF(noAuth *httpexpect.Expect) {
	setCookie := noAuth.GET("/login").
		Expect().
		Status(http.StatusOK).
		Header("Set-Cookie").Raw()

	// Set-Cookie:
	// =================================== match[0] ================================
	// =	 																	   =
	// _csrf=M5CtIigue53Mcesal2vhW26OOfeOdGTq; Expires=Wed, 05 Jan 2022 10:47:03 GMT
	//		 --------------------------------	       -----------------------------
	//					match[1]							     match[2]
	regex := regexp.MustCompile(`_csrf\=(.*); Expires\=(.*)$`)
	match := regex.FindStringSubmatch(setCookie)

	csrfToken = match[1]
	// var expires string
	// csrfToken, expires = match[1], match[2]
	// csrfTokenExpires, _ = time.Parse(time.RFC1123, expires)

}

// Setup test server no authentication and CSRF-Token
// request with cookie: csrf
func setupTestServerNoAuthCSRF(e *httpexpect.Expect) (noAuthCSRF *httpexpect.Expect) {
	noAuthCSRF = e.Builder(func(request *httpexpect.Request) {
		request.WithCookie("_csrf", csrfToken)
	})
	return
}

// Setup test server authentication
// request with cookie session and csrf
//
// @type is_admin: 1 admin and 0 user.
func setupTestServerAuth(e *httpexpect.Expect, is_admin int) (auth *httpexpect.Expect) {
	auth = e.Builder(func(request *httpexpect.Request) {
		var session string
		if is_admin == 1 {
			session = session_admin
		} else if is_admin == 0 {
			session = session_user
		} else {
			panic("func setupTestServerAuth is type is_admin: 1=admin or 0=user")
		}

		request.WithCookies(map[string]string{
			"_csrf":   csrfToken,
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
// session_user: 23 Jan 2022
// username: ockibagusp
const session_user = "MTY0MjkzNDIwNnxEdi1CQkFFQ180SUFBUkFCRUFBQVNfLUNBQUlHYzNSeWFXNW5EQ" +
	"W9BQ0hWelpYSnVZVzFsQm5OMGNtbHVad3dNQUFwdlkydHBZbUZuZFhOd0JuTjBjbWx1Wnd3T0FBeHBjMTl" +
	"oZFhSb1gzUjVjR1VEYVc1MEJBSUFCQT09fDoaOeOnXeVm_zXJUWYidClXXXB3KevfkiI4v2O33QQ-"

// session_admin: 6 Feb 2022
// username: admin
const session_admin = "MTY0NDE0ODU3MHxEdi1CQkFFQ180SUFBUkFCRUFBQVJ2LUNBQUlHYzNSeWFXNW5E" +
	"QW9BQ0hWelpYSnVZVzFsQm5OMGNtbHVad3dIQUFWaFpHMXBiZ1p6ZEhKcGJtY01EZ0FNYVhOZllYVjBhRj" +
	"kwZVhCbEEybHVkQVFDQUFJPXxtxAIODyK4IVnBC8QT410I7adyvV1ziyqjm5jqsIoN0A=="

/*
	Cross Site Request Forgery (CSRF)

	HTTP Req. Headers:
	Set-Cookie: _csrf=M5CtIigue53Mcesal2vhW26OOfeOdGTq; Expires=Wed, 05 Jan 2022 10:47:03 GMT
*/
// 					 ________________________________
// Set-Cookie: _csrf=M5CtIigue53Mcesal2vhW26OOfeOdGTq; ...
//				  	 --------------------------------
var csrfToken string

// (?)
// 		   					_____________________________
// Set-Cookie: ...; Expires=Wed, 05 Jan 2022 10:47:03 GMT
// 		   					-----------------------------
// var csrfTokenExpires time.Time

func TestServer(t *testing.T) {
	//
}

func TestMain(m *testing.M) {
	// TODO: go test main_test.go ?
	// ----
	// cannot find package "." in:
	// /home/ockibagusp/go/src/github.com/ockibagusp/golang-website-example/vendor/main_test.go
	exit := m.Run()
	os.Exit(exit)
}
