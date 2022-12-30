package tests

import (
	"net/http"
	"regexp"
	"testing"

	"github.com/ockibagusp/golang-website-example/models"
	"github.com/ockibagusp/golang-website-example/types"
	"github.com/stretchr/testify/assert"
)

const GET int = 1

// POST int = 2
const POST = 2

func TestLogin(t *testing.T) {
	no_auth := setupTestServer(t)

	// test for db users
	truncateUsers(db)
	models.User{
		Username: "ockibagusp",
		Email:    "ocki.bagus.p@gmail.com",
		Password: "$2a$10$Y3UewQkjw808Ig90OPjuq.zFYIUGgFkWBuYiKzwLK8n3t9S8RYuYa",
		Name:     "Ocki Bagus Pratama",
	}.Save(db)

	test_cases := []struct {
		method int
		name   string
		user   types.LoginForm
		flash  regex
		status int
	}{
		{
			method: GET,
			name:   "login get",
			// HTTP response status: 200 OK
			status: http.StatusOK,
		},
		{
			method: POST,
			name:   "login success",
			user: types.LoginForm{
				Username: "ockibagusp",
				Password: "user123",
			},
			// HTTP response status: 200 OK
			status: http.StatusOK,
		},
		{
			method: POST,
			name:   "login failure",
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
			if test.method == GET {
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
