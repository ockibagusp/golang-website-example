package test

import (
	"net/http"
	"testing"
)

func TestAboutSuccess(t *testing.T) {
	noAuth := setupTestServer(t)

	t.Run("about success", func(t *testing.T) {
		auth := setupTestServerAuth(noAuth, 0)
		auth.GET("/about").
			Expect().
			Status(http.StatusOK)
	})

	t.Run("about failure", func(t *testing.T) {
		noAuth.GET("/about").
			Expect().
			Status(http.StatusOK)
	})
}
