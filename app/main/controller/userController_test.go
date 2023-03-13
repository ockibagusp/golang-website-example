package controller_test

import (
	"net/http"
	"regexp"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/ockibagusp/golang-website-example/app/main/controller/mock/method"
	modelsTest "github.com/ockibagusp/golang-website-example/app/main/controller/mock/models"
	"github.com/ockibagusp/golang-website-example/app/main/types"
	"github.com/stretchr/testify/assert"
)

func TestUsersController(t *testing.T) {
	assert := assert.New(t)

	no_auth := setupTestServer(t)

	// test for SetSession = false
	method.SetSession = false
	// test for db users
	truncateUsers()

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

func TestCreateUserController(t *testing.T) {
	no_auth := setupTestServer(t)

	// test for SetSession = false
	method.SetSession = false
	// test for db users
	truncateUsers()

	// TODO: flash with redirect on failure, insyaallah

	test_cases := []struct {
		name   string
		method int    // method: 1=GET or 2=POST
		expect string // auth or no-auth
		form   types.UserForm
		status int

		// body navbar
		html_navbar regex
		// body heading
		html_heading regex
		// flash message
		html_flash_success regex
		html_flash_error   regex
	}{
		/*
			create it [admin]
		*/
		// GET
		{
			name:   "users [admin] to GET create it success",
			expect: ADMIN,
			method: method.HTTP_REQUEST_GET,
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
				actual:       `<h2 class="mt-4">New User</h2>`,
			},
		},
		// POST
		{
			name:   "user [admin] to POST create it success",
			expect: ADMIN,
			method: method.HTTP_REQUEST_POST,
			form: types.UserForm{
				Role:            "user",
				Username:        "unit-test",
				Email:           "unit-test@exemple.com",
				Name:            "Unit Test",
				Password:        "unit-test",
				ConfirmPassword: "unit-test",
			},
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
			// flash message success
			html_flash_success: regex{
				must_compile: `<strong>success:</strong> (.*)`,
				actual:       `<strong>success:</strong> success new user: unit-test!`,
			},
		},
		// Database: " Error 1062: Duplicate entry 'unit-test@exemple.com' for key 'email_UNIQUE' " v
		//			-> " Error 1062: Duplicate entry 'unit-test@exemple.com' for key 'users.email_UNIQUE' " x
		{
			name:   "users [admin] to POST create it failure: Duplicate entry",
			expect: ADMIN,
			method: method.HTTP_REQUEST_POST,
			form: types.UserForm{
				Role:            "user",
				Username:        "unit-test",
				Email:           "unit-test@exemple.com",
				Name:            "Unit Test",
				Password:        "unit-test",
				ConfirmPassword: "unit-test",
			},
			// HTTP response status: 400 Bad Request
			status: http.StatusBadRequest,
			// body navbar
			html_navbar: regex{
				must_compile: `<a class="btn">(.*)</a>`,
				actual:       `<a class="btn">ADMIN</a>`,
			},
			// body heading
			html_heading: regex{
				must_compile: `<h2 class="mt-4">(.*)</h2>`,
				actual:       `<h2 class="mt-4">New User</h2>`,
			},
			// flash message error
			html_flash_error: regex{
				must_compile: `<strong>error:</strong> (.*)`,
				actual:       `<strong>error:</strong> Error 1062 (23000): Duplicate entry &#39;unit-test@exemple.com&#39; for key &#39;email_UNIQUE&#39;!`,
			},
		},

		/*
			create it [sugriwa]
		*/
		// GET
		{
			name:   "users [sugriwa] to GET create it failure",
			expect: SUGRIWA,
			method: method.HTTP_REQUEST_GET,
			// redirect @route: /users/read/:id
			// HTTP response status: 200 OK
			status: http.StatusOK,
			// body heading
			html_heading: regex{
				must_compile: `<h2 class="mt-4">(.*)</h2>`,
				actual:       `<h2 class="mt-4">User: Sugriwa</h2>`,
			},
		},
		// POST
		{
			name:   "user [sugriwa] to POST create it failure",
			expect: SUGRIWA,
			method: method.HTTP_REQUEST_POST,
			// redirect @route: /users/read/:id
			// HTTP response status: 200 OK
			status: http.StatusOK,
			// body heading
			html_heading: regex{
				must_compile: `<h2 class="mt-4">(.*)</h2>`,
				actual:       `<h2 class="mt-4">User: Sugriwa</h2>`,
			},
			// flash message error
			html_flash_error: regex{
				must_compile: `<strong>error:</strong> (.*)`,
				actual:       `<strong>error:</strong> 403 Forbidden!`,
			},
		},

		/*
			create it [no-auth]
		*/
		// GET
		{
			name:   "users [no-auth] to GET create it success",
			expect: "anonymous",
			method: method.HTTP_REQUEST_GET,
			// HTTP response status: 200 OK
			status: http.StatusOK,
			// body navbar
			html_navbar: regex{
				must_compile: `<a href="/login" (.*)>Login</a>`,
				actual:       `<a href="/login" class="btn btn-outline-success my-2 my-sm-0">Login</a>`,
			},
			// body heading
			html_heading: regex{
				must_compile: `<h2 class="mt-4">(.*)</h2>`,
				actual:       `<h2 class="mt-4">New User</h2>`,
			},
		},
		// POST
		{
			name:   "user [no-auth] to POST create it success",
			expect: "anonymous",
			method: method.HTTP_REQUEST_POST,
			form: types.UserForm{
				Role:            "user",
				Username:        "example",
				Email:           "example@example.com",
				Name:            "Example",
				Password:        "example123",
				ConfirmPassword: "example123",
			},
			// HTTP response status: 200 OK
			status: http.StatusOK,
			// flash message success
			html_flash_success: regex{
				must_compile: `<strong>success:</strong> (.*)`,
				actual:       `<strong>success:</strong> success new user: example!`,
			},
			// TODO: difficult html_navbar and html_heading, insyaallah
		},
	}

	for _, test := range test_cases {
		modelsTest.UserSelectTest = test.expect

		t.Run(test.name, func(t *testing.T) {
			var result *httpexpect.Response
			if test.method == method.HTTP_REQUEST_GET {
				result = no_auth.GET("/users/add").
					WithForm(test.form).
					Expect().
					Status(test.status)
			} else if test.method == method.HTTP_REQUEST_POST {
				result = no_auth.POST("/users/add").
					WithForm(test.form).
					Expect().
					Status(test.status)
			} else {
				panic("method: 1=GET or 2=POST")
			}

			result_body := result.Body().Raw()

			var (
				must_compile, actual, match string
				regex                       *regexp.Regexp
			)

			if test.html_flash_success.must_compile != "" {
				must_compile = test.html_flash_success.must_compile
				actual = test.html_flash_success.actual

				regex = regexp.MustCompile(must_compile)
				match = regex.FindString(result_body)

				assert.Equal(t, match, actual)
			}

			if test.html_flash_error.must_compile != "" {
				must_compile = test.html_flash_error.must_compile
				actual = test.html_flash_error.actual

				regex = regexp.MustCompile(must_compile)
				match = regex.FindString(result_body)

				assert.Equal(t, match, actual)
			}

			assert := assert.New(t)
			if test.html_navbar.must_compile != "" {
				must_compile = test.html_navbar.must_compile
				actual = test.html_navbar.actual

				regex = regexp.MustCompile(must_compile)
				match = regex.FindString(result_body)

				assert.Equal(match, actual)
			}

			if test.html_heading.must_compile != "" {
				must_compile = test.html_heading.must_compile
				actual = test.html_heading.actual

				regex := regexp.MustCompile(must_compile)
				match := regex.FindString(result_body)

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
