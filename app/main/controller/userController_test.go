package controller_test

import (
	"net/http"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/labstack/echo/v4"
	"github.com/ockibagusp/golang-website-example/app/main/types"
	"github.com/ockibagusp/golang-website-example/business/auth"
	"github.com/stretchr/testify/assert"
)

func TestCreateUsers_WithInputPOSTForSuccess(t *testing.T) {
	noAuth := setupTestServer(t)

	testCases := []struct {
		name   string
		form   types.UserForm
		token  echo.Map
		status int
	}{
		/*
			users admin [admin]
		*/
		{
			name: "user admin [admin] for POST create to success",
			token: echo.Map{
				"username": "admin",
				"role":     "admin",
			},
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
		},
		{
			name: "user no-auth [user] for POST create to success",
			token: echo.Map{
				"username": "anonymous",
				"role":     "anonymous",
			},
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
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			token, _ := auth.GenerateToken(conf.AppJWTAuthSign, 1, "admin", "admin")
			noAuth.POST("/users/add").
				WithCookie("token", token).
				WithForm(types.UserForm{
					Role:            test.form.Role,
					Username:        test.form.Username,
					Email:           test.form.Email,
					Name:            test.form.Name,
					Password:        test.form.Password,
					ConfirmPassword: test.form.ConfirmPassword,
				}).
				Expect().
				Status(test.status)
		})
	}

	// test for db users
	truncateUsers()
}

func TestCreateUsers_WithInputPOSTFormRoleWrongFailure(t *testing.T) {
	// assert
	assert := assert.New(t)

	noAuth := setupTestServer(t)

	var result *httpexpect.Response
	t.Run(`user anonymous for POST create wrong "role" failure`, func(t *testing.T) {
		// Error "role" wrong
		result = noAuth.POST("/users/add").
			WithForm(types.UserForm{
				Role:            "failure",
				Username:        "unit-test",
				Email:           "unit-test@exemple.com",
				Name:            "Unit Test",
				Password:        "unit-test",
				ConfirmPassword: "unit-test",
			}).
			Expect().
			Status(http.StatusBadRequest)

		assert.Contains(result.Body().Raw(), "Error 1265 (01000): Data truncated for column &#39;role&#39; at row 1!")
	})

	t.Run(`user admin for POST create wrong "role" failure`, func(t *testing.T) {
		token, _ := auth.GenerateToken(conf.AppJWTAuthSign, 1, "admin", "admin")

		//  Error "role" wrong
		result = noAuth.POST("/users/add").
			WithCookie("token", token).
			WithForm(types.UserForm{
				Role:            "failure",
				Username:        "unit-test",
				Email:           "unit-test@exemple.com",
				Name:            "Unit Test",
				Password:        "unit-test",
				ConfirmPassword: "unit-test",
			}).
			Expect().
			Status(http.StatusBadRequest)

		assert.Contains(result.Body().Raw(), "Error 1265 (01000): Data truncated for column &#39;role&#39; at row 1!")
	})

	// test for db users
	truncateUsers()
}

