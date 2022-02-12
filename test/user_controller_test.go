package test

import (
	"fmt"
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
	auth_user := setupTestServerAuth(no_auth, 0)

	testCases := []struct {
		name   string
		expect *httpexpect.Expect // auth_admin, auth_user or no-auth
		navbar regex
	}{
		{
			name:   "users [admin] to GET it success",
			expect: auth_admin,
			navbar: regex{
				must_compile: `<a class="btn">(.*)</a>`,
				actual:       `<a class="btn">ADMIN</a>`,
			},
		},
		{
			name:   "users [user] to GET it success",
			expect: auth_user,
			navbar: regex{
				must_compile: `<a href="/users" (.*)</a>`,
				actual:       `<a href="/users" class="btn btn-outline-secondary my-2 my-sm-0">Users</a>`,
			},
		},
		{
			name:   "users [no-auth] to GET it failure: login",
			expect: no_auth,
			navbar: regex{
				must_compile: `<p class="text-danger">*(.*)!</p>`,
				actual:       `<p class="text-danger">*login process failed!</p>`,
			},
		},
	}

	for _, test := range testCases {
		var result *httpexpect.Response
		expect := test.expect // auth_admin, auth_user or no-auth

		t.Run(test.name, func(t *testing.T) {
			result = expect.GET("/users").
				Expect().
				Status(http.StatusOK)

			result_body := result.Body().Raw()

			// navbar nav
			regex := regexp.MustCompile(test.navbar.must_compile)
			match := regex.FindString(result_body)

			assert.Equal(match, test.navbar.actual)
		})
	}
}

