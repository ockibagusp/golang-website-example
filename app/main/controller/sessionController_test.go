package controller_test

import (
	"net/http"
	"regexp"
	"testing"

	methodTest "github.com/ockibagusp/golang-website-example/app/main/controller/mock/method"
	modelsTest "github.com/ockibagusp/golang-website-example/app/main/controller/mock/models"
	"github.com/ockibagusp/golang-website-example/app/main/types"
	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
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

				assert.Equal(t, match, test.flash.actual)
			}
		})
	}
}