func TestCreateUsers_WithInputPOSTNotForFormFailure(t *testing.T) {
	// assert
	assert := assert.New(t)

	noAuth := setupTestServer(t)
	var result *httpexpect.Response

	testCases := []struct {
		name       string
		token      echo.Map
		form       types.UserForm
		status     int
		flashError string
	}{
		/*
			create form username failure
		*/
		{
			name: "user admin [admin] for POST form create to username too long failure: 2",
			token: echo.Map{
				"username": "admin",
				"role":     "admin",
			},
			form: types.UserForm{
				Role:            "user",
				Username:        "subali_copy_failure",
				Email:           "unit-test@exemple.com",
				Name:            "Unit Test",
				Password:        "unit-test",
				ConfirmPassword: "unit-test",
			},
			status:     http.StatusBadRequest,
			flashError: "<strong>error:</strong> username: the length must be between 4 and 15.!",
		},
		{
			name: "user anonymous for POST form create to username too long failure: 2",
			token: echo.Map{
				"username": "anonymous",
				"role":     "anonymous",
			},
			form: types.UserForm{
				Role:            "user",
				Username:        "anonymous_failure",
				Email:           "unit-test@exemple.com",
				Name:            "Unit Test",
				Password:        "unit-test",
				ConfirmPassword: "unit-test",
			},
			status:     http.StatusBadRequest,
			flashError: "<strong>error:</strong> username: the length must be between 4 and 15.!",
		},
		/*
			create form email failure
		*/
		{
			name: "user admin [admin] for POST form create to without email failure: 3",
			token: echo.Map{
				"username": "admin",
				"role":     "admin",
			},
			form: types.UserForm{
				Role:            "user",
				Username:        "subali_failure",
				Email:           "unit-test@.com",
				Name:            "Unit Test",
				Password:        "unit-test",
				ConfirmPassword: "unit-test",
			},
			status:     http.StatusBadRequest,
			flashError: "<strong>error:</strong> email: must be a valid email address.!",
		},
		{
			name: "user anonymous for POST form create to email too long failure: 2",
			token: echo.Map{
				"username": "anonymous",
				"role":     "anonymous",
			},
			form: types.UserForm{
				Role:            "user",
				Username:        "anony_failure",
				Email:           "unit-test@exemplecom",
				Name:            "Unit Test",
				Password:        "unit-test",
				ConfirmPassword: "unit-test",
			},
			status:     http.StatusBadRequest,
			flashError: "<strong>error:</strong> email: must be a valid email address.!",
		},
		// {
		// 	name:   "users [admin] to admin POST create it success: id=1",
		// 	expect: ADMIN,
		// 	path:   "1", // admin: 1 admin
		// 	form: types.UserForm{
		// 		Role:     "admin",
		// 		Username: "admin-success",
		// 	},
		// 	// redirect @route: /users
		// 	// HTTP response status: 200 OK
		// 	status: http.StatusOK,
		// 	// body navbar
		// 	htmlNavbar: `<a class="btn">ADMIN</a>`,
		// 	// body heading
		// 	htmlHeading: `<h2 class="mt-4">Users: All</h2>`,
		// 	// flash message success
		// 	htmlFlashSuccess: `<strong>success:</strong> success create user: admin-success!`,
		// },
		// {
		// 	name:   "users [admin] to user POST create it success: id=2",
		// 	expect: ADMIN,
		// 	path:   "2", // user: 2 sugriwa
		// 	form: types.UserForm{
		// 		// id=2 username: sugriwa
		// 		Role:     "user",
		// 		Username: "sugriwa",
		// 		Name:     "Sugriwa Success",
		// 	},
		// 	// redirect @route: /users
		// 	// HTTP response status: 200 OK
		// 	status: http.StatusOK,
		// 	// body navbar
		// 	htmlNavbar: `<a class="btn">ADMIN</a>`,
		// 	// body heading
		// 	htmlHeading: `<h2 class="mt-4">Users: All</h2>`,
		// 	// flash message success
		// 	// [admin] id=2 username: sugriwa
		// 	htmlFlashSuccess: `<strong>success:</strong> success create user: sugriwa!`,
		// },
		// {
		// 	name:   "users [admin] to POST create it failure: id=-1",
		// 	expect: ADMIN,
		// 	path:   "-1",
		// 	form:   types.UserForm{},
		// 	// HTTP response status: 404 Not Found
		// 	status:           http.StatusNotFound,
		// 	jsonMessageError: `{"message":"User Not Found"}`,
		// },

		// /*
		// 	create it [sugriwa]
		// */
		// // GET
		// {
		// 	name:   "users [sugriwa] to GET create it success: id=2",
		// 	expect: SUGRIWA,
		// 	method: http.MethodGet,
		// 	path:   "2", // user: 2 sugriwa ok
		// 	// HTTP response status: 200 OK
		// 	status: http.StatusOK,
		// 	// body heading
		// 	htmlHeading: regex{
		// 		mustCompile: `<h2 class="mt-4">(.*)</h2>`,
		// 		actual:      `<h2 class="mt-4">User: Sugriwa Success</h2>`,
		// 	},
		// },
		// {
		// 	name:   "users [sugriwa] to GET create it failure: id=-2",
		// 	expect: SUGRIWA,
		// 	method: http.MethodGet,
		// 	path:   "-2",
		// 	// HTTP response status: 404 Not Found
		// 	status: http.StatusNotFound,
		// 	jsonMessageError: regex{
		// 		mustCompile: `{"message":"(.*)"}`,
		// 		actual:      `{"message":"User Not Found"}`,
		// 	},
		// },
		// {
		// 	name:   "users [sugriwa] to GET create it failure: id=3",
		// 	expect: SUGRIWA,
		// 	method: http.MethodGet,
		// 	path:   "3", // user: 2 sugriwa no
		// 	// HTTP response status: 403 Forbidden,
		// 	status: http.StatusForbidden,
		// 	jsonMessageError: regex{
		// 		mustCompile: `{"message":"(.*)"}`,
		// 		actual:      `{"message":"Forbidden"}`,
		// 	},
		// },
		// // POST
		// // ?
		// {
		// 	name:   "users [sugriwa] to sugriwa POST create it success",
		// 	expect: SUGRIWA,
		// 	path:   "2", // user: 2 sugriwa
		// 	form: types.UserForm{
		// 		Username: "sugriwa", // admin: "sugriwa-success" to sugriwa: "sugriwa"
		// 		Name:     "Sugriwa",
		// 	},
		// 	// redirect @route: /
		// 	// HTTP response status: 200 OK
		// 	status: http.StatusOK,
		// 	// body heading
		// 	htmlHeading: regex{
		// 		mustCompile: `<h1 class="display-4">(.*)</h1>`,
		// 		actual:      `<h1 class="display-4">Hello Sugriwa!</h1>`,
		// 	},
		// 	// flash message success
		// 	htmlFlashSuccess: regex{
		// 		mustCompile: `<strong>success:</strong> (.*)`,
		// 		actual:      `<strong>success:</strong> success create user: sugriwa!`,
		// 	},
		// },
		// {
		// 	name:   "users [sugriwa] to POST create it failure",
		// 	expect: SUGRIWA,
		// 	path:   "3", // user: 2 sugriwa no
		// 	form: types.UserForm{
		// 		Username: "subali-failure",
		// 	},
		// 	// HTTP response status: 403 Forbidden
		// 	status: http.StatusForbidden,
		// 	jsonMessageError: regex{
		// 		mustCompile: `{"message":"(.*)"}`,
		// 		actual:      `{"message":"Forbidden"}`,
		// 	},
		// },

		// /*
		// 	create it [no-auth]
		// */
		// // GET
		// {
		// 	name:   "users [no-auth] to GET create it failure: id=1",
		// 	expect: "anonymous",
		// 	method: http.MethodGet,
		// 	path:   "1",
		// 	// redirect @route: /login
		// 	// HTTP response status: 200 OK
		// 	status: http.StatusOK,
		// 	// flash message
		// 	htmlFlashError: regex{
		// 		mustCompile: `<p class="text-danger">*(.*)</p>`,
		// 		actual:      `<p class="text-danger">*login process failed!</p>`,
		// 	},
		// },
		// {
		// 	name:   "users [no-auth] to GET create it failure: id=-1",
		// 	expect: "anonymous",
		// 	method: http.MethodGet,
		// 	path:   "-1",
		// 	// redirect @route: /login
		// 	// HTTP response status: 200 OK: 3 session and id
		// 	status: http.StatusOK,
		// },
		// // POST
		// {
		// 	name:   "users [no-auth] to POST create it failure: id=2",
		// 	expect: "anonymous",
		// 	path:   "2",
		// 	form: types.UserForm{
		// 		Username: "sugriwa-failure",
		// 	},
		// 	// redirect @route: /login
		// 	// HTTP response status: 200 OK: 3 session and id
		// 	status: http.StatusOK,
		// },
		// {
		// 	name:   "users [no-auth] to POST create it failure: id=-2",
		// 	expect: "anonymous",
		// 	path:   "-2",
		// 	form: types.UserForm{
		// 		Username: "sugriwa-failure",
		// 	},
		// 	// redirect @route: /login
		// 	// HTTP response status: 200 OK: 3 session and id
		// 	status: http.StatusOK,
		// },
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			token, _ := auth.GenerateToken(
				conf.AppJWTAuthSign,
				0,
				test.token["username"].(string),
				test.token["role"].(string),
			)
			auth := noAuth.Builder(func(req *httpexpect.Request) {
				req.WithCookie("token", token)
			})

			result = auth.POST("/users/add").
				WithForm(test.form).
				Expect().
				Status(test.status)

			assert.Contains(result.Body().Raw(), test.flashError)

			// resultBody := result.Body().Raw()

			// // assert.Equal(t, match, actual)
			// //
			// // or,
			// //
			// // assert := assert.New(t)
			// // ...
			// // assert.Equal(match, actual)
			// if test.htmlNavbar != "" {
			// 	assert.Equal(resultBody, test.htmlNavbar)
			// }
			// assert.Equal(resultBody, test.htmlHeading)

			// if test.htmlFlashSuccess != "" {
			// 	assert.Equal(resultBody, test.htmlFlashSuccess)
			// }

			// if test.htmlFlashError != "" {
			// 	assert.Equal(resultBody, test.htmlFlashError)
			// }

			// if test.jsonMessageError != "" {
			// 	assert.Equal(resultBody, test.jsonMessageError)
			// }

			// statusCode := result.Raw().StatusCode
			// if test.status != statusCode {
			// 	t.Logf(
			// 		"got: %d but expect %d", test.status, statusCode,
			// 	)
			// 	t.Fail()
			// }
		})
	}

	// test for db users
	truncateUsers()
}

