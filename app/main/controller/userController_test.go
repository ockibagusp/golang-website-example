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
		name  string
		token echo.Map
		form  types.UserForm
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
				// HTTP response status: 200 OK
				Status(http.StatusOK)
		})
	}

	// test for db users
	truncateUsers()
}

// Database: " Error 1062: Duplicate entry 'unit-test@exemple.com' for key 'email_UNIQUE' " v
//
//	-> " Error 1062: Duplicate entry 'unit-test@exemple.com' for key 'users.email_UNIQUE' " x
func TestCreateUsers_WithInputPOSTFormEmailDuplicateEntryFailure(t *testing.T) {
	// assert
	assert := assert.New(t)

	noAuth := setupTestServer(t)
	// new users "unit-test"
	noAuth.POST("/users/add").
		WithForm(types.UserForm{
			Role:            "user",
			Username:        "unit-test",
			Email:           "unit-test@exemple.com",
			Name:            "Unit Test",
			Password:        "unit-test",
			ConfirmPassword: "unit-test",
		}).
		Expect().
		Status(http.StatusOK)

	// Database: " Error 1062: Duplicate entry 'unit-test@exemple.com' for key 'email_UNIQUE'
	noAuth = setupTestServer(t)
	result := noAuth.POST("/users/add").
		WithForm(types.UserForm{
			Role:            "user",
			Username:        "unit-test",
			Email:           "unit-test@exemple.com",
			Name:            "Unit Test",
			Password:        "unit-test",
			ConfirmPassword: "unit-test",
		}).
		Expect().
		Status(http.StatusBadRequest)

	assert.Contains(result.Body().Raw(), "<strong>error:</strong> Error 1062 (23000): Duplicate entry &#39;unit-test@exemple.com&#39; for key &#39;users.email_UNIQUE&#39;!")

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
		name           string
		token          echo.Map
		form           types.UserForm
		htmlFlashError string
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
				Username:        "subali_copy_failure", // look
				Email:           "unit-test@exemple.com",
				Name:            "Unit Test",
				Password:        "unit-test",
				ConfirmPassword: "unit-test",
			},
			htmlFlashError: "<strong>error:</strong> username: the length must be between 4 and 15.!",
		},
		{
			name: "user anonymous for POST form create to username too long failure: 2",
			token: echo.Map{
				"username": "anonymous",
				"role":     "anonymous",
			},
			form: types.UserForm{
				Role:            "user",
				Username:        "anonymous_failure", // look
				Email:           "unit-test@exemple.com",
				Name:            "Unit Test",
				Password:        "unit-test",
				ConfirmPassword: "unit-test",
			},
			htmlFlashError: "<strong>error:</strong> username: the length must be between 4 and 15.!",
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
				Email:           "unit-test@.com", // look
				Name:            "Unit Test",
				Password:        "unit-test",
				ConfirmPassword: "unit-test",
			},
			htmlFlashError: "<strong>error:</strong> email: must be a valid email address.!",
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
				Email:           "unit-test@exemplecom", // look
				Name:            "Unit Test",
				Password:        "unit-test",
				ConfirmPassword: "unit-test",
			},
			htmlFlashError: "<strong>error:</strong> email: must be a valid email address.!",
		},
		/*
			create form password it's not confirm_password failure
		*/
		{
			name: "user admin [admin] for POST password it's not confirm_password failure: 4",
			token: echo.Map{
				"username": "admin",
				"role":     "admin",
			},
			form: types.UserForm{
				Role:            "user",
				Username:        "subali",
				Email:           "unit-test@exemple.com",
				Name:            "Unit Test",
				Password:        "unit-failure", // look
				ConfirmPassword: "unit-test",    // look
			},
			htmlFlashError: " <strong>error:</strong> password: passwords don&#39;t match.!",
		},
		{
			name: "user anonymous for POST password it's not confirm_password failure: 5",
			token: echo.Map{
				"username": "anonymous",
				"role":     "anonymous",
			},
			form: types.UserForm{
				Role:            "user",
				Username:        "anonymous",
				Email:           "unit-test@exemple.com",
				Name:            "Unit Test",
				Password:        "unit-test",    // look
				ConfirmPassword: "unit-failure", // look
			},
			htmlFlashError: " <strong>error:</strong> password: passwords don&#39;t match.!",
		},
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
				Status(http.StatusBadRequest)

			assert.Contains(result.Body().Raw(), test.htmlFlashError)
		})
	}

	// test for db users
	truncateUsers()
}

func TestUpdateUser_WithInputPOSTForSuccess(t *testing.T) {
	// assert
	assert := assert.New(t)

	noAuth := setupTestServer(t)
	testCases := []struct {
		name             string
		path             string // id=string. Exemple, id="1"
		form             types.UserForm
		htmlFlashSuccess string
		token            echo.Map
	}{
		{
			name: "user admin [admin] for POST update to success",
			token: echo.Map{
				"username": "admin",
				"role":     "admin",
			},
			path: "1", // admin: 1 admin
			form: types.UserForm{
				Role:     "admin",
				Username: "admin-success",
			},
			htmlFlashSuccess: `<strong>success:</strong> success update user: admin-success!`,
		},
		{
			name: "user subali [user] for POST update to success",
			token: echo.Map{
				"username": "subali",
				"role":     "user",
			},
			path: "3", // user: 3 subali
			form: types.UserForm{
				Username: "subali-success",
				Name:     "Subali Success",
			},
			htmlFlashSuccess: `<strong>success:</strong> success update user: subali-success!`,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			token, _ := auth.GenerateToken(conf.AppJWTAuthSign, 1, "admin", "admin")
			result := noAuth.POST("/users/view/{id}").
				WithPath("id", test.path).
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
				// HTTP response status: 200 OK
				Status(http.StatusOK)

			assert.Contains(result.Body().Raw(), test.htmlFlashSuccess)
		})
	}

	// test for db users
	truncateUsers()
}

