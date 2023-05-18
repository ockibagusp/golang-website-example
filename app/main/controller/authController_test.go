package controller_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/ockibagusp/golang-website-example/app/main/types"
	"github.com/stretchr/testify/assert"
)

func TestLogin_WithInputPOSTForSuccess(t *testing.T) {
	// assert
	assert := assert.New(t)

	// echo setup
	e := echo.New()

	testCases := []struct {
		name   string
		user   types.LoginForm
		status int
	}{
		/*
			users admin [admin]
		*/
		{
			name: "user admin [admin] for POST login to success",
			user: types.LoginForm{
				Username: "admin",
				Password: "admin123",
			},
			// HTTP response status: 200 OK
			status: http.StatusOK,
		},
		/*
			users sugriwa [user]
		*/
		{
			name: "user sugriwa [user] for POST login to success",
			user: types.LoginForm{
				Username: "sugriwa",
				Password: "user123",
			},
			// HTTP response status: 200 OK
			status: http.StatusOK,
		},
		/*
			users subali [user]
		*/
		{
			name: "user subali [user] for POST login to success",
			user: types.LoginForm{
				Username: "subali",
				Password: "user123",
			},
			// HTTP response status: 200 OK
			status: http.StatusOK,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// test data
			expected := make(url.Values)
			expected.Set("username", test.user.Username)
			expected.Set("password", test.user.Password)

			request := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(expected.Encode()))
			request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
			recorder := httptest.NewRecorder()
			c := e.NewContext(request, recorder)

			// internal setup
			mockApp := setupTestController()

			// act
			statusCode := recorder.Code
			if assert.NoError(mockApp.Login(c)) {
				assert.Equalf(test.status, statusCode, "got: %d but expect %d", test.status, statusCode)
			}
		})
	}
}

func TestLogin_WithInputPOSTNotForUsernameFailure(t *testing.T) {
	// assert
	assert := assert.New(t)

	// echo setup
	e := echo.New()

	testCases := []struct {
		name   string
		user   types.LoginForm
		status int
	}{
		/*
			user admin_failure
		*/
		{
			name: "user admin_failure for POST login to failure",
			user: types.LoginForm{
				Username: "admin_failure",
				Password: "admin123",
			},
			// HTTP response status: 406 Not Acceptable?

			// HTTP response status: 200 OK
			status: http.StatusOK,
		},
		/*
			user sugriwa_failure
		*/
		{
			name: "user sugriwa_failure for POST login to failure",
			user: types.LoginForm{
				Username: "sugriwa_failure",
				Password: "user123",
			},
			// HTTP response status: 200 OK
			status: http.StatusOK,
		},
		/*
			user anonymous
		*/
		{
			name: "user anonymous for POST login to failure",
			user: types.LoginForm{
				Username: "anonymous",
				Password: "user123",
			},
			// HTTP response status: 200 OK
			status: http.StatusOK,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// test data
			expected := make(url.Values)
			expected.Set("username", test.user.Username)
			expected.Set("password", test.user.Password)

			request := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(expected.Encode()))
			request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
			recorder := httptest.NewRecorder()
			c := e.NewContext(request, recorder)

			// internal setup
			mockApp := setupTestController()

			// act
			statusCode := recorder.Code
			if assert.Error(mockApp.Login(c)) {
				assert.Equalf(test.status, statusCode, "got: %d but expect %d", test.status, statusCode)
				// ? assert.Contains(recorder.Body.String(), "*username or password not match")
			}
		})
	}
}