// Database: " Error 1062: Duplicate entry 'unit-test@exemple.com' for key 'email_UNIQUE' " v
//			-> " Error 1062: Duplicate entry 'unit-test@exemple.com' for key 'users.email_UNIQUE' " x

// func TestUpdateUserController(t *testing.T) {
// 	// test for db users
// 	truncateUsers()

// 	// assert
// 	assert := assert.New(t)

// 	noAuth := setupTestServer(t)
// 	testCases := []struct {
// 		name   string
// 		expect string // auth or no-auth
// 		method string // method: 1=GET or 2=POST
// 		path   string // id=string. Exemple, id="1"
// 		form   types.UserForm
// 		status int

// 		htmlNavbar  string
// 		htmlHeading string
// 		// flash message
// 		htmlFlashSuccess string
// 		htmlFlashError   string

// 		jsonMessageError string
// 	}{
// 		/*
// 			update it [admin]
// 		*/
// 		// GET
// 		{
// 			name:   "users [admin] to admin GET update it success: id=1",
// 			expect: ADMIN,
// 			method: http.MethodGet,
// 			path:   "1", // admin: 1 admin
// 			// HTTP response status: 200 OK
// 			status: http.StatusOK,
// 			// body navbar
// 			htmlNavbar: `<a class="btn">ADMIN</a>`,
// 			// body heading
// 			htmlHeading: `<h2 class="mt-4">User: Admin</h2>`,
// 		},
// 		{
// 			name:   "users [admin] to user GET update it success: id=2",
// 			expect: ADMIN,
// 			method: http.MethodGet,
// 			path:   "2", // user: 2 sugriwa
// 			// HTTP response status: 200 OK
// 			status: http.StatusOK,
// 			// body navbar
// 			htmlNavbar: `<a class="btn">ADMIN</a>`,
// 			// body heading
// 			htmlHeading: `<h2 class="mt-4">User: Sugriwa</h2>`,
// 		},
// 		{
// 			name:   "users [admin] to -1 GET update it failure: id=-1",
// 			expect: ADMIN,
// 			method: http.MethodGet,
// 			path:   "-1",
// 			// HTTP response status: 404 Not Found
// 			status:           http.StatusNotFound,
// 			jsonMessageError: `{"message":"User Not Found"}`,
// 		},
// 		// POST
// 		{
// 			name:   "users [admin] to admin POST update it success: id=1",
// 			expect: ADMIN,
// 			method: http.MethodPost,
// 			path:   "1", // admin: 1 admin
// 			form: types.UserForm{
// 				Role:     "admin",
// 				Username: "admin-success",
// 			},
// 			// redirect @route: /users
// 			// HTTP response status: 200 OK
// 			status: http.StatusOK,
// 			// body navbar
// 			htmlNavbar: `<a class="btn">ADMIN</a>`,
// 			// body heading
// 			htmlHeading: `<h2 class="mt-4">Users: All</h2>`,
// 			// flash message success
// 			htmlFlashSuccess: `<strong>success:</strong> success update user: admin-success!`,
// 		},
// 		{
// 			name:   "users [admin] to user POST update it success: id=2",
// 			expect: ADMIN,
// 			method: http.MethodPost,
// 			path:   "2", // user: 2 sugriwa
// 			form: types.UserForm{
// 				// id=2 username: sugriwa
// 				Role:     "user",
// 				Username: "sugriwa",
// 				Name:     "Sugriwa Success",
// 			},
// 			// redirect @route: /users
// 			// HTTP response status: 200 OK
// 			status: http.StatusOK,
// 			// body navbar
// 			htmlNavbar: `<a class="btn">ADMIN</a>`,
// 			// body heading
// 			htmlHeading: `<h2 class="mt-4">Users: All</h2>`,
// 			// flash message success
// 			htmlFlashSuccess: `<strong>success:</strong> success update user: sugriwa!`,
// 		},
// 		{
// 			name:   "users [admin] to POST update it failure: id=-1",
// 			expect: ADMIN,
// 			method: http.MethodPost,
// 			path:   "-1",
// 			form:   types.UserForm{},
// 			// HTTP response status: 404 Not Found
// 			status:           http.StatusNotFound,
// 			jsonMessageError: `{"message":"User Not Found"}`,
// 		},

