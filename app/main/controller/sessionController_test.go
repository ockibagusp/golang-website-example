package controller_test

import (
	"net/http"
	"regexp"
	"testing"

	"github.com/ockibagusp/golang-website-example/app/main/types"
	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	no_auth := setupTestServer(t)

	// test for db users
	truncateUsers()

	test_cases := []struct {
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
			method: HTTP_REQUEST_GET,
			expect: ADMIN,
			// HTTP response status: 200 OK
			status: http.StatusOK,
		},
		{
			name:   "users [admin] to POST login success",
			method: HTTP_REQUEST_POST,
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
			method: HTTP_REQUEST_POST,
			expect: ADMIN,
			user: types.LoginForm{
				Username: "admin",
				Password: "<bad password>",
			},
			flash: regex{
				must_compile: `<p class="text-danger">*(.*)</p>`,
				actual:       `<p class="text-danger">*username or password not match</p>`,
			},
			// HTTP response status: 403 Forbidden
			status: http.StatusForbidden,
		},

		/*
			users [ockibagusp]
		*/
		{
			name:   "users [ockibagusp] to GET login",
			method: HTTP_REQUEST_GET,
			expect: OCKIBAGUSP,
			// HTTP response status: 200 OK
			status: http.StatusOK,
		},
		{
			name:   "users [ockibagusp] to POST login success",
			method: HTTP_REQUEST_POST,
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
			method: HTTP_REQUEST_POST,
			expect: OCKIBAGUSP,
			user: types.LoginForm{
				Username: "ockibagusp",
				Password: "<bad password>",
			},
			flash: regex{
				must_compile: `<p class="text-danger">*(.*)</p>`,
				actual:       `<p class="text-danger">*username or password not match</p>`,
			},
			// HTTP response status: 403 Forbidden
			status: http.StatusForbidden,
		},
	}

	for _, test := range test_cases {
		t.Run(test.name, func(t *testing.T) {
			userSelectTest = test.expect // ADMIN and OCKIBAGUSP

			if test.method == HTTP_REQUEST_GET {
				no_auth.GET("/login").
					Expect().
					Status(test.status)
				return
			}
			// tc.method == POST
			result := no_auth.POST("/login").
				WithForm(test.user).
				Expect().
				Status(test.status)

			// flash message: "username or password not match"
			if (test.flash.must_compile == "") && (test.flash.actual == "") {
				result_body := result.Body().Raw()

				regex := regexp.MustCompile(test.flash.must_compile)
				match := regex.FindString(result_body)

				assert.Equal(t, match, test.flash.actual)
			}
		})
	}
}