func TestUpdateUser_WithInputPOSTForFailure(t *testing.T) {
	// assert
	assert := assert.New(t)

	noAuth := setupTestServer(t)
	testCases := []struct {
		name             string
		token            echo.Map
		path             string // id=string. Exemple, id="1"
		form             types.UserForm
		jsonMessageError string
		status           int
	}{
		{
			name: "user admin [admin] for POST update to role error failure id=1",
			token: echo.Map{
				"username": "admin",
				"role":     "admin",
			},
			path: "1", // admin: 1 admin
			form: types.UserForm{
				Role:     "failure", // -> look
				Username: "admin-error",
				Location: 0,
			},
			jsonMessageError: `{"message":"Internal Server Error"}`,
			status:           http.StatusInternalServerError,
		},
		{
			name: "user subali [user] for POST update it failure: id=-1",
			token: echo.Map{
				"username": "subali",
				"role":     "user",
			},
			path: "-1", // -> look
			form: types.UserForm{
				Location: 0,
			},
			jsonMessageError: `{"message":"User Not Found"}`,
			status:           http.StatusNotFound,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			token, _ := auth.GenerateToken(
				conf.AppJWTAuthSign,
				0,
				test.token["username"].(string),
				test.token["role"].(string),
			)
			result := noAuth.POST("/users/view/{id}").
				WithPath("id", test.path).
				WithCookie("token", token).
				WithForm(test.form).
				Expect().
				// HTTP response status: 200 OK
				Status(test.status)

			resultBody := result.Body().Raw()
			assert.Contains(resultBody, test.jsonMessageError)
		})
	}

	// test for db users
	truncateUsers()
}

func TestUpdateUserByPassword_WithInputPOSTForSuccess(t *testing.T) {
	// assert
	assert := assert.New(t)

	noAuth := setupTestServer(t)
	testCases := []struct {
		name             string
		path             string // id=string. Exemple, id="1"
		form             types.NewPasswordForm
		htmlFlashSuccess string
		token            echo.Map
	}{
		{
			name: "user admin [admin] for POST update by password to success",
			token: echo.Map{
				"username": "admin",
				"role":     "admin",
			},
			path: "1", // admin: 1 admin
			form: types.NewPasswordForm{
				OldPassword:        "admin123",
				NewPassword:        "admin-success",
				ConfirmNewPassword: "admin-success",
			},
			htmlFlashSuccess: "<strong>success:</strong> success update user by password: admin!",
		},
		{
			name: "user subali [user] for POST update by password to success",
			token: echo.Map{
				"username": "subali",
				"role":     "user",
			},
			path: "3", // user: 3 subali
			form: types.NewPasswordForm{
				OldPassword:        "user123",
				NewPassword:        "user-success",
				ConfirmNewPassword: "user-success",
			},
			htmlFlashSuccess: `<strong>success:</strong> success update user by password: subali!`,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			token, _ := auth.GenerateToken(conf.AppJWTAuthSign, 1, "admin", "admin")
			result := noAuth.POST("/users/view/{id}/password").
				WithPath("id", test.path).
				WithCookie("token", token).
				WithForm(types.NewPasswordForm{
					OldPassword:        test.form.OldPassword,
					NewPassword:        test.form.NewPassword,
					ConfirmNewPassword: test.form.ConfirmNewPassword,
				}).
				Expect().
				// HTTP response status: 200 OK
				Status(http.StatusOK)

			assert.Contains(result.Body().Raw(), test.htmlFlashSuccess)
		})
	}

	// test for db users
	truncateUsers()
}

func TestUpdateUserByPassword_WithInputPOSTForFailure(t *testing.T) {
	// assert
	assert := assert.New(t)

	noAuth := setupTestServer(t)
	testCases := []struct {
		name             string
		token            echo.Map
		path             string // id=string. Exemple, id="1"
		form             types.UserForm
		jsonMessageError string
		status           int
	}{
		{
			name: "user admin [admin] for POST update to role error failure id=1",
			token: echo.Map{
				"username": "admin",
				"role":     "admin",
			},
			path: "1", // admin: 1 admin
			form: types.UserForm{
				Role:     "failure", // -> look
				Username: "admin-error",
				Location: 0,
			},
			jsonMessageError: `{"message":"Internal Server Error"}`,
			status:           http.StatusInternalServerError,
		},
		{
			name: "user subali [user] for POST update it failure: id=-1",
			token: echo.Map{
				"username": "subali",
				"role":     "user",
			},
			path: "-1", // -> look
			form: types.UserForm{
				Location: 0,
			},
			jsonMessageError: `{"message":"User Not Found"}`,
			status:           http.StatusNotFound,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			token, _ := auth.GenerateToken(
				conf.AppJWTAuthSign,
				0,
				test.token["username"].(string),
				test.token["role"].(string),
			)
			result := noAuth.POST("/users/view/{id}").
				WithPath("id", test.path).
				WithCookie("token", token).
				WithForm(test.form).
				Expect().
				// HTTP response status: 200 OK
				Status(test.status)

			resultBody := result.Body().Raw()
			assert.Contains(resultBody, test.jsonMessageError)
		})
	}

	// test for db users
	truncateUsers()
}

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