// 		/*
// 			update it [sugriwa]
// 		*/
// 		// GET
// 		{
// 			name:   "users [sugriwa] to GET update it success: id=2",
// 			expect: SUGRIWA,
// 			method: http.MethodGet,
// 			path:   "2", // user: 2 sugriwa ok
// 			// HTTP response status: 200 OK
// 			status: http.StatusOK,
// 			// body heading
// 			htmlHeading: `<h2 class="mt-4">User: Sugriwa Success</h2>`,
// 		},
// 		{
// 			name:   "users [sugriwa] to GET update it failure: id=-2",
// 			expect: SUGRIWA,
// 			method: http.MethodGet,
// 			path:   "-2",
// 			// HTTP response status: 404 Not Found
// 			status:           http.StatusNotFound,
// 			jsonMessageError: `{"message":"User Not Found"}`,
// 		},
// 		{
// 			name:   "users [sugriwa] to GET update it failure: id=3",
// 			expect: SUGRIWA,
// 			method: http.MethodGet,
// 			path:   "3", // user: 2 sugriwa no
// 			// HTTP response status: 403 Forbidden,
// 			status:           http.StatusForbidden,
// 			jsonMessageError: `{"message":"Forbidden"}`,
// 		},
// 		// POST
// 		// ?
// 		{
// 			name:   "users [sugriwa] to sugriwa POST update it success",
// 			expect: SUGRIWA,
// 			method: http.MethodPost,
// 			path:   "2", // user: 2 sugriwa
// 			form: types.UserForm{
// 				Username: "sugriwa", // admin: "sugriwa-success" to sugriwa: "sugriwa"
// 				Name:     "Sugriwa",
// 			},
// 			// redirect @route: /
// 			// HTTP response status: 200 OK
// 			status: http.StatusOK,
// 			// body heading
// 			htmlHeading: `<h1 class="display-4">Hello Sugriwa!</h1>`,
// 			// flash message success
// 			htmlFlashSuccess: `<strong>success:</strong> success update user: sugriwa!`,
// 		},
// 		{
// 			name:   "users [sugriwa] to POST update it failure",
// 			expect: SUGRIWA,
// 			method: http.MethodPost,
// 			path:   "3", // user: 2 sugriwa no
// 			form: types.UserForm{
// 				Username: "subali-failure",
// 			},
// 			// HTTP response status: 403 Forbidden
// 			status:           http.StatusForbidden,
// 			jsonMessageError: `{"message":"Forbidden"}`,
// 		},

// 		/*
// 			update it [no-auth]
// 		*/
// 		// GET
// 		{
// 			name:   "users [no-auth] to GET update it failure: id=1",
// 			expect: "anonymous",
// 			method: http.MethodGet,
// 			path:   "1",
// 			// redirect @route: /login
// 			// HTTP response status: 200 OK
// 			status: http.StatusOK,
// 			// flash message
// 			htmlFlashError: `<p class="text-danger">*login process failed!</p>`,
// 		},
// 		{
// 			name:   "users [no-auth] to GET update it failure: id=-1",
// 			expect: "anonymous",
// 			method: http.MethodGet,
// 			path:   "-1",
// 			// redirect @route: /login
// 			// HTTP response status: 200 OK: 3 session and id
// 			status: http.StatusOK,
// 		},
// 		// POST
// 		{
// 			name:   "users [no-auth] to POST update it failure: id=2",
// 			expect: "anonymous",
// 			method: http.MethodPost,
// 			path:   "2",
// 			form: types.UserForm{
// 				Username: "sugriwa-failure",
// 			},
// 			// redirect @route: /login
// 			// HTTP response status: 200 OK: 3 session and id
// 			status: http.StatusOK,
// 		},
// 		{
// 			name:   "users [no-auth] to POST update it failure: id=-2",
// 			expect: "anonymous",
// 			method: http.MethodPost,
// 			path:   "-2",
// 			form: types.UserForm{
// 				Username: "sugriwa-failure",
// 			},
// 			// redirect @route: /login
// 			// HTTP response status: 200 OK: 3 session and id
// 			status: http.StatusOK,
// 		},
// 	}

// 	for _, test := range testCases {
// 		modelsTest.UserSelectTest = test.expect // ADMIN and SUGRIWA

// 		t.Run(test.name, func(t *testing.T) {
// 			var result *httpexpect.Response
// 			if test.method == http.MethodGet {
// 				// same:
// 				//
// 				// noAuth.GET("/users/view/{id}").
// 				//	WithPath("id", test.path).
// 				// ...
// 				result = noAuth.GET("/users/view/{id}", test.path).
// 					WithForm(test.form).
// 					Expect().
// 					Status(test.status)
// 			} else if test.method == http.MethodPost {
// 				result = noAuth.POST("/users/view/{id}").
// 					WithPath("id", test.path).
// 					WithForm(test.form).
// 					Expect().
// 					Status(test.status)
// 			} else {
// 				panic("method: 1=GET or 2=POST")
// 			}

// 			resultBody := result.Body().Raw()

// 			assert.Equal(resultBody, test.htmlNavbar)
// 			assert.Equal(resultBody, test.htmlHeading)

// 			if test.htmlFlashSuccess != "" {
// 				assert.Equal(resultBody, test.htmlFlashSuccess)
// 			}

// 			if test.htmlFlashError != "" {
// 				assert.Equal(resultBody, test.htmlFlashError)
// 			}

// 			if test.jsonMessageError != "" {
// 				assert.Equal(resultBody, test.jsonMessageError)
// 			}

// 			statusCode := result.Raw().StatusCode
// 			if test.status != statusCode {
// 				t.Logf(
// 					"got: %d but expect %d", test.status, statusCode,
// 				)
// 				t.Fail()
// 			}
// 		})
// 	}

// 	// test for db users
// 	truncateUsers()
// }

