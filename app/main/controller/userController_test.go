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
			token, _ := auth.GenerateToken(
				conf.AppJWTAuthSign,
				0,
				test.token["username"].(string),
				test.token["role"].(string),
			)
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
		token            echo.Map
		path             string // id=string. Exemple, id="1"
		form             types.UserForm
		htmlFlashSuccess string
	}{
		{
			name: "user admin [admin] for POST update to success: uid=1",
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
			name: "user subali [user] for POST update to success: uid=3",
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
			token, _ := auth.GenerateToken(
				conf.AppJWTAuthSign,
				0,
				test.token["username"].(string),
				test.token["role"].(string),
			)
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
			name: "user admin [admin] for POST update to role error failure: uid=1",
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
			name: "user subali [user] for POST update it failure: uid=-1",
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
		token            echo.Map
		path             string // id=string. Exemple, id="1"
		form             types.NewPasswordForm
		htmlFlashSuccess string
	}{
		{
			name: "user admin [admin] for POST update by password to success: uid=1",
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
			name: "user subali [user] for POST update by password to success: uid=3",
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
			token, _ := auth.GenerateToken(
				conf.AppJWTAuthSign,
				0,
				test.token["username"].(string),
				test.token["role"].(string),
			)
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
	noAuth := setupTestServer(t)
	testCases := []struct {
		name  string
		token echo.Map
		path  string
		form  types.NewPasswordForm
	}{
		{
			name: "user admin [admin] for POST update by password to old password error failure: uid=1",
			token: echo.Map{
				"username": "admin",
				"role":     "admin",
			},
			path: "1", // admin: 1 admin
			form: types.NewPasswordForm{
				OldPassword:        "admin-failure", // -> look
				NewPassword:        "admin",
				ConfirmNewPassword: "admin",
			},
		},
		{
			name: "user sugriwa [user] for POST update by password to old password error failure: uid=2",
			token: echo.Map{
				"username": "sugriwa",
				"role":     "user",
			},
			path: "2",
			form: types.NewPasswordForm{
				OldPassword:        "user-failure",
				NewPassword:        "user",
				ConfirmNewPassword: "user",
			},
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
			noAuth.POST("/users/view/{id}/password").
				WithPath("id", test.path).
				WithCookie("token", token).
				WithForm(types.NewPasswordForm{
					OldPassword:        test.form.OldPassword,
					NewPassword:        test.form.NewPassword,
					ConfirmNewPassword: test.form.ConfirmNewPassword,
				}).
				Expect().
				// HTTP response status: 403 Forbidden
				Status(http.StatusForbidden)
		})
	}

	// test for db users
	truncateUsers()
}

func TestDelete_WithInputPOSTForSuccess(t *testing.T) {
	// assert
	assert := assert.New(t)

	noAuth := setupTestServer(t)
	testCases := []struct {
		name             string
		token            echo.Map
		path             string // id=string. Exemple, id="1"
		htmlFlashSuccess string
	}{
		/*
			delete admin [admin] to success
		*/
		{
			name: "user admin [admin] for GET delete sugriwa [user] to success: uid=2",
			token: echo.Map{
				"username": "admin",
				"role":     "admin",
			},
			path:             "2",
			htmlFlashSuccess: "<strong>success:</strong> success delete user: sugriwa!",
		},
		/*
			delete subali [user] to success
		*/
		{
			name: "user subali [user] for GET delete to success: uid=3",
			token: echo.Map{
				"username": "subali",
				"role":     "user",
			},
			path:             "3",
			htmlFlashSuccess: "<strong>success:</strong> success delete user: subali!",
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
			result := noAuth.GET("/users/delete/{id}", test.path).
				WithCookie("token", token).
				Expect().
				// HTTP response status: 200 OK
				Status(http.StatusOK)

			assert.Contains(result.Body().Raw(), test.htmlFlashSuccess)
		})
	}

	// test for db users
	truncateUsers()
}

func TestDelete_WithInputPOSTForFailure(t *testing.T) {
	// assert
	assert := assert.New(t)

	noAuth := setupTestServer(t)
	testCases := []struct {
		name             string
		token            echo.Map
		path             string // id=string. Exemple, id="1"
		jsonMessageError string
		htmlFlashError   string
		status           int
	}{
		/*
			delete admin [admin] to admin [admin] failure: uid=1
		*/
		{
			name: "user admin [admin] for GET delete admin [admin] to failure: uid=1",
			token: echo.Map{
				"username": "admin",
				"role":     "admin",
			},
			path:             "1",
			jsonMessageError: `{"message":"Forbidden"}`,
			status:           http.StatusForbidden,
		},
		/*
			delete subali [user] to failure: uid=2
		*/
		{
			name: "user subali [user] for GET delete failure: uid=-1",
			token: echo.Map{
				"username": "subali",
				"role":     "user",
			},
			path:             "-1",
			jsonMessageError: `{"message":"User Not Found"}`,
			status:           http.StatusNotFound,
		},
		/*
			delete anonymous to failure: uid=2
		*/
		{
			name: "user anonymous for GET delete failure: uid=2",
			token: echo.Map{
				"username": "anonymous",
				"role":     "anonymous",
			},
			path:           "2",
			htmlFlashError: "*login process failed!",
			status:         http.StatusOK,
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
			result := noAuth.GET("/users/delete/{id}", test.path).
				WithCookie("token", token).
				Expect().
				// HTTP response status: 200 OK
				Status(test.status)

			resultBody := result.Body().Raw()
			assert.Contains(resultBody, test.htmlFlashError)
			assert.Contains(resultBody, test.jsonMessageError)
		})
	}

	// test for db users
	truncateUsers()
}
