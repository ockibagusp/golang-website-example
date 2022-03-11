package test

import (
	"net/http"
	"regexp"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/ockibagusp/golang-website-example/models"
	"github.com/ockibagusp/golang-website-example/types"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// truncate Users
//
// parameter: db *gorm.DB or not available:
// func truncateUsers() {...}, just the same
func truncateUsers(db *gorm.DB) {
	db.Exec("TRUNCATE users")

	// database: just `users.username` varchar 15
	users := []models.User{
		{
			Username: "admin",
			Email:    "admin@website.com",
			Password: "$2a$10$XJAj65HZ2c.n1iium4qUEeGarW0PJsqVcedBh.PDGMXdjqfOdN1hW",
			Name:     "Admin",
			IsAdmin:  1,
		},
		{
			Username: "sugriwa",
			Email:    "sugriwa@wanara.com",
			Password: "$2a$10$bVVMuFHe/iaydX9yO2AttOPT8WyhMPe9F8nDflEqEyJbGRD5.guFu",
			Name:     "Sugriwa",
		},
		{
			Username: "subali",
			Email:    "subali@wanara.com",
			Password: "$2a$10$eO8wPLSfBU.8KLUh/T9kDeBm0vIRjiCvsmWe8ou5fZHJ3cYAUcg6y",
			Name:     "Subali",
		},
	}

	tx := db.Begin()
	// *gorm.DB
	if err := tx.Create(&users).Error; err != nil {
		tx.Rollback()
		panic(err.Error())
	}
	tx.Commit()
}

// TODO: types users error
// // type: users test cases
// type usersTestCases []struct {
// 	name   string
// 	expect *httpexpect.Expect // auth or no-auth
// 	method int                // method: 1=GET or 2=POST
// 	path   int                // id=int. Exemple, id=1
// 	form   struct{} ?
// 	status int
// }

func TestUsersController(t *testing.T) {
	assert := assert.New(t)

	no_auth := setupTestServer(t)
	auth_admin := setupTestServerAuth(no_auth, 1)
	auth_sugriwa := setupTestServerAuth(no_auth, 2)

	// test for db users
	truncateUsers(db)

	// TODO: rows all, admin dan user

	test_cases := []struct {
		name         string
		expect       *httpexpect.Expect // auth_admin, session_sugriwa or no-auth
		url_query    string             // @route: exemple /users?admin=all
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
			expect: auth_admin,
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
			expect:    auth_admin,
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
			expect:    auth_admin,
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
			expect:    auth_admin,
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
			expect: auth_sugriwa,
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
			expect:    auth_sugriwa,
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
			expect: no_auth,
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
		expect := test.expect // auth_admin, session_sugriwa or no-auth

		t.Run(test.name, func(t *testing.T) {

			// @route: exemple "/users?admin=all"
			if test.url_query != "" {
				result = expect.GET("/users").
					WithQuery(test.url_query, "all").
					Expect().
					Status(test.status)
			} else {
				// @route: "/users"
				result = expect.GET("/users").
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
	auth_admin := setupTestServerAuth(no_auth, 1)
	auth_sugriwa := setupTestServerAuth(no_auth, 2)

	// test for db users
	truncateUsers(db)

	// TODO: flash with redirect on failure

	test_cases := []struct {
		name   string
		expect *httpexpect.Expect // auth or no-auth
		method int                // method: 1=GET or 2=POST
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
			expect: auth_admin,
			method: GET,
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
			expect: auth_admin,
			method: POST,
			form: types.UserForm{
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
		// Database: " Error 1062: Duplicate entry 'unit-test@exemple.com' for key 'users.email_UNIQUE' "
		{
			name:   "users [admin] to POST create it failure: Duplicate entry",
			expect: auth_admin,
			method: POST,
			form: types.UserForm{
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
				actual:       `<strong>error:</strong> Error 1062: Duplicate entry &#39;unit-test@exemple.com&#39; for key &#39;users.email_UNIQUE&#39;!`,
			},
		},

		/*
			create it [sugriwa]
		*/
		// GET
		{
			name:   "users [sugriwa] to GET create it failure",
			expect: auth_sugriwa,
			method: GET,
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
			expect: auth_sugriwa,
			method: POST,
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
			expect: no_auth,
			method: GET,
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
			expect: no_auth,
			method: POST,
			form: types.UserForm{
				Username:        "ockibagusp",
				Email:           "ocki.bagus.p@gmail.com",
				Name:            "Ocki Bagus Pratama",
				Password:        "user123",
				ConfirmPassword: "user123",
			},
			// HTTP response status: 200 OK
			status: http.StatusOK,
			// flash message success
			html_flash_success: regex{
				must_compile: `<strong>success:</strong> (.*)`,
				actual:       `<strong>success:</strong> success new user: ockibagusp!`,
			},
			// TODO: difficult html_navbar and html_heading
		},
	}

	for _, test := range test_cases {
		expect := test.expect // auth or no-auth

		t.Run(test.name, func(t *testing.T) {
			var result *httpexpect.Response
			if test.method == GET {
				result = expect.GET("/users/add").
					WithForm(test.form).
					Expect().
					Status(test.status)
			} else if test.method == POST {
				result = expect.POST("/users/add").
					WithForm(test.form).
					WithFormField("X-CSRF-Token", csrf_token).
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

func TestReadUserController(t *testing.T) {
	assert := assert.New(t)

	no_auth := setupTestServer(t)
	auth_admin := setupTestServerAuth(no_auth, 1)
	auth_sugriwa := setupTestServerAuth(no_auth, 2)

	// test for db users
	truncateUsers(db)

	test_cases := []struct {
		name         string
		expect       *httpexpect.Expect // auth or no-auth
		method       int                // method: 1=GET or 2=POST
		path         string
		status       int
		html_navbar  regex
		html_heading regex
		flash_error  regex
	}{
		/*
			read it [admin]
		*/
		{
			name:   "users [admin] to GET read it success",
			expect: auth_admin,
			method: GET,
			path:   "1",
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
				actual:       `<h2 class="mt-4">User: Admin</h2>`,
			},
		},
		{
			name:   "users [admin] to GET read it failure",
			expect: auth_admin,
			method: GET,
			path:   "-1",
			// HTTP response status: 406 Not Acceptable
			status: http.StatusNotAcceptable,
		},

		/*
			read it [sugriwa]
		*/
		{
			name:   "users [sugriwa] to GET read it success",
			expect: auth_sugriwa,
			method: GET,
			path:   "2",
			// HTTP response status: 200 OK
			status: http.StatusOK,
			// body heading
			html_heading: regex{
				must_compile: `<h2 class="mt-4">(.*)</h2>`,
				actual:       `<h2 class="mt-4">User: Sugriwa</h2>`,
			},
		},
		{
			name:   "users [sugriwa] to GET read it failure",
			expect: auth_sugriwa,
			method: GET,
			path:   "-1",
			// HTTP response status: 406 Not Acceptable
			status: http.StatusNotAcceptable,
		},

		/*
			read it [no-auth]
		*/
		{
			name:   "users [no-auth] to GET read it failure",
			expect: no_auth,
			method: GET,
			path:   "1",
			// redirect @route: /login
			// HTTP response status: 200 OK
			status: http.StatusOK,
			// flash message
			flash_error: regex{
				must_compile: `<p class="text-danger">*(.*)</p>`,
				actual:       `<p class="text-danger">*login process failed!</p>`,
			},
		},
		{
			name:   "users [no-auth] to GET read it failure: 4 session and no-id",
			expect: no_auth,
			method: GET,
			path:   "-1",
			// redirect @route: /login
			// HTTP response status: 200 OK: 3 session and id
			status: http.StatusOK,
			// flash message
			flash_error: regex{
				must_compile: `<p class="text-danger">*(.*)</p>`,
				actual:       `<p class="text-danger">*login process failed!</p>`,
			},
		},
	}

	for _, test := range test_cases {
		var result *httpexpect.Response
		expect := test.expect // auth or no-auth

		t.Run(test.name, func(t *testing.T) {
			if test.method == GET {
				// same:
				//
				// expect.GET("/users/read/{id}").
				//	WithPath("id", test.path).
				// ...
				result = expect.GET("/users/read/{id}", test.path).
					Expect().
					Status(test.status)

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

					assert.Equal(match, actual)
				}

				if test.html_heading.must_compile != "" {
					must_compile = test.html_heading.must_compile
					actual = test.html_heading.actual

					regex = regexp.MustCompile(must_compile)
					match = regex.FindString(result_body)

					assert.Equal(match, actual)
				}

				if test.flash_error.must_compile != "" {
					must_compile = test.flash_error.must_compile
					actual = test.flash_error.actual

					regex = regexp.MustCompile(must_compile)
					match = regex.FindString(result_body)

					assert.Equal(match, actual)
				}
			} else {
				panic("method: 1=GET")
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

func TestUpdateUserController(t *testing.T) {
	no_auth := setupTestServer(t)
	auth_admin := setupTestServerAuth(no_auth, 1)
	auth_sugriwa := setupTestServerAuth(no_auth, 2)

	// test for db users
	truncateUsers(db)

	test_cases := []struct {
		name   string
		expect *httpexpect.Expect // auth or no-auth
		method int                // method: 1=GET or 2=POST
		path   string             // id=string. Exemple, id="1"
		form   types.UserForm
		status int

		html_navbar  regex
		html_heading regex
		// flash message
		html_flash_success regex
		html_flash_error   regex
	}{
		/*
			update it [admin]
		*/
		// GET
		{
			name:   "users [admin] to admin GET update it success: id=1",
			expect: auth_admin,
			method: GET,
			path:   "1", // admin: 1 admin
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
				actual:       `<h2 class="mt-4">User: Admin</h2>`,
			},
		},
		{
			name:   "users [admin] to user GET update it success: id=2",
			expect: auth_admin,
			method: GET,
			path:   "2", // user: 2 sugriwa
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
				actual:       `<h2 class="mt-4">User: Sugriwa</h2>`,
			},
		},
		{
			name:   "users [admin] to -1 GET update it failure: id=-1",
			expect: auth_admin,
			method: GET,
			path:   "-1",
			// HTTP response status: 404 Not Found
			status: http.StatusNotFound,
		},
		// POST
		{
			name:   "users [admin] to admin POST update it success: id=1",
			expect: auth_admin,
			method: POST,
			path:   "1", // admin: 1 admin
			form: types.UserForm{
				Username: "admin-success",
			},
			// redirect @route: /users
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
				actual:       `<strong>success:</strong> success update user: admin-success!`,
			},
		},
		{
			name:   "users [admin] to user POST update it success: id=2",
			expect: auth_admin,
			method: POST,
			path:   "2", // user: 2 sugriwa
			form: types.UserForm{
				// id=2 username: sugriwa
				Username: "sugriwa",
				Name:     "Sugriwa Success",
			},
			// redirect @route: /users
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
				// [admin] id=2 username: sugriwa
				actual: `<strong>success:</strong> success update user: sugriwa!`,
			},
		},
		{
			name:   "users [admin] to POST update it failure: id=-1",
			expect: auth_admin,
			method: POST,
			path:   "-1",
			form:   types.UserForm{},
			// HTTP response status: 404 Not Found
			status: http.StatusNotFound,
		},

		/*
			update it [sugriwa]
		*/
		// GET
		{
			name:   "users [sugriwa] to GET update it success: id=2",
			expect: auth_sugriwa,
			method: GET,
			path:   "2", // user: 2 sugriwa ok
			// HTTP response status: 200 OK
			status: http.StatusOK,
			// body heading
			html_heading: regex{
				must_compile: `<h2 class="mt-4">(.*)</h2>`,
				actual:       `<h2 class="mt-4">User: Sugriwa Success</h2>`,
			},
		},
		{
			name:   "users [sugriwa] to GET update it failure: id=-2",
			expect: auth_sugriwa,
			method: GET,
			path:   "-2",
			// HTTP response status: 404 Not Found
			status: http.StatusNotFound,
		},
		{
			name:   "users [sugriwa] to GET update it failure: id=3",
			expect: auth_sugriwa,
			method: GET,
			path:   "3", // user: 2 sugriwa no
			// HTTP response status: 403 Forbidden,
			status: http.StatusForbidden,
		},
		// POST
		{
			name:   "users [sugriwa] to sugriwa POST update it success",
			expect: auth_sugriwa,
			method: POST,
			path:   "2", // user: 2 sugriwa
			form: types.UserForm{
				Username: "sugriwa", // admin: "sugriwa-success" to sugriwa: "sugriwa"
				Name:     "Sugriwa",
			},
			// redirect @route: /
			// HTTP response status: 200 OK
			status: http.StatusOK,
			// body heading
			html_heading: regex{
				must_compile: `<h1 class="display-4">(.*)</h1>`,
				actual:       `<h1 class="display-4">Hello Sugriwa!</h1>`,
			},
			// flash message success
			html_flash_success: regex{
				must_compile: `<strong>success:</strong> (.*)`,
				actual:       `<strong>success:</strong> success update user: sugriwa!`,
			},
		},
		{
			name:   "users [sugriwa] to POST update it failure",
			expect: auth_sugriwa,
			method: POST,
			path:   "3", // user: 2 sugriwa no
			form: types.UserForm{
				Username: "subali-failure",
			},
			// HTTP response status: 403 Forbidden
			status: http.StatusForbidden,
		},

		/*
			update it [no-auth]
		*/
		// GET
		{
			name:   "users [no-auth] to GET update it failure: id=1",
			expect: no_auth,
			method: GET,
			path:   "1",
			// redirect @route: /login
			// HTTP response status: 200 OK
			status: http.StatusOK,
			// flash message
			html_flash_error: regex{
				must_compile: `<p class="text-danger">*(.*)</p>`,
				actual:       `<p class="text-danger">*login process failed!</p>`,
			},
		},
		{
			name:   "users [no-auth] to GET update it failure: id=-1",
			expect: no_auth,
			method: GET,
			path:   "-1",
			// redirect @route: /login
			// HTTP response status: 200 OK: 3 session and id
			status: http.StatusOK,
		},
		// POST
		{
			name:   "users [no-auth] to POST update it failure: id=2",
			expect: no_auth,
			method: POST,
			path:   "2",
			form: types.UserForm{
				Username: "sugriwa-failure",
			},
			// redirect @route: /login
			// HTTP response status: 200 OK: 3 session and id
			status: http.StatusOK,
		},
		{
			name:   "users [no-auth] to POST update it failure: id=-2",
			expect: no_auth,
			method: POST,
			path:   "-2",
			form: types.UserForm{
				Username: "sugriwa-failure",
			},
			// redirect @route: /login
			// HTTP response status: 200 OK: 3 session and id
			status: http.StatusOK,
		},
	}

	for _, test := range test_cases {
		expect := test.expect // auth or no-auth

		t.Run(test.name, func(t *testing.T) {
			var result *httpexpect.Response
			if test.method == GET {
				// same:
				//
				// expect.GET("/users/view/{id}").
				//	WithPath("id", test.path).
				// ...
				result = expect.GET("/users/view/{id}", test.path).
					WithForm(test.form).
					Expect().
					Status(test.status)
			} else if test.method == POST {
				result = expect.POST("/users/view/{id}").
					WithPath("id", test.path).
					WithForm(test.form).
					WithFormField("X-CSRF-Token", csrf_token).
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
				assert.Equal(t, match, actual)
			}

			if test.html_heading.must_compile != "" {
				must_compile = test.html_heading.must_compile
				actual = test.html_heading.actual

				regex = regexp.MustCompile(must_compile)
				match = regex.FindString(result_body)

				assert.Equal(t, match, actual)
			}

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

func TestUpdateUserByPasswordUserController(t *testing.T) {
	no_auth := setupTestServer(t)
	auth_admin := setupTestServerAuth(no_auth, 1)
	auth_sugriwa := setupTestServerAuth(no_auth, 2)

	// test for db users
	truncateUsers(db)

	test_cases := []struct {
		name   string
		expect *httpexpect.Expect // auth or no-auth
		method int                // method: 1=GET or 2=POST
		path   string             // id=string. Exemple, id="1"
		form   types.NewPasswordForm
		status int

		html_heading regex
		// flash message
		html_flash_success regex
		html_flash_error   regex
	}{
		/*
			update by password it [admin]
		*/
		// GET
		{
			name:   "users [admin] to GET update user by password it success: id=1",
			expect: auth_admin,
			method: GET,
			path:   "1",
			// HTTP response status: 200 OK
			status: http.StatusOK,
			// body heading
			html_heading: regex{
				must_compile: `<h3 class="mt-4">(.*)</h3>`,
				actual:       `<h3 class="mt-4">User: Admin</h3>`,
			},
		},
		{
			name:   "users [admin] to [sugriwa] GET update user by password it success: id=2",
			expect: auth_admin,
			method: GET,
			path:   "2",
			// HTTP response status: 200 OK
			status: http.StatusOK,
			// body heading
			html_heading: regex{
				must_compile: `<h3 class="mt-4">(.*)</h3>`,
				actual:       `<h3 class="mt-4">User: Sugriwa</h3>`,
			},
		},
		{
			name: "users [admin] to GET update user by password it failure: id=-1" +
				" GET passwords don't match",
			expect: auth_admin,
			method: GET,
			path:   "-1",
			// HTTP response status: 404 Not Found
			status: http.StatusNotFound,
		},
		// POST
		{
			name:   "users [admin] to POST update user by password it success: id=1",
			expect: auth_admin,
			method: POST,
			path:   "1",
			form: types.NewPasswordForm{
				OldPassword:        "admin123",
				NewPassword:        "admin_success",
				ConfirmNewPassword: "admin_success",
			},
			// HTTP response status: 200 OK
			status: http.StatusOK,
			// body heading
			html_heading: regex{
				must_compile: `<h2 class="mt-4">(.*)</h2>`,
				actual:       `<h2 class="mt-4">Users: All</h2>`,
			},
			// flash message success
			html_flash_success: regex{
				must_compile: `<strong>success:</strong> (.*)`,
				actual:       `<strong>success:</strong> success update user by password: admin!`,
			},
		},
		{
			name:   "users [admin] to [sugriwa] POST update user by password it success: id=2",
			expect: auth_admin,
			method: POST,
			path:   "2",
			form: types.NewPasswordForm{
				OldPassword:        "user123",
				NewPassword:        "user_success",
				ConfirmNewPassword: "user_success",
			},
			// HTTP response status: 200 OK
			status: http.StatusOK,
			// body heading
			html_heading: regex{
				must_compile: `<h2 class="mt-4">(.*)</h2>`,
				actual:       `<h2 class="mt-4">Users: All</h2>`,
			},
			// flash message success
			html_flash_success: regex{
				must_compile: `<strong>success:</strong> (.*)`,
				actual:       `<strong>success:</strong> success update user by password: sugriwa!`,
			},
		},
		{
			name: "users [auth] to POST update user by password it failure: id=1" +
				" POST passwords don't match",
			expect: auth_admin,
			method: POST,
			path:   "1",
			form: types.NewPasswordForm{
				OldPassword:        "admin_success",
				NewPassword:        "admin_success_",
				ConfirmNewPassword: "admin_failure",
			},
			// HTTP response status: 403 Forbidden
			status: http.StatusForbidden,
		},
		{
			name: "users [auth] to [sugriwa] POST update user by password it failure: id=2" +
				" POST passwords don't match",
			expect: auth_admin,
			method: POST,
			path:   "2",
			form: types.NewPasswordForm{
				OldPassword:        "admin_password_success",
				NewPassword:        "admin_password_failure",
				ConfirmNewPassword: "admin_password_success_",
			},
			// HTTP response status: 403 Forbidden
			status: http.StatusForbidden,
		},
		{
			name:   "users [auth] to POST update user by password it failure: id=-1",
			expect: auth_admin,
			method: POST,
			path:   "-1",
			form:   types.NewPasswordForm{},
			// HTTP response status: 404 Not Found
			status: http.StatusNotFound,
		},

		/*
			update by password it [sugriwa]
		*/
		// GET
		{
			name:   "users [sugriwa] to GET update user by password it success: id=2",
			expect: auth_sugriwa,
			method: GET,
			path:   "2",
			// HTTP response status: 200 OK
			status: http.StatusOK,
			// body heading
			html_heading: regex{
				must_compile: `<h3 class="mt-4">(.*)</h3>`,
				actual:       `<h3 class="mt-4">User: Sugriwa</h3>`,
			},
		},
		{
			name:   "users [sugriwa] to [admin] GET update user by password it failure: id=1",
			expect: auth_sugriwa,
			method: GET,
			path:   "1",
			// HTTP response status: 403 Forbidden
			status: http.StatusForbidden,
		},
		{
			name:   "users [sugriwa] to [subali] GET update user by password it failure: id=3",
			expect: auth_sugriwa,
			method: GET,
			path:   "3",
			// HTTP response status: 403 Forbidden
			status: http.StatusForbidden,
		},
		{
			name:   "users [sugriwa] to GET update user by password it failure: id=-1",
			expect: auth_sugriwa,
			method: GET,
			path:   "-1",
			// HTTP response status: 404 Not Found
			status: http.StatusNotFound,
		},
		// POST
		{
			name:   "users [sugriwa] to POST update user by password it success: id=2",
			expect: auth_sugriwa,
			method: POST,
			path:   "2",
			form: types.NewPasswordForm{
				OldPassword:        "user_success",
				NewPassword:        "user123",
				ConfirmNewPassword: "user123",
			},
			// HTTP response status: 200 OK
			status: http.StatusOK,
			// flash message success
			html_flash_success: regex{
				must_compile: `<strong>success:</strong> (.*)`,
				actual:       `<strong>success:</strong> success update user by password: sugriwa!`,
			},
		},
		{
			name:   "users [sugriwa] to [admin] POST update user by password it failure: id=1",
			expect: auth_sugriwa,
			method: POST,
			path:   "1",
			form:   types.NewPasswordForm{},
			// HTTP response status: 403 Forbidden,
			status: http.StatusForbidden,
		},
		{
			name:   "users [sugriwa] to [subali] POST update user by password it failure: id=3",
			expect: auth_sugriwa,
			method: POST,
			path:   "3",
			form:   types.NewPasswordForm{},
			// HTTP response status: 403 Forbidden,
			status: http.StatusForbidden,
		},
		{
			name:   "users [sugriwa] to POST update user by password it failure: id=-1",
			expect: auth_sugriwa,
			method: POST,
			path:   "-1",
			form:   types.NewPasswordForm{},
			// HTTP response status: 404 Not Found
			status: http.StatusNotFound,
		},

		/*
			update by password it [no-auth]
		*/
		// GET
		{
			name:   "users [no-auth] to GET update user by password it failure: id=1",
			expect: no_auth,
			method: GET,
			path:   "1",
			// redirect @route: /login
			// HTTP response status: 200 OK
			status: http.StatusOK,
			// flash message
			html_flash_error: regex{
				must_compile: `<p class="text-danger">*(.*)</p>`,
				actual:       `<p class="text-danger">*login process failed!</p>`,
			},
		},
		{
			name:   "users [no-auth] to POST update user by password it failure: id=-1",
			expect: no_auth,
			method: GET,
			path:   "-1",
			// redirect @route: /login
			// HTTP response status: 200 OK
			status: http.StatusOK,
			// flash message
			html_flash_error: regex{
				must_compile: `<p class="text-danger">*(.*)</p>`,
				actual:       `<p class="text-danger">*login process failed!</p>`,
			},
		},
		// POST
		{
			name:   "users [no-auth] to POST update user by password it failure: id=1",
			expect: no_auth,
			method: POST,
			path:   "1",
			// redirect @route: /login
			// HTTP response status: 200 OK
			status: http.StatusOK,
			form:   types.NewPasswordForm{},
			// flash message
			html_flash_error: regex{
				must_compile: `<p class="text-danger">*(.*)</p>`,
				actual:       `<p class="text-danger">*login process failed!</p>`,
			},
		},
		{
			name:   "users [no-auth] to POST update user by password it success: id=-1",
			expect: no_auth,
			method: POST,
			path:   "-1",
			// redirect @route: /login
			// HTTP response status: 200 OK
			status: http.StatusOK,
			form:   types.NewPasswordForm{},
			// flash message
			html_flash_error: regex{
				must_compile: `<p class="text-danger">*(.*)</p>`,
				actual:       `<p class="text-danger">*login process failed!</p>`,
			},
		},
	}

	// for...{...}, same:
	//
	// t.Run("users [auth] to POST update user by password it success", func(t *testing.T) {
	// 	auth.POST("/users/view/{id}/password").
	// 		WithPath("id", "1").
	// 		WithForm(types.NewPasswordForm{
	// 			...
	// 		}).
	// 		Expect().
	// 		Status(http.StatusOK)
	// })
	//
	// ...
	//
	// t.Run("users [no-auth] to POST update user by password it failure: 4"+
	// 	" no session", func(t *testing.T) {
	// 	noAuth.POST("/users/view/{id}/password").
	// 		WithPath("id", "1").
	// 		Expect().
	// 		// redirect @route: /login
	// 		// HTTP response status: 200 OK
	// 		Status(http.StatusOK)
	// })
	for _, test := range test_cases {
		var result *httpexpect.Response
		expect := test.expect // auth or no-auth

		t.Run(test.name, func(t *testing.T) {
			if test.method == GET {
				// same:
				//
				// expect.POST("/users/view/{id}/password").
				//	WithPath("id", test.path).
				// ...
				result = expect.GET("/users/view/{id}/password", test.path).
					WithForm(test.form).
					Expect().
					Status(test.status)
			} else if test.method == POST {
				result = expect.POST("/users/view/{id}/password").
					WithPath("id", test.path).
					WithForm(test.form).
					WithFormField("X-CSRF-Token", csrf_token).
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

			if test.html_heading.must_compile != "" {
				must_compile = test.html_heading.must_compile
				actual = test.html_heading.actual

				regex = regexp.MustCompile(must_compile)
				match = regex.FindString(result_body)

				assert.Equal(t, match, actual)
			}

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

func TestDeleteUserController(t *testing.T) {
	no_auth := setupTestServer(t)
	auth_admin := setupTestServerAuth(no_auth, 1)
	auth_subali := setupTestServerAuth(no_auth, 3)

	// test for db users
	truncateUsers(db)

	test_cases := []struct {
		name   string
		expect *httpexpect.Expect // auth or no-auth
		path   string             // id=string. Exemple, id="1"
		status int

		// flash message
		html_flash_success regex
	}{
		// GET all
		/*
			delete it [admin]
		*/
		{
			name:   "users [admin] to [admin] DELETE it failure: id=1",
			expect: auth_admin,
			path:   "1",
			// HTTP response status: 403 Forbidden,
			status: http.StatusForbidden,
		},
		{
			name:   "users [admin] to [sugriwa] DELETE it success: id=2",
			expect: auth_admin,
			path:   "2",
			// redirect @route: /users
			// HTTP response status: 200 OK
			status: http.StatusOK,
			// flash message success
			html_flash_success: regex{
				must_compile: `<strong>success:</strong> (.*)`,
				actual:       `<strong>success:</strong> success delete user: sugriwa!`,
			},
		},
		{
			name:   "users [admin] to [sugriwa] DELETE it failure: id=2 delete exists",
			expect: auth_admin,
			path:   "2",
			// HTTP response status: 404 Not Found
			status: http.StatusNotFound,
		},
		{
			name:   "users [admin] to DELETE it failure: 2 (id=-1)",
			expect: auth_admin,
			path:   "-1",
			// HTTP response status: 404 Not Found
			status: http.StatusNotFound,
		},

		/*
			delete it [subali]
		*/
		{
			name:   "users [subali] to [admin] DELETE it failure: id=1",
			expect: auth_subali,
			path:   "1",
			// HTTP response status: 403 Forbidden,
			status: http.StatusForbidden,
		},
		{
			name:   "users [subali] to DELETE it failure: id=-1",
			expect: auth_subali,
			path:   "-1",
			// HTTP response status: 404 Not Found
			status: http.StatusNotFound,
		},
		{
			name:   "users [subali] to [subali] DELETE it success: id=3",
			expect: auth_subali,
			path:   "3",
			// redirect @route: /
			// HTTP response status: 200 OK
			status: http.StatusOK,
			// flash message success
			html_flash_success: regex{
				must_compile: `<strong>success:</strong> (.*)`,
				actual:       `<strong>success:</strong> success delete user: subali!`,
			},
		},

		/*
			delete it [na-auth]
		*/
		{
			name:   "users [no-auth] to DELETE it failure: id=1",
			expect: no_auth,
			path:   "1",
			// redirect @route: /login
			// HTTP response status: 200 OK
			status: http.StatusOK,
		},
		{
			name:   "users [no-auth] to DELETE it failure: id=-1",
			expect: no_auth,
			path:   "-1",
			// redirect @route: /login
			// HTTP response status: 200 OK
			status: http.StatusOK,
		},
		{
			name:   "users [no-auth] to DELETE it failure: id=error",
			expect: no_auth,
			path:   "error",
			// redirect @route: /login
			// HTTP response status: 200 OK
			status: http.StatusOK,
		},
	}

	for _, test := range test_cases {
		var result *httpexpect.Response
		expect := test.expect // auth or no-auth

		t.Run(test.name, func(t *testing.T) {
			result = expect.GET("/users/delete/{id}", test.path).
				Expect().
				Status(test.status)

			if test.html_flash_success.must_compile != "" {
				regex := regexp.MustCompile(test.html_flash_success.must_compile)
				match := regex.FindString(result.Body().Raw())

				assert.Equal(t, match, test.html_flash_success.actual)
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