// func TestUpdateUserByPasswordUserController(t *testing.T) {
// 	assert := assert.New(t)

// 	noAuth := setupTestServer(t)

// 	// test for SetSession = false
// 	method.SetSession = false
// 	// test for db users
// 	truncateUsers()

// 	testCases := []struct {
// 		name   string
// 		expect string // ADMIN and SUGRIWA
// 		method string // method: 1=GET or 2=POST
// 		path   string // id=string. Exemple, id="1"
// 		form   types.NewPasswordForm
// 		status int

// 		htmlHeading regex
// 		// flash message
// 		htmlFlashSuccess regex
// 		htmlFlashError   regex

// 		jsonMessageError regex
// 	}{
// 		/*
// 			update by password it [admin]
// 		*/
// 		// GET
// 		{
// 			name:   "users [admin] to GET update user by password it success: id=1",
// 			expect: ADMIN,
// 			method: http.MethodGet,
// 			path:   "1",
// 			// HTTP response status: 200 OK
// 			status: http.StatusOK,
// 			// body heading
// 			htmlHeading: regex{
// 				mustCompile: `<h3 class="mt-4">(.*)</h3>`,
// 				actual:      `<h3 class="mt-4">User: Admin</h3>`,
// 			},
// 		},
// 		{
// 			name:   "users [admin] to [sugriwa] GET update user by password it success: id=2",
// 			expect: ADMIN,
// 			method: http.MethodGet,
// 			path:   "2",
// 			// HTTP response status: 200 OK
// 			status: http.StatusOK,
// 			// body heading
// 			htmlHeading: regex{
// 				mustCompile: `<h3 class="mt-4">(.*)</h3>`,
// 				actual:      `<h3 class="mt-4">User: Sugriwa</h3>`,
// 			},
// 		},
// 		{
// 			name: "users [admin] to GET update user by password it failure: id=-1" +
// 				" GET passwords don't match",
// 			expect: ADMIN,
// 			method: http.MethodGet,
// 			path:   "-1",
// 			// HTTP response status: 404 Not Found
// 			status: http.StatusNotFound,
// 			jsonMessageError: regex{
// 				mustCompile: `{"message":"(.*)"}`,
// 				actual:      `{"message":"User Not Found"}`,
// 			},
// 		},
// 		// POST
// 		{
// 			name:   "users [admin] to POST update user by password it success: id=1",
// 			expect: ADMIN,
// 			method: http.MethodPost,
// 			path:   "1",
// 			form: types.NewPasswordForm{
// 				OldPassword:        "admin123",
// 				NewPassword:        "admin_success",
// 				ConfirmNewPassword: "admin_success",
// 			},
// 			// HTTP response status: 200 OK
// 			status: http.StatusOK,
// 			// body heading
// 			htmlHeading: regex{
// 				mustCompile: `<h2 class="mt-4">(.*)</h2>`,
// 				actual:      `<h2 class="mt-4">Users: All</h2>`,
// 			},
// 			// flash message success
// 			htmlFlashSuccess: regex{
// 				mustCompile: `<strong>success:</strong> (.*)`,
// 				actual:      `<strong>success:</strong> success update user by password: admin!`,
// 			},
// 		},
// 		{
// 			name:   "users [admin] to [sugriwa] POST update user by password it success: id=2",
// 			expect: ADMIN,
// 			method: http.MethodPost,
// 			path:   "2",
// 			form: types.NewPasswordForm{
// 				OldPassword:        "user123",
// 				NewPassword:        "user_success",
// 				ConfirmNewPassword: "user_success",
// 			},
// 			// HTTP response status: 200 OK
// 			status: http.StatusOK,
// 			// body heading
// 			htmlHeading: regex{
// 				mustCompile: `<h2 class="mt-4">(.*)</h2>`,
// 				actual:      `<h2 class="mt-4">Users: All</h2>`,
// 			},
// 			// flash message success
// 			htmlFlashSuccess: regex{
// 				mustCompile: `<strong>success:</strong> (.*)`,
// 				actual:      `<strong>success:</strong> success update user by password: sugriwa!`,
// 			},
// 		},
// 		{
// 			name: "users [admin] to POST update user by password it failure: id=1" +
// 				" POST passwords don't match",
// 			expect: ADMIN,
// 			method: http.MethodPost,
// 			path:   "1",
// 			form: types.NewPasswordForm{
// 				OldPassword:        "admin_success",
// 				NewPassword:        "admin_success_",
// 				ConfirmNewPassword: "admin_failure",
// 			},
// 			// HTTP response status: 403 Forbidden
// 			status: http.StatusForbidden,
// 		},
// 		{
// 			name: "users [admin] to [sugriwa] POST update user by password it failure: id=2" +
// 				" POST passwords don't match",
// 			expect: ADMIN,
// 			method: http.MethodPost,
// 			path:   "2",
// 			form: types.NewPasswordForm{
// 				OldPassword:        "admin_password_success",
// 				NewPassword:        "admin_password_failure",
// 				ConfirmNewPassword: "admin_password_success_",
// 			},
// 			// HTTP response status: 403 Forbidden
// 			status: http.StatusForbidden,
// 		},
// 		{
// 			name:   "users [admin] to POST update user by password it failure: id=-1",
// 			expect: ADMIN,
// 			method: http.MethodPost,
// 			path:   "-1",
// 			form:   types.NewPasswordForm{},
// 			// HTTP response status: 404 Not Found
// 			status: http.StatusNotFound,
// 			jsonMessageError: regex{
// 				mustCompile: `{"message":"(.*)"}`,
// 				actual:      `{"message":"User Not Found"}`,
// 			},
// 		},

