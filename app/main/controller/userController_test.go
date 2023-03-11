package controller_test

import (
	"net/http"
	"regexp"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/ockibagusp/golang-website-example/app/main/controller/mock/method"
	modelsTest "github.com/ockibagusp/golang-website-example/app/main/controller/mock/models"
	"github.com/stretchr/testify/assert"
)

func TestUsersController(t *testing.T) {
	assert := assert.New(t)

	no_auth := setupTestServer(t)

	// test for SetSession = false
	method.SetSession = false
	// test for db users
	truncateUsers()

	// TODO: rows all, admin dan user, insyaallah

	test_cases := []struct {
		name         string
		expect       string // expect: admin, sugriwa
		url_query    string // @route: exemple /users?admin=all
		status       int
		html_navbar  regex
		html_heading regex
		html_table   regex
	}{
		/*
			users [admin]
		*/
		{
			name:   "users [admin] to GET it success: all",
			expect: ADMIN,
			// HTTP response status: 200 OK
			status: http.StatusOK,
			// body navbar
			html_navbar: regex{
				must_compile: `<a class="btn">(.*)</a>`,
				actual:       `<a class="btn">ADMIN</a>`,
			},
			// body heading
			html_heading: regex{
				must_compile: `<h2 class="mt-4">(.*)</h2>`,
				actual:       `<h2 class="mt-4">Users: All</h2>`,
			},
			html_table: regex{
				/*
					<tr>
						...
					    <td>
					        admin
					    </td>
						....
					</tr>
				*/
			},
		},
		{
			name:      "users [admin] to GET it success: admin",
			expect:    ADMIN,
			url_query: "admin",
			// HTTP response status: 200 OK
			status: http.StatusOK,
			// body navbar
			html_navbar: regex{
				must_compile: `<a class="btn">(.*)</a>`,
				actual:       `<a class="btn">ADMIN</a>`,
			},
			// body heading
			html_heading: regex{
				must_compile: `<h2 class="mt-4">(.*)</h2>`,
				actual:       `<h2 class="mt-4">Users: Admin</h2>`,
			},
		},
		{
			name:      "users [admin] to GET it success: user",
			expect:    ADMIN,
			url_query: "user",
			// HTTP response status: 200 OK
			status: http.StatusOK,
			// body navbar
			html_navbar: regex{
				must_compile: `<a class="btn">(.*)</a>`,
				actual:       `<a class="btn">ADMIN</a>`,
			},
			// body heading
			html_heading: regex{
				must_compile: `<h2 class="mt-4">(.*)</h2>`,
				actual:       `<h2 class="mt-4">Users: User</h2>`,
			},
		},
		{
			name:      "users [admin] to GET it failed: all",
			expect:    ADMIN,
			url_query: "false",
			status:    http.StatusOK,
			// body heading
			html_heading: regex{
				must_compile: `<h2 class="mt-4">(.*)</h2>`,
				actual:       `<h2 class="mt-4">Users: All</h2>`,
			},
		},

		/*
			users [user]
		*/
		{
			name:   "users [user] to GET it redirect success: sugriwa",
			expect: SUGRIWA,
			// redirect @route: /user/read/2 [sugriwa: 2]
			// HTTP response status: 200 OK
			status: http.StatusOK,
			// body heading
			html_heading: regex{
				must_compile: `<h2 class="mt-4">(.*)</h2>`,
				actual:       `<h2 class="mt-4">User: Sugriwa</h2>`,
			},
		},
		{
			name:      "users [user] to GET it redirect success: admin failed",
			expect:    SUGRIWA,
			url_query: "admin",
			// redirect @route: /user/read/2 [sugriwa: 2]
			// HTTP response status: 200 OK
			status: http.StatusOK,
			// body heading
			html_heading: regex{
				must_compile: `<h2 class="mt-4">(.*)</h2>`,
				actual:       `<h2 class="mt-4">User: Sugriwa</h2>`,
			},
		},

		/*
			No Auth
		*/
		{
			name:   "users [no-auth] to GET it failure: login",
			expect: ANONYMOUS,
			// HTTP response status: 200 OK
			status: http.StatusOK,
			// body navbar
			html_navbar: regex{
				must_compile: `<p class="text-danger">*(.*)!</p>`,
				actual:       `<p class="text-danger">*login process failed!</p>`,
			},
		},
	}

	for _, test := range test_cases {
		var result *httpexpect.Response

		t.Run(test.name, func(t *testing.T) {
			modelsTest.UserSelectTest = test.expect // admin, sugriwa

			// @route: exemple "/users?admin=all"
			if test.url_query != "" {
				result = no_auth.GET("/users").
					WithQuery(test.url_query, "all").
					Expect().
					Status(test.status)
			} else {
				// @route: "/users"
				result = no_auth.GET("/users").
					Expect().
					Status(test.status)
			}

			result_body := result.Body().Raw()

			var (
				must_compile, actual, match string
				regex                       *regexp.Regexp
			)

			if test.html_navbar.must_compile != "" {
				must_compile = test.html_navbar.must_compile
				actual = test.html_navbar.actual

				regex = regexp.MustCompile(must_compile)
				match = regex.FindString(result_body)

				// assert.Equal(t, match, actual)
				//
				// or,
				//
				// assert := assert.New(t)
				// ...
				// assert.Equal(match, actual)
				assert.Equal(match, actual)
			}

			if test.html_heading.must_compile != "" {
				must_compile = test.html_heading.must_compile
				actual = test.html_heading.actual

				regex = regexp.MustCompile(must_compile)
				match = regex.FindString(result_body)

				assert.Equal(match, actual)
			}

			statusCode := result.Raw().StatusCode
			if test.status != statusCode {
				t.Logf(
					"got: %d but expect %d", test.status, statusCode,
				)
				t.Fail()
			}
		})
	}
}
