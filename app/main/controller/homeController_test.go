package controller_test

import (
	"net/http"
	"regexp"
	"testing"

	"github.com/gavv/httpexpect/v2"
	methodTest "github.com/ockibagusp/golang-website-example/app/main/controller/mock/method"
	modelsTest "github.com/ockibagusp/golang-website-example/app/main/controller/mock/models"
	"github.com/stretchr/testify/assert"
)

func TestHomeController(t *testing.T) {
	assert := assert.New(t)

	no_auth := setupTestServer(t)

	// test for SetSession = false
	methodTest.SetSession = false
	// test for db users
	truncateUsers()

	test_cases := []struct {
		name           string
		expect         string // admin, sugriwa
		html_navbar    regex
		html_jumbotron regex
	}{
		{
			name:   "home [no-auth] success",
			expect: ANONYMOUS,
			html_navbar: regex{
				must_compile: `<a href="/login" (.*)</a>`,
				actual:       `<a href="/login" class="btn btn-outline-success my-2 my-sm-0">Login</a>`,
			},
			html_jumbotron: regex{
				must_compile: `<p class="lead">(.*)</p>`,
				actual:       `<p class="lead">Test.</p>`,
			},
		},
		{
			name:   "home [admin] success",
			expect: ADMIN,
			html_navbar: regex{
				must_compile: `<a class="btn">(.*)</a>`,
				actual:       `<a class="btn">ADMIN</a>`,
			},
			html_jumbotron: regex{
				must_compile: `<p class="lead">(.*)</p>`,
				actual:       `<p class="lead">Admin.</p>`,
			},
		},
		{
			name:   "home [user] success",
			expect: SUGRIWA,
			html_navbar: regex{
				must_compile: `<a href="/users" (.*)</a>`,
				actual:       `<a href="/users" class="btn btn-outline-secondary my-2 my-sm-0">Users</a>`,
			},
			html_jumbotron: regex{
				must_compile: `<p class="lead">(.*)</p>`,
				actual:       `<p class="lead">User.</p>`,
			},
		},
	}

	for _, test := range test_cases {
		var result *httpexpect.Response
		modelsTest.UserSelectTest = test.expect // auth_admin, auth_sugriwa or no-auth

		t.Run(test.name, func(t *testing.T) {
			result = no_auth.GET("/").
				Expect().
				Status(http.StatusOK)
			result_body := result.Body().Raw()

			// why?
			var regex *regexp.Regexp
			var match string
			if test.html_navbar.must_compile != "" {
				// navbar nav
				regex = regexp.MustCompile(test.html_navbar.must_compile)
				match = regex.FindString(result_body)

				assert.Equal(match, test.html_navbar.actual)
			}

			if test.html_jumbotron.must_compile != "" {
				// main: jumbotron
				regex = regexp.MustCompile(test.html_jumbotron.must_compile)
				match = regex.FindString(result_body)

				assert.Equal(match, test.html_jumbotron.actual)
			}
		})
	}
}