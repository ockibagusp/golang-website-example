package controller_test

import (
	"net/http"
	"regexp"
	"testing"

	"github.com/gavv/httpexpect/v2"
	methodTest "github.com/ockibagusp/golang-website-example/app/main/controller/mock/method"
	modelsTest "github.com/ockibagusp/golang-website-example/app/main/controller/mock/models"
	"github.com/stretchr/testify/assert"
)

func TestAboutSuccess(t *testing.T) {
	assert := assert.New(t)

	noAuth := setupTestServer(t)

	// test for SetSession = false
	methodTest.SetSession = false
	// test for db users
	truncateUsers()

	testCases := []struct {
		name       string
		expect     string // admin, sugriwa
		htmlNavbar regex
	}{
		{
			name:   "about [no-auth] success",
			expect: ANONYMOUS,
			htmlNavbar: regex{
				mustCompile: `<a href="/login" class="btn btn-outline-success my-2 my-sm-0">(.*)</a>`,
				actual:      `<a href="/login" class="btn btn-outline-success my-2 my-sm-0">Login</a>`,
			},
		},
		{
			name:   "about [admin] success",
			expect: ADMIN,
			htmlNavbar: regex{
				mustCompile: `<a class="btn">(.*)</a>`,
				actual:      `<a class="btn">ADMIN</a>`,
			},
		},
		{
			name:   "home [user] success",
			expect: SUGRIWA,
			htmlNavbar: regex{
				mustCompile: `<a href="/users" class="btn btn-outline-secondary my-2 my-sm-0">(.*)</a>`,
				actual:      `<a href="/users" class="btn btn-outline-secondary my-2 my-sm-0">Users</a>`,
			},
		},
	}

	for _, test := range testCases {
		var result *httpexpect.Response
		modelsTest.UserSelectTest = test.expect // auth_admin, auth_sugriwa or no-auth

		t.Run(test.name, func(t *testing.T) {
			result = noAuth.GET("/about").
				Expect().
				Status(http.StatusOK)

			resultBody := result.Body().Raw()

			// navbar nav
			regex := regexp.MustCompile(test.htmlNavbar.mustCompile)
			match := regex.FindString(resultBody)

			assert.Equal(match, test.htmlNavbar.actual)
		})
	}
}
