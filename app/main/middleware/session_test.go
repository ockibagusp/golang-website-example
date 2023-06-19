package middleware_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	. "golang-website-example/app/main/middleware"
)

func TestSessionFlashSuccessSuccess(t *testing.T) {
	// assert
	assert := assert.New(t)

	// echo setup
	e := echo.New()

	request := httptest.NewRequest(http.MethodGet, "/user/add", nil)
	recorder := httptest.NewRecorder()
	c := e.NewContext(request, recorder)

	SetFlashSuccess(c, fmt.Sprintf("success new user: %s!", "sugriwa"))
	// test data
	expected := GetFlashSuccess(c)

	assert.Equal(expected, []string([]string{"success new user: sugriwa!"}))
}

func TestSessionFlashErrorSuccess(t *testing.T) {
	// assert
	assert := assert.New(t)

	// echo setup
	e := echo.New()

	request := httptest.NewRequest(http.MethodGet, "/login", nil)
	recorder := httptest.NewRecorder()
	c := e.NewContext(request, recorder)

	SetFlash(c, "error", "username or password not match")

	// test data
	expected := GetFlashError(c)

	assert.Equal(expected, []string([]string{"username or password not match"}))
}