func TestLogin_WithInputPOSTPasswordTooShortFailure(t *testing.T) {
	// assert
	assert := assert.New(t)

	// echo setup
	e := echo.New()

	testCases := []struct {
		name   string
		user   types.LoginForm
		status int
	}{
		/*
			users admin [admin]
		*/
		{
			name: "user admin [admin] for POST login to failure",
			user: types.LoginForm{
				Username: "admin",
				Password: "admi",
			},
			// HTTP response status: 406 Not Acceptable?

			// HTTP response status: 200 OK
			status: http.StatusOK,
		},
		/*
			users sugriwa [user]
		*/
		{
			name: "user sugriwa [user] for POST login to failure",
			user: types.LoginForm{
				Username: "sugriwa",
				Password: "user",
			},
			// HTTP response status: 200 OK
			status: http.StatusOK,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// test data
			expected := make(url.Values)
			expected.Set("username", test.user.Username)
			expected.Set("password", test.user.Password)

			request := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(expected.Encode()))
			request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
			recorder := httptest.NewRecorder()
			c := e.NewContext(request, recorder)

			// internal setup
			mockApp := setupTestController()

			// act
			statusCode := recorder.Code
			if assert.Error(mockApp.Login(c)) {
				assert.Equalf(test.status, statusCode, "got: %d but expect %d", test.status, statusCode)
				// ? assert.Contains(recorder.Body.String(), "*username or password not match")
			}
		})
	}
}

func TestLogin_WithInputPOSTPasswordTooLongFailure(t *testing.T) {
	// assert
	assert := assert.New(t)

	// echo setup
	e := echo.New()

	testCases := []struct {
		name   string
		user   types.LoginForm
		status int
	}{
		/*
			users admin [admin]
		*/
		{
			name: "user admin [admin] for POST login to failure",
			user: types.LoginForm{
				Username: "admin",
				Password: "<bad password too long>",
			},
			// HTTP response status: 406 Not Acceptable?

			// HTTP response status: 200 OK
			status: http.StatusOK,
		},
		/*
			users sugriwa [user]
		*/
		{
			name: "user sugriwa [user] for POST login to failure",
			user: types.LoginForm{
				Username: "sugriwa",
				Password: "<bad password too long>",
			},
			// HTTP response status: 200 OK
			status: http.StatusOK,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// test data
			expected := make(url.Values)
			expected.Set("username", test.user.Username)
			expected.Set("password", test.user.Password)

			request := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(expected.Encode()))
			request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
			recorder := httptest.NewRecorder()
			c := e.NewContext(request, recorder)

			// internal setup
			mockApp := setupTestController()

			// act
			statusCode := recorder.Code
			if assert.Error(mockApp.Login(c)) {
				assert.Equalf(test.status, statusCode, "got: %d but expect %d", test.status, statusCode)
				// ? assert.Contains(recorder.Body.String(), "*username or password not match")
			}
		})
	}
}

func TestLogin_WithInputPOSTPasswordTooShort(t *testing.T) {
	// assert
	assert := assert.New(t)

	// echo setup
	e := echo.New()

	testCases := []struct {
		name   string
		user   types.LoginForm
		status int
	}{
		/*
			users admin [admin]
		*/
		{
			name: "user admin [admin] for POST login to success",
			user: types.LoginForm{
				Username: "admin",
				Password: "<bad>",
			},
			// HTTP response status: 200 OK
			status: http.StatusOK,
		},
		/*
			users sugriwa [user]
		*/
		{
			name: "user sugriwa [user] for POST login to success",
			user: types.LoginForm{
				Username: "sugriwa",
				Password: "<bad>",
			},
			// HTTP response status: 200 OK
			status: http.StatusOK,
		},
		/*
			users subali [user]
		*/
		{
			name: "user subali [user] for POST login to success",
			user: types.LoginForm{
				Username: "subali",
				Password: "user",
			},
			// HTTP response status: 200 OK
			status: http.StatusOK,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			// test data
			expected := make(url.Values)
			expected.Set("username", test.user.Username)
			expected.Set("password", test.user.Password)

			request := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(expected.Encode()))
			request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
			recorder := httptest.NewRecorder()
			c := e.NewContext(request, recorder)

			// internal setup
			mockApp := setupTestController()

			// act
			statusCode := recorder.Code
			if assert.Error(mockApp.Login(c)) {
				assert.Equalf(test.status, statusCode, "got: %d but expect %d", test.status, statusCode)
			}
		})
	}
}
