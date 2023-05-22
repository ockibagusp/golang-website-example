package controller_test

import (
	"net/http"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/labstack/echo/v4"
	"github.com/ockibagusp/golang-website-example/business"
	"github.com/ockibagusp/golang-website-example/business/auth"
	"github.com/stretchr/testify/assert"
)

func TestAdminRestoreByID_WithInputGETForSuccess(t *testing.T) {
	assert := assert.New(t)
	noAuth := setupTestServer(t)

	// delete user sugriwa (uid=2)
	err := newUserService(db).Delete(business.InternalContext{}, 2)
	if err != nil {
		panic("sugriwa: username not already: " + err.Error())
	}

	t.Run("delete restore by uid [admin] to GET it success: uid=2", func(t *testing.T) {
		token, _ := auth.GenerateToken(conf.AppJWTAuthSign, 1, "admin", "admin")
		result := noAuth.GET("/admin/restore/{uid}", "2").
			WithCookie("token", token).
			Expect().
			// HTTP response status: 200 OK
			Status(http.StatusOK)

		assert.Contains(result.Body().Raw(), "<strong>success:</strong> success restore user: sugriwa!")
	})

	// test for db users
	truncateUsers()
}

func TestAdminRestoreByID_WithInputGETForFailure(t *testing.T) {
	noAuth := setupTestServer(t)

	// delete user sugriwa (uid=2)
	err := newUserService(db).Delete(business.InternalContext{}, 2)
	if err != nil {
		panic("sugriwa: username not already: " + err.Error())
	}

	testCases := []struct {
		name             string
		token            echo.Map
		path             string
		jsonMessageError string
		status           int
	}{
		/*
			delete restore by uid [admin]
		*/
		{
			name: "delete restore by uid [admin] to GET it failure: uid=1",
			token: echo.Map{
				"username": "admin",
				"role":     "admin",
			},
			path:             "1",
			jsonMessageError: `{"message":"Forbidden"}`,
			// HTTP response status: 403 Forbidden
			status: http.StatusForbidden,
		},
		{
			name: "delete restore by uid [admin] to GET it failure: uid=-1",
			token: echo.Map{
				"username": "admin",
				"role":     "admin",
			},
			path:             "-1",
			jsonMessageError: `{"message":"User Not Found"}`,
			// HTTP response status: 404 Not Found
			status: http.StatusNotFound,
		},

		/*
			delete restore by uid [sugriwa]
		*/
		{
			name: "delete restore by uid [sugriwa] to GET it success: uid=1",
			token: echo.Map{
				"username": "sugriwa",
				"role":     "user",
			},
			path:             "1",
			jsonMessageError: `{"message":"Not Found"}`,
			// HTTP response status: 404 Not Found
			status: http.StatusNotFound,
		},

		/*
			delete restore by uid anonymous
		*/
		{
			name: "delete restore by uid [no-auth] to GET it failure: uid=1",
			token: echo.Map{
				"username": "anonymous",
				"role":     "anonymous",
			},
			path:             "1",
			jsonMessageError: `{"message":"Not Found"}`,
			// HTTP response status: 404 Not Found
			status: http.StatusNotFound,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			token, _ := auth.GenerateToken(
				conf.AppJWTAuthSign,
				0,
				test.token["username"].(string),
				test.token["role"].(string),
			)
			noAuth.GET("/admin/restore/{uid}", test.path).
				WithCookie("token", token).
				Expect().
				Status(test.status)
		})
	}

	// test for db users
	truncateUsers()
}

func TestAdminDeletePermanentlyByID_WithInputGETForSuccess(t *testing.T) {
	assert := assert.New(t)
	noAuth := setupTestServer(t)

	// delete user sugriwa (uid=2)
	err := newUserService(db).Delete(business.InternalContext{}, 2)
	if err != nil {
		panic("subali: username not already: " + err.Error())
	}

	t.Run("delete permanently by uid [admin] to GET it success: uid=2", func(t *testing.T) {
		token, _ := auth.GenerateToken(conf.AppJWTAuthSign, 1, "admin", "admin")
		result := noAuth.GET("/admin/delete/permanently/{uid}", "2").
			WithCookie("token", token).
			Expect().
			// HTTP response status: 200 OK
			Status(http.StatusOK)

		assert.Contains(result.Body().Raw(), "<strong>success:</strong> success permanently user: sugriwa!")
	})

	// test for db users
	truncateUsers()
}

func TestAdminDeletePermanentlyByID_WithInputGETForFailure(t *testing.T) {
	assert := assert.New(t)
	noAuth := setupTestServer(t)

	// delete user sugriwa (uid=2)
	err := newUserService(db).Delete(business.InternalContext{}, 2)
	if err != nil {
		panic("subali: username not already: " + err.Error())
	}

	testCases := []struct {
		name             string
		token            echo.Map
		urlQuery         string
		jsonMessageError string
		status           int
	}{
		/*
			delete permanently [sugriwa]
		*/
		{
			name: "delete permanently test [user] to GET it failure: all",
			token: echo.Map{
				"username": "test",
				"role":     "user",
			},
			jsonMessageError: `{"message":"Not Found"}`,
			// HTTP response status: 404 Not Found
			status: http.StatusNotFound,
		},

		/*
			anonymous
		*/
		{
			name: "delete permanently anonymous to GET it failure",
			token: echo.Map{
				"username": "anonymous",
				"role":     "anonymous",
			},
			jsonMessageError: `{"message":"Not Found"}`,
			// HTTP response status: 404 Not Found
			status: http.StatusNotFound,
		},
	}

	for _, test := range testCases {
		var result *httpexpect.Response
		token, _ := auth.GenerateToken(
			conf.AppJWTAuthSign,
			0,
			test.token["username"].(string),
			test.token["role"].(string),
		)
		auth := noAuth.Builder(func(req *httpexpect.Request) {
			req.WithCookie("token", token)
		})

		t.Run(test.name, func(t *testing.T) {
			result = auth.GET("/admin/delete/permanently/{uid}", test.urlQuery).
				Expect().
				Status(test.status)

			resultBody := result.Body().Raw()

			assert.Contains(resultBody, test.jsonMessageError)
		})
	}

	// test for db users
	truncateUsers()
}
