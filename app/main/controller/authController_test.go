package controller_test

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	methodTest "github.com/ockibagusp/golang-website-example/app/main/controller/mock/method"
	modelsTest "github.com/ockibagusp/golang-website-example/app/main/controller/mock/models"
	"github.com/ockibagusp/golang-website-example/app/main/types"
	"github.com/stretchr/testify/assert"
)

func TestLogin_Success(t *testing.T) {
	// echo setup
	e := echo.New()

	// test data
	expected := "username=admin&password=admin123"

	request := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(expected))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	recorder := httptest.NewRecorder()
	c := e.NewContext(request, recorder)

	// internal setup
	mockApp := setupTestController()

	// act
	mockApp.Login(c)

	// assert
	assert := assert.New(t)
	assert.Equal(http.StatusOK, recorder.Code)
}

func TestRegisterAccount_PasswordTooShort(t *testing.T) {
	// test data
	// expected := `{"username":"user1","email":"123"}`

	// echo setup
	e := echo.New()
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	// req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recorder := httptest.NewRecorder()
	c := e.NewContext(request, recorder)

	// internal setup
	mockApp := setupTestController()

	// act
	mockApp.Home(c)

	// assert
	assert := assert.New(t)
	assert.Equal(http.StatusOK, recorder.Code)
	// assert.Contains(recorder.Body.String(), "A validation error occurred")
	// assert.Contains(recorder.Body.String(), "Password must be 8 characters")
}

func TestLogin(t *testing.T) {
	assert := assert.New(t)

	noAuth := setupTestServer(t)

	// test for db users
	truncateUsers()

	testCases := []struct {
		name   string
		method int
		expect string
		user   types.LoginForm
		flash  regex
		status int
	}{
		/*
			users [admin]
		*/
		{
			name:   "users [admin] to GET login",
			method: methodTest.HTTP_REQUEST_GET,
			expect: ADMIN,
			// HTTP response status: 200 OK
			status: http.StatusOK,
		},
		{
			name:   "users [admin] to POST login success",
			method: methodTest.HTTP_REQUEST_POST,
			expect: ADMIN,
			user: types.LoginForm{
				Username: "admin",
				Password: "admin123",
			},
			// HTTP response status: 200 OK
			status: http.StatusOK,
		},
		{
			name:   "users [admin] to POST login failure",
			method: methodTest.HTTP_REQUEST_POST,
			expect: ADMIN,
			user: types.LoginForm{
				Username: "admin",
				Password: "<bad password>",
			},
			flash: regex{
				mustCompile: `<p class="text-danger">*(.*)</p>`,
				actual:      `<p class="text-danger">*username or password not match</p>`,
			},
			// HTTP response status: 403 Forbidden
			status: http.StatusForbidden,
		},

		/*
			users [ockibagusp]
		*/
		{
			name:   "users [ockibagusp] to GET login",
			method: methodTest.HTTP_REQUEST_GET,
			expect: OCKIBAGUSP,
			// HTTP response status: 200 OK
			status: http.StatusOK,
		},
		{
			name:   "users [ockibagusp] to POST login success",
			method: methodTest.HTTP_REQUEST_POST,
			expect: OCKIBAGUSP,
			user: types.LoginForm{
				Username: "ockibagusp",
				Password: "user123",
			},
			// HTTP response status: 200 OK
			status: http.StatusOK,
		},
		{
			name:   "users [ockibagusp] to POST login failure",
			method: methodTest.HTTP_REQUEST_POST,
			expect: OCKIBAGUSP,
			user: types.LoginForm{
				Username: "ockibagusp",
				Password: "<bad password>",
			},
			flash: regex{
				mustCompile: `<p class="text-danger">*(.*)</p>`,
				actual:      `<p class="text-danger">*username or password not match</p>`,
			},
			// HTTP response status: 403 Forbidden
			status: http.StatusForbidden,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			modelsTest.UserSelectTest = test.expect // ADMIN and OCKIBAGUSP

			if test.method == methodTest.HTTP_REQUEST_GET {
				noAuth.GET("/login").
					Expect().
					Status(test.status)
				return
			}
			// tc.method == POST
			result := noAuth.POST("/login").
				WithForm(test.user).
				Expect().
				Status(test.status)

			// flash message: "username or password not match"
			if (test.flash.mustCompile == "") && (test.flash.actual == "") {
				resultBody := result.Body().Raw()

				regex := regexp.MustCompile(test.flash.mustCompile)
				match := regex.FindString(resultBody)

				assert.Equal(match, test.flash.actual)
			}
		})
	}
}
