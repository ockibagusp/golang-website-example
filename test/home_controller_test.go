package test

import (
	"net/http"
	"regexp"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/assert"
)

func TestHomeController(t *testing.T) {
	assert := assert.New(t)

	no_auth := setupTestServer(t)
	auth_admin := setupTestServerAuth(no_auth, 1)
	auth_sugriwa := setupTestServerAuth(no_auth, 2)

	test_cases := []struct {
		name      string
		expect    *httpexpect.Expect // auth_admin, session_sugriwa or no-auth
		navbar    regex
		jumbotron regex
	}{
		{
			name:   "home [no-auth] success",
			expect: no_auth,
			navbar: regex{
				must_compile: `<a href="/login" (.*)</a>`,
				actual:       `<a href="/login" class="btn btn-outline-success my-2 my-sm-0">Login</a>`,
			},
			jumbotron: regex{
				must_compile: `<p class="lead">(.*)</p>`,
				actual:       `<p class="lead">Test.</p>`,
			},
		},
		{
			name:   "home [admin] success",
			expect: auth_admin,
			navbar: regex{
				must_compile: `<a class="btn">(.*)</a>`,
				actual:       `<a class="btn">ADMIN</a>`,
			},
			jumbotron: regex{
				must_compile: `<p class="lead">(.*)</p>`,
				actual:       `<p class="lead">Admin.</p>`,
			},
		},
		{
			name:   "home [user] success",
			expect: auth_sugriwa,
			navbar: regex{
				must_compile: `<a href="/users" (.*)</a>`,
				actual:       `<a href="/users" class="btn btn-outline-secondary my-2 my-sm-0">Users</a>`,
			},
			jumbotron: regex{
				must_compile: `<p class="lead">(.*)</p>`,
				actual:       `<p class="lead">User.</p>`,
			},
		},
	}

	for _, test := range test_cases {
		var result *httpexpect.Response
		expect := test.expect // auth_admin, auth_sugriwa or no-auth

		t.Run(test.name, func(t *testing.T) {
			result = expect.GET("/").
				Expect().
				Status(http.StatusOK)

			result_body := result.Body().Raw()

			// TODO: why?
			var regex *regexp.Regexp
			var match string
			if test.navbar.must_compile != "" {
				// navbar nav
				regex = regexp.MustCompile(test.navbar.must_compile)
				match = regex.FindString(result_body)

				assert.Equal(match, test.navbar.actual)
			}

			if test.jumbotron.must_compile != "" {
				// main: jumbotron
				regex = regexp.MustCompile(test.jumbotron.must_compile)
				match = regex.FindString(result_body)

				assert.Equal(match, test.jumbotron.actual)
			}
		})
	}
}
