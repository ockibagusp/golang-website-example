package tests

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gavv/httpexpect/v2"
)

const (
	ADMIN      string = "admin"
	SUGRIWA           = "sugriwa"
	SUBALI            = "subali"
	OCKIBAGUSP        = "ockibagusp"
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
