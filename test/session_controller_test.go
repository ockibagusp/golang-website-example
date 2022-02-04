package test

import (
	"fmt"
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
	noAuth := setupTestServer(t)
	noAuthCSRF := setupTestServerNoAuthCSRF(noAuth)

	// test for db users
	truncateUsers(db)
	models.User{
		Username: "ockibagusp",
		Email:    "ocki.bagus.p@gmail.com",
		Password: "$2a$10$Y3UewQkjw808Ig90OPjuq.zFYIUGgFkWBuYiKzwLK8n3t9S8RYuYa",
		Name:     "Ocki Bagus Pratama",
	}.Save(db)

	testCases := []struct {
		method       int
		name         string
		user         types.LoginForm
		flashMessage bool
		// TODO: flash error for string
		flashErrorForString string
		status              int
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
			flashMessage:        true,
			flashErrorForString: "username or password not match",
			// HTTP response status: 403 Forbidden
			status: http.StatusForbidden,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.method == GET {
				noAuth.GET("/login").
					Expect().
					Status(tc.status)
				return
			}
			// tc.method == POST
			flashError := noAuthCSRF.POST("/login").
				WithForm(tc.user).
				WithFormField("X-CSRF-Token", csrfToken).
				Expect().
				Status(tc.status).
				Body().Raw()

			actual := fmt.Sprintf(`<p class="text-danger">*%s</p>`, tc.flashErrorForString)

			regex := regexp.MustCompile(`<p class\="text-danger">\*(.*)</p>`)
			match := regex.FindString(flashError)

			// flash message: "username or password not match"
			if tc.flashMessage {
				assert.Equal(t, match, actual)
			}
		})
	}
}
