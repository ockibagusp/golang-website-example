package tests

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
func setupTestServer(t *testing.T, debug ...bool) (no_auth *httpexpect.Expect) {
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

	setupTestSetCookieCSRF(no_auth)

	return
}

// Setup test server to set cookie
func setupTestSetCookie(noAuth *httpexpect.Expect) {
	// TODO: set cookie to user and CSRF
}

// Setup test server to set cookie CSRF-Token
func setupTestSetCookieCSRF(noAuth *httpexpect.Expect) {
	set_cookie := noAuth.GET("/login").
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
	match := regex.FindStringSubmatch(set_cookie)

	csrf_token = match[1]
	// var expires string
	// csrf_token, expires = match[1], match[2]
	// csrf_token_expires, _ = time.Parse(time.RFC1123, expires)

}

// Setup test server no authentication and CSRF-Token
// request with cookie: csrf
func setupTestServerNoAuthCSRF(e *httpexpect.Expect) (no_auth_CSRF *httpexpect.Expect) {
	no_auth_CSRF = e.Builder(func(request *httpexpect.Request) {
		request.WithCookie("_csrf", csrf_token)
	})
	return
}

// Setup test server authentication
// request with cookie session and csrf
//
// @type is_user: 1 admin, 2 sugriwa and 3 subali.
func setupTestServerAuth(e *httpexpect.Expect, is_user int) (auth *httpexpect.Expect) {
	auth = e.Builder(func(request *httpexpect.Request) {
		var session string
		if is_user == 1 {
			session = session_admin
		} else if is_user == 2 {
			session = session_sugriwa
		} else if is_user == 3 {
			session = session_subali
		} else {
			panic("func setupTestServerAuth is type is_user: 1=admin, 2=sugriwa or 3=subali")
		}

		request.WithCookies(map[string]string{
			"_csrf":   csrf_token,
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
// session_admin: 13 Feb 2022
// session_admin: 17 Mar 2022
// username: admin
const session_admin = "MTY0NzUyOTA4OHxEdi1CQkFFQ180SUFBUkFCRUFBQVJ2LUNBQUlHYzNSeWFXNW5EQW9BQ0hWelpYSnVZVzFsQm5OMGNtbHVad3dIQUFWaFpHMXBiZ1p6ZEhKcGJtY01EZ0FNYVhOZllYVjBhRjkwZVhCbEEybHVkQVFDQUFJPXzLA59NKEOjtyt3H_fa63VIUy6Zd-KlT2Ag7iV1WjFCoQ=="

// session_sugriwa: 13 Feb 2022
// session_sugriwa: 17 Mar 2022
// username: sugriwa
const session_sugriwa = "MTY0NzUyOTE2MnxEdi1CQkFFQ180SUFBUkFCRUFBQVNQLUNBQUlHYzNSeWFXNW5EQW9BQ0hWelpYSnVZVzFsQm5OMGNtbHVad3dKQUFkemRXZHlhWGRoQm5OMGNtbHVad3dPQUF4cGMxOWhkWFJvWDNSNWNHVURhVzUwQkFJQUJBPT18RJKPkkUj4w1Xcz-nHAbVGbxjiNXkOJ9U-g5yTR4QOZM="

// session_subali: 13 Feb 2022
// session_subali: 17 Mar 2022
// username: subali
const session_subali = "MTY0NzUyOTI1OXxEdi1CQkFFQ180SUFBUkFCRUFBQVJfLUNBQUlHYzNSeWFXNW5EQW9BQ0hWelpYSnVZVzFsQm5OMGNtbHVad3dJQUFaemRXSmhiR2tHYzNSeWFXNW5EQTRBREdselgyRjFkR2hmZEhsd1pRTnBiblFFQWdBRXz8m__UUdCh3j937h2M3oGOaT91qMofTJPWbVoodkURPg=="

/*
	Cross Site Request Forgery (CSRF)

	HTTP Req. Headers:
	Set-Cookie: _csrf=M5CtIigue53Mcesal2vhW26OOfeOdGTq; Expires=Wed, 05 Jan 2022 10:47:03 GMT
*/
// 					 ________________________________
// Set-Cookie: _csrf=M5CtIigue53Mcesal2vhW26OOfeOdGTq; ...
//				  	 --------------------------------
var csrf_token string

// (?)
// 		   					_____________________________
// Set-Cookie: ...; Expires=Wed, 05 Jan 2022 10:47:03 GMT
// 		   					-----------------------------
// var csrfTokenExpires time.Time

func TestServer(t *testing.T) {
	//
}

func TestMain(m *testing.M) {
	exit := m.Run()
	// why?
	os.Exit(exit)
}