// 		/*
// 			update by password it [sugriwa]
// 		*/
// 		// GET
// 		{
// 			name:   "users [sugriwa] to GET update user by password it success: id=2",
// 			expect: SUGRIWA,
// 			method: http.MethodGet,
// 			path:   "2",
// 			// HTTP response status: 200 OK
// 			status: http.StatusOK,
// 			// body heading
// 			htmlHeading: regex{
// 				mustCompile: `<h3 class="mt-4">(.*)</h3>`,
// 				actual:      `<h3 class="mt-4">User: Sugriwa</h3>`,
// 			},
// 		},
// 		{
// 			name:   "users [sugriwa] to [admin] GET update user by password it failure: id=1",
// 			expect: SUGRIWA,
// 			method: http.MethodGet,
// 			path:   "1",
// 			// HTTP response status: 403 Forbidden
// 			status: http.StatusForbidden,
// 			jsonMessageError: regex{
// 				mustCompile: `{"message":"(.*)"}`,
// 				actual:      `{"message":"User Not Found"}`,
// 			},
// 		},
// 		{
// 			name:   "users [sugriwa] to [subali] GET update user by password it failure: id=3",
// 			expect: SUGRIWA,
// 			method: http.MethodGet,
// 			path:   "3",
// 			// HTTP response status: 403 Forbidden
// 			status: http.StatusForbidden,
// 			jsonMessageError: regex{
// 				mustCompile: `{"message":"(.*)"}`,
// 				actual:      `{"message":"User Not Found"}`,
// 			},
// 		},
// 		{
// 			name:   "users [sugriwa] to GET update user by password it failure: id=-1",
// 			expect: SUGRIWA,
// 			method: http.MethodGet,
// 			path:   "-1",
// 			// HTTP response status: 404 Not Found
// 			status: http.StatusNotFound,
// 			jsonMessageError: regex{
// 				mustCompile: `{"message":"(.*)"}`,
// 				actual:      `{"message":"User Not Found"}`,
// 			},
// 		},
// 		// POST
// 		{
// 			name:   "users [sugriwa] to POST update user by password it success: id=2",
// 			expect: SUGRIWA,
// 			method: http.MethodPost,
// 			path:   "2",
// 			form: types.NewPasswordForm{
// 				OldPassword:        "user_success",
// 				NewPassword:        "user123",
// 				ConfirmNewPassword: "user123",
// 			},
// 			// HTTP response status: 200 OK
// 			status: http.StatusOK,
// 			// flash message success
// 			htmlFlashSuccess: regex{
// 				mustCompile: `<strong>success:</strong> (.*)`,
// 				actual:      `<strong>success:</strong> success update user by password: sugriwa!`,
// 			},
// 		},
// 		{
// 			name:   "users [sugriwa] to [admin] POST update user by password it failure: id=1",
// 			expect: SUGRIWA,
// 			method: http.MethodPost,
// 			path:   "1",
// 			form:   types.NewPasswordForm{},
// 			// HTTP response status: 403 Forbidden,
// 			status: http.StatusForbidden,
// 			jsonMessageError: regex{
// 				mustCompile: `{"message":"(.*)"}`,
// 				actual:      `{"message":"User Not Found"}`,
// 			},
// 		},
// 		{
// 			name:   "users [sugriwa] to [subali] POST update user by password it failure: id=3",
// 			expect: SUGRIWA,
// 			method: http.MethodPost,
// 			path:   "3",
// 			form:   types.NewPasswordForm{},
// 			// HTTP response status: 403 Forbidden,
// 			status: http.StatusForbidden,
// 			jsonMessageError: regex{
// 				mustCompile: `{"message":"(.*)"}`,
// 				actual:      `{"message":"User Not Found"}`,
// 			},
// 		},
// 		{
// 			name:   "users [sugriwa] to POST update user by password it failure: id=-1",
// 			expect: SUGRIWA,
// 			method: http.MethodPost,
// 			path:   "-1",
// 			form:   types.NewPasswordForm{},
// 			// HTTP response status: 404 Not Found
// 			status: http.StatusNotFound,
// 			jsonMessageError: regex{
// 				mustCompile: `{"message":"(.*)"}`,
// 				actual:      `{"message":"User Not Found"}`,
// 			},
// 		},

// 		/*
// 			update by password it [no-auth]
// 		*/
// 		// GET
// 		{
// 			name:   "users [no-auth] to GET update user by password it failure: id=1",
// 			expect: ANONYMOUS,
// 			method: http.MethodGet,
// 			path:   "1",
// 			// redirect @route: /login
// 			// HTTP response status: 200 OK
// 			status: http.StatusOK,
// 			// flash message
// 			htmlFlashError: regex{
// 				mustCompile: `<p class="text-danger">*(.*)</p>`,
// 				actual:      `<p class="text-danger">*login process failed!</p>`,
// 			},
// 		},
// 		{
// 			name:   "users [no-auth] to POST update user by password it failure: id=-1",
// 			expect: ANONYMOUS,
// 			method: http.MethodGet,
// 			path:   "-1",
// 			// redirect @route: /login
// 			// HTTP response status: 200 OK
// 			status: http.StatusOK,
// 			// flash message
// 			htmlFlashError: regex{
// 				mustCompile: `<p class="text-danger">*(.*)</p>`,
// 				actual:      `<p class="text-danger">*login process failed!</p>`,
// 			},
// 		},
// 		// POST
// 		{
// 			name:   "users [no-auth] to POST update user by password it failure: id=1",
// 			expect: ANONYMOUS,
// 			method: http.MethodPost,
// 			path:   "1",
// 			// redirect @route: /login
// 			// HTTP response status: 200 OK
// 			status: http.StatusOK,
// 			form:   types.NewPasswordForm{},
// 			// flash message
// 			htmlFlashError: regex{
// 				mustCompile: `<p class="text-danger">*(.*)</p>`,
// 				actual:      `<p class="text-danger">*login process failed!</p>`,
// 			},
// 		},
// 		{
// 			name:   "users [no-auth] to POST update user by password it success: id=-1",
// 			expect: ANONYMOUS,
// 			method: http.MethodPost,
// 			path:   "-1",
// 			// redirect @route: /login
// 			// HTTP response status: 200 OK
// 			status: http.StatusOK,
// 			form:   types.NewPasswordForm{},
// 			// flash message
// 			htmlFlashError: regex{
// 				mustCompile: `<p class="text-danger">*(.*)</p>`,
// 				actual:      `<p class="text-danger">*login process failed!</p>`,
// 			},
// 		},
// 	}

