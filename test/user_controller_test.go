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
	noAuth := setupTestServer(t)
	auth := setupTestServerAuth(noAuth)

	t.Run("users [auth] to GET it success", func(t *testing.T) {
		auth.GET("/users").
			Expect().
			// HTTP response status: 200 OK
			Status(http.StatusOK)
	})

	t.Run("users [no auth] to GET it failure: login", func(t *testing.T) {
		flashError := noAuth.GET("/users").
			Expect().
			// redirect @route: /login
			// HTTP response status: 200 OK
			Status(http.StatusOK).
			Body().Raw()

		regex := regexp.MustCompile(`<p class\="text-danger">\*(.*)</p>`)
		match := regex.FindString(flashError)

		actual := `<p class="text-danger">*login!</p>`

		// flash message: "login!"
		assert.Equal(t, match, actual)
	})
}

func TestCreateUserController(t *testing.T) {
	noAuth := setupTestServer(t)
	auth := setupTestServerAuth(noAuth)

	// test for db users
	truncateUsers(db)

	// database: just `users.username` varchar 15
	userForm := types.UserForm{
		Username:        "sugriwa",
		Email:           "sugriwa@wanara.com",
		Name:            "Sugriwa",
		Password:        "user123",
		ConfirmPassword: "user123",
	}

	testCases := []struct {
		name   string
		expect *httpexpect.Expect // auth or no-auth
		method int                // method: 1=GET or 2=POST
		form   types.UserForm
		status int

		// flash message
		isFlashSuccess     bool
		flashSuccessActual string

		isFlashError     bool
		flashErrorActual string

		// or,
		//
		// flash struct {
		// 	success string
		// 	error   string
		// }
	}{
		{
			name:   "users [auth] to GET create it success",
			expect: auth,
			method: GET,
			// HTTP response status: 200 OK
			status: http.StatusOK,
		},
		{
			name:   "users [no auth] to GET create it success",
			expect: noAuth,
			method: GET,
			// HTTP response status: 200 OK
			status: http.StatusOK,
		},
		{
			name:   "user [auth] to POST create it success",
			expect: auth,
			method: POST,
			form:   userForm,
			// HTTP response status: 200 OK
			status: http.StatusOK,
			// flash message success
			isFlashSuccess:     true,
			flashSuccessActual: "success new user: sugriwa",
		},
		// Database: " Error 1062: Duplicate entry 'sugriwa@wanara.com' for key 'users.email_UNIQUE' "
		{
			name:   "users [no auth] to POST create it failure: Duplicate entry",
			expect: noAuth,
			method: POST,
			form:   userForm,
			// HTTP response status: 400 Bad Request
			status: http.StatusBadRequest,
			// flash message error
			isFlashError:     true,
			flashErrorActual: "Error 1062: Duplicate entry",

			// or,
			//
			// flash: struct {
			// 	success string
			// 	error   string
			// }{
			// 	error: "Error 1062: Duplicate entry",
			// },
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

				if test.isFlashSuccess {
					successMessage := result.Body().Raw()

					regex := regexp.MustCompile(`<strong>success:</strong> (.*)`)
					match := regex.FindString(successMessage)

					actual := fmt.Sprintf("<strong>success:</strong> %s!", test.flashSuccessActual)

					assert.Equal(t, match, actual)
				}

				if test.isFlashError {
					errorMessage := result.Body().Raw()

					actual := fmt.Sprintf("<strong>error:</strong> %s", test.flashErrorActual)
					regex := regexp.MustCompile(actual)
					match := regex.FindString(errorMessage)

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
	noAuth := setupTestServer(t)
	auth := setupTestServerAuth(noAuth)

	// test for db users
	truncateUsers(db)
	// database: just `users.username` varchar 15
	models.User{
		Username: "sugriwa",
		Email:    "sugriwa@wanara.com",
		Name:     "Sugriwa",
	}.Save(db)

	testCases := []struct {
		name   string
		expect *httpexpect.Expect // auth or no-auth
		method int                // method: 1=GET or 2=POST
		path   string
		status int
	}{
		{
			name:   "users [auth] to GET read it success",
			expect: auth,
			method: GET,
			path:   "1",
			// HTTP response status: 200 OK
			status: http.StatusOK,
		},
		{
			name:   "users [auth] to GET read it failure: 1 session and no-id",
			expect: auth,
			method: GET,
			path:   "-1",
			// HTTP response status: 406 Not Acceptable
			status: http.StatusNotAcceptable,
		},
		{
			name:   "users [no auth] to GET read it failure: 2 no-session and id",
			expect: noAuth,
			method: GET,
			path:   "1",
			// redirect @route: /login
			// HTTP response status: 200 OK
			status: http.StatusOK,
		},
		{
			name:   "users [no auth] to GET read it failure: 3 no-session and no-id",
			expect: noAuth,
			method: GET,
			path:   "-1",
			// redirect @route: /login
			// HTTP response status: 200 OK
			status: http.StatusOK,
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
	auth := setupTestServerAuth(noAuth)

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
					successMessage := result.Body().Raw()

					regex := regexp.MustCompile(`<strong>success:</strong> (.*)`)
					match := regex.FindString(successMessage)

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
	auth := setupTestServerAuth(noAuth)

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
					successMessage := result.Body().Raw()

					regex := regexp.MustCompile(`<strong>success:</strong> (.*)`)
					match := regex.FindString(successMessage)

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
	auth := setupTestServerAuth(noAuth)

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
				successMessage := result.Body().Raw()

				regex := regexp.MustCompile(`<strong>success:</strong> (.*)`)
				match := regex.FindString(successMessage)

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
