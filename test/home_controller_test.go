package test

import (
	"net/http"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHomeController(t *testing.T) {
	assert := assert.New(t)

	no_auth := setupTestServer(t)
	t.Run("home [no-auth] success", func(t *testing.T) {
		result := no_auth.GET("/").
			Expect().
			Status(http.StatusOK)

		no_auth_text := result.Body().Raw()

		// TODO: why?

		// navbar nav
		regex := regexp.MustCompile(`<a href="/login" class="btn btn-outline-success my-2 my-sm-0">(.*)</a>`)
		match := regex.FindString(no_auth_text)

		actual := `<a href="/login" class="btn btn-outline-success my-2 my-sm-0">Login</a>`
		assert.Equal(match, actual)

		// main: jumbotron
		regex = regexp.MustCompile(`<p class="lead">(.*)</p>`)
		match = regex.FindString(no_auth_text)

		assert.Equal(match, `<p class="lead">Test.</p>`)
	})

	auth_admin := setupTestServerAuth(no_auth, 1)
	t.Run("home [admin] success", func(t *testing.T) {
		result := auth_admin.GET("/").
			Expect().
			Status(http.StatusOK)

		admin_text := result.Body().Raw()

		// navbar nav
		regex := regexp.MustCompile(`<a class="btn">(.*)</a>`)
		match := regex.FindString(admin_text)

		assert.Equal(match, `<a class="btn">ADMIN</a>`)

		// main: jumbotron
		regex = regexp.MustCompile(`<p class="lead">(.*)</p>`)
		match = regex.FindString(admin_text)

		assert.Equal(match, `<p class="lead">Admin.</p>`)
	})

	auth_user := setupTestServerAuth(no_auth, 0)
	t.Run("home [user] success", func(t *testing.T) {
		result := auth_user.GET("/").
			Expect().
			Status(http.StatusOK)

		user_text := result.Body().Raw()

		// navbar nav
		regex := regexp.MustCompile(`<p class="lead">(.*)</p>`)
		match := regex.FindString(user_text)

		assert.Equal(match, `<p class="lead">User.</p>`)

		// main: jumbotron
		regex = regexp.MustCompile(`<p class="lead">(.*)</p>`)
		match = regex.FindString(user_text)

		assert.Equal(match, `<p class="lead">User.</p>`)
	})
}