// 	// for...{...}, equal:
// 	//
// 	// t.Run("users [auth] to POST update user by password it success", func(t *testing.T) {
// 	// 	auth.POST("/users/view/{id}/password").
// 	// 		WithPath("id", "1").
// 	// 		WithForm(types.NewPasswordForm{
// 	// 			...
// 	// 		}).
// 	// 		Expect().
// 	// 		Status(http.StatusOK)
// 	// })
// 	//
// 	// ...
// 	//
// 	// t.Run("users [no-auth] to POST update user by password it failure: 4"+
// 	// 	" no session", func(t *testing.T) {
// 	// 	noAuth.POST("/users/view/{id}/password").
// 	// 		WithPath("id", "1").
// 	// 		Expect().
// 	// 		// redirect @route: /login
// 	// 		// HTTP response status: 200 OK
// 	// 		Status(http.StatusOK)
// 	// })
// 	for _, test := range testCases {
// 		modelsTest.UserSelectTest = test.expect // ADMIN and SUGRIWA

// 		var result *httpexpect.Response
// 		t.Run(test.name, func(t *testing.T) {
// 			if test.method == http.MethodGet {
// 				// equal:
// 				//
// 				// noAuth.POST("/users/view/{id}/password").
// 				//	WithPath("id", test.path).
// 				// ...
// 				result = noAuth.GET("/users/view/{id}/password", test.path).
// 					WithForm(test.form).
// 					Expect().
// 					Status(test.status)
// 			} else if test.method == http.MethodPost {
// 				result = noAuth.POST("/users/view/{id}/password").
// 					WithPath("id", test.path).
// 					WithForm(test.form).
// 					Expect().
// 					Status(test.status)
// 			} else {
// 				panic("method: 1=GET or 2=POST")
// 			}

// 			resultBody := result.Body().Raw()

// 			var (
// 				mustCompile, actual, match string
// 				regex                      *regexp.Regexp
// 			)

// 			if test.htmlHeading.mustCompile != "" {
// 				mustCompile = test.htmlHeading.mustCompile
// 				actual = test.htmlHeading.actual

// 				regex = regexp.MustCompile(mustCompile)
// 				match = regex.FindString(resultBody)

// 				assert.Equal(match, actual)
// 			}

// 			if test.htmlFlashSuccess.mustCompile != "" {
// 				mustCompile = test.htmlFlashSuccess.mustCompile
// 				actual = test.htmlFlashSuccess.actual

// 				regex = regexp.MustCompile(mustCompile)
// 				match = regex.FindString(resultBody)

// 				assert.Equal(match, actual)
// 			}

// 			if test.htmlFlashError.mustCompile != "" {
// 				mustCompile = test.htmlFlashError.mustCompile
// 				actual = test.htmlFlashError.actual

// 				regex = regexp.MustCompile(mustCompile)
// 				match = regex.FindString(resultBody)

// 				assert.Equal(match, actual)
// 			}

// 			if test.jsonMessageError.mustCompile != "" {
// 				mustCompile = test.jsonMessageError.mustCompile
// 				actual = test.jsonMessageError.actual

// 				regex = regexp.MustCompile(mustCompile)
// 				match = regex.FindString(resultBody)

// 				assert.Equal(match, actual)
// 			}

// 			statusCode := result.Raw().StatusCode
// 			if test.status != statusCode {
// 				t.Logf(
// 					"got: %d but expect %d", test.status, statusCode,
// 				)
// 				t.Fail()
// 			}
// 		})
// 	}

// 	// test for db users
// 	truncateUsers()
// }

// // TODO: Test Delete User Controller, insyaallah
// func TestDeleteUserController(t *testing.T) {
// 	assert := assert.New(t)

// 	noAuth := setupTestServer(t)

// 	// test for SetSession = false
// 	method.SetSession = false
// 	// test for db users
// 	truncateUsers()

// 	testCases := []struct {
// 		name           string
// 		expect         string // ADMIN and SUBALI
// 		path           string // id=string. Exemple, id="1"
// 		setSessionTrue bool
// 		status         int

// 		htmlHeading regex
// 		// flash message
// 		htmlFlashSuccess regex
// 		htmlFlashError   regex
// 		jsonMessageError regex
// 	}{
// 		// GET all
// 		/*
// 			delete it [admin]
// 		*/
// 		{
// 			name:   "users [admin] to [admin] DELETE it failure: id=1",
// 			expect: ADMIN,
// 			path:   "1",
// 			// HTTP response status: 403 Forbidden,
// 			status: http.StatusForbidden,
// 			jsonMessageError: regex{
// 				mustCompile: `{"message":"(.*)"}`,
// 				actual:      `{"message":"Forbidden"}`,
// 			},
// 		},
// 		{
// 			name:   "users [admin] to [sugriwa] DELETE it success: id=2",
// 			expect: ADMIN,
// 			path:   "2",
// 			// redirect @route: /users
// 			// HTTP response status: 200 OK
// 			status: http.StatusOK,
// 			// body heading
// 			htmlHeading: regex{
// 				mustCompile: `<h2 class="mt-4">(.*)</h2>`,
// 				actual:      `<h2 class="mt-4">Users: All</h2>`,
// 			},
// 			// flash message success
// 			htmlFlashSuccess: regex{
// 				mustCompile: `<strong>success:</strong> (.*)`,
// 				actual:      `<strong>success:</strong> success delete user: sugriwa!`,
// 			},
// 		},
// 		{
// 			name:   "users [admin] to [sugriwa] DELETE it failure: id=2 delete exists",
// 			expect: ADMIN,
// 			path:   "2",
// 			// HTTP response status: 404 Not Found
// 			status: http.StatusNotFound,
// 			jsonMessageError: regex{
// 				mustCompile: `{"message":"(.*)"}`,
// 				actual:      `{"message":"User Not Found"}`,
// 			},
// 		},
// 		{
// 			name:   "users [admin] to DELETE it failure: 2 (id=-1)",
// 			expect: ADMIN,
// 			path:   "-1",
// 			// HTTP response status: 404 Not Found
// 			status: http.StatusNotFound,
// 			jsonMessageError: regex{
// 				mustCompile: `{"message":"(.*)"}`,
// 				actual:      `{"message":"User Not Found"}`,
// 			},
// 		},