func TestCreateUserController(t *testing.T) {
	no_auth := setupTestServer(t)
	auth_admin := setupTestServerAuth(no_auth, 1)
	auth_user := setupTestServerAuth(no_auth, 0)

	// test for db users
	truncateUsers(db)

	// database: just `users.username` varchar 15
	user_form_sugriwa := types.UserForm{
		Username:        "sugriwa",
		Email:           "sugriwa@wanara.com",
		Name:            "Sugriwa",
		Password:        "user123",
		ConfirmPassword: "user123",
	}

	user_form_subali := user_form_sugriwa
	user_form_subali.Username = "subali"
	user_form_subali.Email = "subali@wanara.com"
	user_form_subali.Name = "subali"

	testCases := []struct {
		name   string
		expect *httpexpect.Expect // auth or no-auth
		method int                // method: 1=GET or 2=POST
		form   types.UserForm
		status int

		// flash message
		flash_success regex
		flash_error   regex
	}{
		/*
			GET create it success
		*/
		{
			name:   "users [admin] to GET create it success",
			expect: auth_admin,
			method: GET,
			// HTTP response status: 200 OK
			status: http.StatusOK,
		},
		{
			name:   "users [user] to GET create it success",
			expect: auth_user,
			method: GET,
			// HTTP response status: 200 OK
			status: http.StatusOK,
		},
		{
			name:   "users [no-auth] to GET create it success",
			expect: no_auth,
			method: GET,
			// HTTP response status: 200 OK
			status: http.StatusOK,
		},
		/*
			POST create it success
		*/
		{
			name:   "user [admin] to POST create it success",
			expect: auth_admin,
			method: POST,
			form:   user_form_sugriwa,
			// HTTP response status: 200 OK
			status: http.StatusOK,
			// flash message success
			flash_success: regex{
				must_compile: `<strong>success:</strong> (.*)`,
				actual:       `<strong>success:</strong> success new user: sugriwa!`,
			},
		},
		{
			name:   "user [user] to POST create it success",
			expect: auth_user,
			method: POST,
			form:   user_form_subali,
			// HTTP response status: 200 OK
			status: http.StatusOK,
			// flash message success
			flash_success: regex{
				must_compile: `<strong>success:</strong> (.*)`,
				actual:       `<strong>success:</strong> success new user: subali!`,
			},
		},

		// Database: " Error 1062: Duplicate entry 'sugriwa@wanara.com' for key 'users.email_UNIQUE' "
		{
			name:   "users [no-auth] to POST create it failure: Duplicate entry",
			expect: no_auth,
			method: POST,
			form:   user_form_sugriwa,
			// HTTP response status: 400 Bad Request
			status: http.StatusBadRequest,
			// flash message error
			flash_error: regex{
				must_compile: `<strong>error:</strong> (.*)`,
				actual:       `<strong>error:</strong> Error 1062: Duplicate entry &#39;sugriwa@wanara.com&#39; for key &#39;users.email_UNIQUE&#39;!.`,
			},
		},
	}

	for _, test := range testCases {
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
					WithFormField("X-CSRF-Token", csrfToken).
					Expect().
					Status(test.status)

				result_body := result.Body().Raw()

				var must_compile, actual string
				var match_actual bool
				if test.flash_success.must_compile != "" {
					match_actual = true
					must_compile = test.flash_success.must_compile
					actual = test.flash_success.actual
				}

				if test.flash_error.must_compile != "" {
					match_actual = true
					must_compile = test.flash_error.must_compile
					actual = test.flash_error.actual
				}

				if match_actual {
					regex := regexp.MustCompile(must_compile)
					match := regex.FindString(result_body)

					assert.Equal(t, match, actual)
				}
			} else {
				panic("method: 1=GET or 2=POST")
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
	no_auth := setupTestServer(t)
	auth_admin := setupTestServerAuth(no_auth, 1)
	auth_user := setupTestServerAuth(no_auth, 0)

	// test for db users
	truncateUsers(db)
	// database: just `users.username` varchar 15
	models.User{
		Username: "sugriwa",
		Email:    "sugriwa@wanara.com",
		Name:     "Sugriwa",
	}.Save(db)

	testCases := []struct {
		name        string
		expect      *httpexpect.Expect // auth or no-auth
		method      int                // method: 1=GET or 2=POST
		path        string
		status      int
		flash_error regex
	}{
		/*
			GET read it success
		*/
		{
			name:   "users [auth_admin] to GET read it success",
			expect: auth_admin,
			method: GET,
			path:   "1",
			// HTTP response status: 200 OK
			status: http.StatusOK,
		},
		{
			name:   "users [auth_user] to GET read it success",
			expect: auth_user,
			method: GET,
			path:   "1",
			// HTTP response status: 200 OK
			status: http.StatusOK,
		},
		/*
			GET read it failure
		*/
		{
			name:   "users [auth_admin] to GET read it failure: 1 session and no-id",
			expect: auth_admin,
			method: GET,
			path:   "-1",
			// HTTP response status: 406 Not Acceptable
			status: http.StatusNotAcceptable,
		},
		{
			name:   "users [auth_user] to GET read it failure: 2 session and no-id",
			expect: auth_user,
			method: GET,
			path:   "-1",
			// HTTP response status: 406 Not Acceptable
			status: http.StatusNotAcceptable,
		},
		{
			name:   "users [no_auth] to GET read it failure: 3 session and id",
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
			name:   "users [no_auth] to GET read it failure: 4 session and no-id",
			expect: no_auth,
			method: GET,
			path:   "-1",
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
			name:   "users [no_auth] to GET read it failure: 5 no-session and id",
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
			name:   "users [no_auth] to GET read it failure: 6 no-session and no-id",
			expect: no_auth,
			method: GET,
			path:   "-1",
			// redirect @route: /login
			// HTTP response status: 200 OK
			status: http.StatusOK,
			// flash message
			flash_error: regex{
				must_compile: `<p class="text-danger">*(.*)</p>`,
				actual:       `<p class="text-danger">*login process failed!</p>`,
			},
		},
	}

	for _, test := range testCases {
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

				var must_compile, actual string
				var match_actual bool

				if test.flash_error.must_compile != "" {
					match_actual = true
					must_compile = test.flash_error.must_compile
					actual = test.flash_error.actual
				}

				if match_actual {
					regex := regexp.MustCompile(must_compile)
					match := regex.FindString(result_body)

					assert.Equal(t, match, actual)
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
	noAuth := setupTestServer(t)
	auth := setupTestServerAuth(noAuth, 0)

	// test for db users
	truncateUsers(db)
	// database: just `users.username` varchar 15
	models.User{
		Username: "subali",
		Email:    "subali@wanara.com",
		Name:     "Subali",
	}.Save(db)

	testCases := []struct {
		name   string
		expect *httpexpect.Expect // auth or no-auth
		method int                // method: 1=GET or 2=POST
		path   string             // id=string. Exemple, id="1"
		form   types.UserForm
		status int

		// flash message
		isFlashSuccess     bool
		flashSuccessActual string
	}{
		{
			name:   "users [auth] to GET update it success",
			expect: auth,
			method: GET,
			path:   "1",
			// HTTP response status: 200 OK
			status: http.StatusOK,
		},
		{
			name:   "users [auth] to POST update it success",
			expect: auth,
			method: POST,
			path:   "1",
			form: types.UserForm{
				Username: "rahwana",
				Email:    "rahwana@rakshasa.com",
				Name:     "Rahwana",
			},
			// redirect @route: /users
			// HTTP response status: 200 OK
			status: http.StatusOK,
			// flash message success
			isFlashSuccess:     true,
			flashSuccessActual: "success update user: rahwana!",
		},
		{
			name:   "users [auth] to GET update it failure: 1 session and no-id",
			expect: auth,
			method: GET,
			path:   "-1",
			// HTTP response status: 406 Not Acceptable
			status: http.StatusNotAcceptable,
		},
		{
			name:   "users [no auth] to GET update it failure: 2 no-session and id",
			expect: noAuth,
			method: GET,
			path:   "1",
			// redirect @route: /login
			// HTTP response status: 200 OK
			status: http.StatusOK,
		},
		{
			name:   "users [no auth] to GET update it failure: 3 no-session and no-id",
			expect: noAuth,
			method: GET,
			path:   "-1",
			// redirect @route: /login
			// HTTP response status: 200 OK
			status: http.StatusOK,
		},
	}

	for _, test := range testCases {
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
					WithFormField("X-CSRF-Token", csrfToken).
					Expect().
					Status(test.status)

				if test.isFlashSuccess {
					flashSuccess := result.Body().Raw()

					regex := regexp.MustCompile(`<strong>success:</strong> (.*)`)
					match := regex.FindString(flashSuccess)

					actual := fmt.Sprintf("<strong>success:</strong> %s", test.flashSuccessActual)

					assert.Equal(t, match, actual)
				}
			} else {
				panic("method: 1=GET or 2=POST")
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
	noAuth := setupTestServer(t)
	auth := setupTestServerAuth(noAuth, 0)

	// test for db users
	truncateUsers(db)
	// database: just `users.username` varchar 15
	users := []models.User{
		{
			Username: "ockibagusp",
			Email:    "ocki.bagus.p@gmail.com",
			Password: "$2a$10$Y3UewQkjw808Ig90OPjuq.zFYIUGgFkWBuYiKzwLK8n3t9S8RYuYa",
			Name:     "Ocki Bagus Pratama",
		},
		{
			Username: "success",
			Email:    "success@exemple.com",
			Name:     "Success",
		},
	}
	// *gorm.DB
	db.Create(&users)

	testCases := []struct {
		name   string
		expect *httpexpect.Expect // auth or no-auth
		method int                // method: 1=GET or 2=POST
		path   string             // id=string. Exemple, id="1"
		form   types.NewPasswordForm
		status int

		// flash message
		isFlashSuccess     bool
		flashSuccessActual string
	}{
		{
			name:   "users [auth] to GET update user by password it success",
			expect: auth,
			method: GET,
			path:   "1",
			// HTTP response status: 200 OK
			status: http.StatusOK,
		},
		{
			name:   "users [auth] to POST update user by password it success",
			expect: auth,
			method: POST,
			path:   "1",
			form: types.NewPasswordForm{
				OldPassword:        "user123",
				NewPassword:        "password_success",
				ConfirmNewPassword: "password_success",
			},
			// HTTP response status: 200 OK
			status: http.StatusOK,
			// flash message success
			isFlashSuccess:     true,
			flashSuccessActual: "success update user by password: ockibagusp!",
		},
		{
			name: "users [auth] to GET update user by password it failure: 1" +
				" GET passwords don't match",
			expect: auth,
			method: GET,
			path:   "2",
			// HTTP response status: 406 Not Acceptabl
			status: http.StatusNotAcceptable,
		},
		{
			name: "users [auth] to POST update user by password it failure: 2" +
				" POST passwords don't match",
			expect: auth,
			method: POST,
			path:   "1",
			form: types.NewPasswordForm{
				OldPassword:        "user123",
				NewPassword:        "password_success",
				ConfirmNewPassword: "password_failure",
			},
			// HTTP response status: 403 Forbidden
			status: http.StatusForbidden,
		},
		{
			name: "users [auth] to POST update user by password it failure: 3" +
				" username don't match",
			expect: auth,
			method: POST,
			path:   "2",
			form: types.NewPasswordForm{
				OldPassword:        "user123",
				NewPassword:        "password_failure",
				ConfirmNewPassword: "password_failure",
			},
			// HTTP response status: 406 Not Acceptable
			status: http.StatusNotAcceptable,
		},
		{
			name: "users [no-auth] to GET update user by password it failure: 4" +
				" no session",
			expect: noAuth,
			method: GET,
			path:   "1",
			// redirect @route: /login
			// HTTP response status: 200 OK
			status: http.StatusOK,
		},
		{
			name: "users [no-auth] to POST update user by password it failure: 5" +
				" no session",
			expect: noAuth,
			method: POST,
			path:   "1",
			// redirect @route: /login
			// HTTP response status: 200 OK
			status: http.StatusOK,
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
	for _, test := range testCases {
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
					WithFormField("X-CSRF-Token", csrfToken).
					Expect().
					Status(test.status)

				if test.isFlashSuccess {
					flashSuccess := result.Body().Raw()

					regex := regexp.MustCompile(`<strong>success:</strong> (.*)`)
					match := regex.FindString(flashSuccess)

					actual := fmt.Sprintf("<strong>success:</strong> %s", test.flashSuccessActual)

					assert.Equal(t, match, actual)
				}
			} else {
				panic("method: 1=GET or 2=POST")
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
	noAuth := setupTestServer(t)
	auth := setupTestServerAuth(noAuth, 0)

	// test for db users
	truncateUsers(db)
	// database: just `users.username` varchar 15
	users := []models.User{
		{
			Username: "rahwana",
			Email:    "rahwana@rakshasa.com",
		},
	}
	// *gorm.DB
	db.Create(&users)

	testCases := []struct {
		name   string
		expect *httpexpect.Expect // auth or no-auth
		path   string             // id=string. Exemple, id="1"
		status int

		// flash message
		isFlashSuccess     bool
		flashSuccessActual string
	}{
		{
			name:   "users [auth] to DELETE it success",
			expect: auth,
			path:   "1",
			// redirect @route: /users
			// HTTP response status: 200 OK
			status: http.StatusOK,
			// flash message success
			isFlashSuccess:     true,
			flashSuccessActual: "success delete user: rahwana!",
		},
		{
			name:   "users [auth] to DELETE it failure: 1 (id=1) delete exists",
			expect: auth,
			path:   "1",
			// HTTP response status: 406 Not Acceptable
			status: http.StatusNotAcceptable,
		},
		{
			name:   "users [auth] to DELETE it failure: 2 (id=-1)",
			expect: auth,
			path:   "-1",
			// HTTP response status: 406 Not Acceptable
			status: http.StatusNotAcceptable,
		},
		{
			name:   "users [no-auth] to DELETE it failure: 3 (id=1)",
			expect: noAuth,
			path:   "1",
			// redirect @route: /login
			// HTTP response status: 200 OK
			status: http.StatusOK,
		},
		{
			name:   "users [no-auth] to DELETE it failure: 4 (id=-1)",
			expect: noAuth,
			path:   "-1",
			// redirect @route: /login
			// HTTP response status: 200 OK
			status: http.StatusOK,
		},
		{
			name:   "users [no-auth] to DELETE it failure: 5 (id=error)",
			expect: noAuth,
			path:   "error",
			// redirect @route: /login
			// HTTP response status: 200 OK
			status: http.StatusOK,
		},
	}

	for _, test := range testCases {
		var result *httpexpect.Response
		expect := test.expect // auth or no-auth

		t.Run(test.name, func(t *testing.T) {
			result = expect.GET("/users/delete/{id}", test.path).
				Expect().
				Status(test.status)

			if test.isFlashSuccess {
				flashSuccess := result.Body().Raw()

				regex := regexp.MustCompile(`<strong>success:</strong> (.*)`)
				match := regex.FindString(flashSuccess)

				actual := fmt.Sprintf("<strong>success:</strong> %s", test.flashSuccessActual)

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