// 		/*
// 		   delete it [subali]
// 		*/
// 		{
// 			name:   "users [subali] to [admin] DELETE it failure: id=1",
// 			expect: SUBALI,
// 			path:   "1",
// 			// HTTP response status: 403 Forbidden,
// 			status: http.StatusForbidden,
// 			jsonMessageError: regex{
// 				mustCompile: `{"message":"(.*)"}`,
// 				actual:      `{"message":"Forbidden"}`,
// 			},
// 		},
// 		{
// 			name:   "users [subali] to DELETE it failure: id=-1",
// 			expect: SUBALI,
// 			path:   "-1",
// 			// HTTP response status: 404 Not Found
// 			status: http.StatusNotFound,
// 			jsonMessageError: regex{
// 				mustCompile: `{"message":"(.*)"}`,
// 				actual:      `{"message":"User Not Found"}`,
// 			},
// 		},
// 		// {
// 		// 	name:             "users [subali] to [subali] DELETE it success: id=3",
// 		// 	expect:           SUBALI,
// 		// 	path:             "3",
// 		// 	set_session_true: true,
// 		// 	// redirect @route: /
// 		// 	// HTTP response status: 200 OK
// 		// 	status: http.StatusOK,
// 		// 	// body heading
// 		// 	htmlHeading: regex{
// 		// 		must_compile: `<p class="lead">(.*)</p>`,
// 		// 		actual:       `<p class="lead">Test.</p>`,
// 		// 	},
// 		// 	// flash message success
// 		// 	html_flash_success: regex{
// 		// 		must_compile: `<strong>success:</strong> (.*)`,
// 		// 		actual:       `<strong>success:</strong> success delete user: subali!`,
// 		// 	},
// 		// },

// 		/*
// 		   delete it [na-auth]
// 		*/
// 		{
// 			name:   "users [no-auth] to DELETE it failure: id=1",
// 			expect: ANONYMOUS,
// 			path:   "1",
// 			// redirect @route: /login
// 			// HTTP response status: 200 OK
// 			status: http.StatusOK,
// 			// flash message
// 			htmlFlashError: regex{
// 				mustCompile: `<p class="text-danger">*(.*)</p>`,
// 				actual:      `<p class="text-danger">*login process failed!</p>`,
// 			},
// 		},
// 		{
// 			name:   "users [no-auth] to DELETE it failure: id=-1",
// 			expect: ANONYMOUS,
// 			path:   "-1",
// 			// redirect @route: /login
// 			// HTTP response status: 200 OK
// 			status: http.StatusOK,
// 			// flash message
// 			htmlFlashError: regex{
// 				mustCompile: `<p class="text-danger">*(.*)</p>`,
// 				actual:      `<p class="text-danger">*login process failed!</p>`,
// 			},
// 		},
// 		{
// 			name:   "users [no-auth] to DELETE it failure: id=error",
// 			expect: ANONYMOUS,
// 			path:   "error",
// 			// redirect @route: /login
// 			// HTTP response status: 200 OK
// 			status: http.StatusOK,
// 			// flash message
// 			htmlFlashError: regex{
// 				mustCompile: `<p class="text-danger">*(.*)</p>`,
// 				actual:      `<p class="text-danger">*login process failed!</p>`,
// 			},
// 		},
// 	}

// 	for _, test := range testCases {
// 		var result *httpexpect.Response

// 		if test.setSessionTrue {
// 			method.SetSession = true
// 		}

// 		t.Run(test.name, func(t *testing.T) {
// 			modelsTest.UserSelectTest = test.expect // ADMIN and SUBALI

// 			result = noAuth.GET("/users/delete/{id}", test.path).
// 				Expect().
// 				Status(test.status)

// 			resultBody := result.Body().Raw()

// 			var (
// 				mustCompile, actual, match string
// 				regex                      *regexp.Regexp
// 			)

// 			if test.htmlHeading.mustCompile != "" {
// 				mustCompile = test.htmlHeading.mustCompile
// 				actual = test.htmlHeading.actual

// 				regex = regexp.MustCompile(mustCompile)
// 				match = regex.FindString(resultBody)

// 				assert.Equal(match, actual)
// 			}

// 			if test.htmlFlashSuccess.mustCompile != "" {
// 				mustCompile = test.htmlFlashSuccess.mustCompile
// 				actual = test.htmlFlashSuccess.actual

// 				regex = regexp.MustCompile(mustCompile)
// 				match = regex.FindString(resultBody)

// 				assert.Equal(match, actual)
// 			}

// 			if test.htmlFlashError.mustCompile != "" {
// 				mustCompile = test.htmlFlashError.mustCompile
// 				actual = test.htmlFlashError.actual

// 				regex = regexp.MustCompile(mustCompile)
// 				match = regex.FindString(resultBody)

// 				assert.Equal(match, actual)
// 			}

// 			if test.jsonMessageError.mustCompile != "" {
// 				mustCompile = test.jsonMessageError.mustCompile
// 				actual = test.jsonMessageError.actual

// 				regex = regexp.MustCompile(mustCompile)
// 				match = regex.FindString(resultBody)

// 				assert.Equal(match, actual)
// 			}

// 			statusCode := result.Raw().StatusCode
// 			if test.status != statusCode {
// 				t.Logf(
// 					"got: %d but expect %d", test.status, statusCode,
// 				)
// 				t.Fail()
// 			}
// 		})
// 	}

// 	// test for db users
// 	truncateUsers()
// }
