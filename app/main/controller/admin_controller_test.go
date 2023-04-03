package controller_test

import (
	"net/http"
	"testing"

	"github.com/gavv/httpexpect/v2"
	methodTest "github.com/ockibagusp/golang-website-example/app/main/controller/mock/method"
	modelsTest "github.com/ockibagusp/golang-website-example/app/main/controller/mock/models"
)

func TestAdminDeletePermanently(t *testing.T) {
	noAuth := setupTestServer(t)

	// test for SetSession = false
	methodTest.SetSession = false
	// test for db users
	truncateUsers()

	testCases := []struct {
		name     string
		expect   string // auth or no-auth
		urlQuery string
		status   int
	}{
		/*
			delete permanently [admin]
		*/
		{
			name:   "delete permanently [admin] to GET it success: all",
			expect: ADMIN,
			// HTTP response status: 200 OK
			status: http.StatusOK,
		},

		/*
			delete permanently [sugriwa]
		*/
		{
			name:   "delete permanently [sugriwa] to GET it success: all",
			expect: SUGRIWA,
			// HTTP response status: 404 Not Found
			status: http.StatusNotFound,
		},

		/*
			No Auth
		*/
		{
			name:   "delete permanently [no-auth] to GET it failure",
			expect: ANONYMOUS,
			// HTTP response status: 404 Not Found
			status: http.StatusNotFound,
		},
	}

	for _, test := range testCases {
		var result *httpexpect.Response
		/*
			expect := test.expect

			or,

			var expect = test.expect
		*/
		modelsTest.UserSelectTest = test.expect // ADMIN and SUGRIWA

		t.Run(test.name, func(t *testing.T) {
			// @route: exemple "/admin/delete-permanently?admin=all"
			if test.urlQuery != "" {
				result = noAuth.GET("/admin/delete-permanently").
					WithQuery(test.urlQuery, "all").
					Expect().
					Status(test.status)
			} else {
				// @route: "/admin/delete-permanently"
				result = noAuth.GET("/admin/delete-permanently").
					Expect().
					Status(test.status)
			}

			statusCode := result.Raw().StatusCode
			if test.status != statusCode {
				t.Logf(
					"got: %d but expect %d", test.status, statusCode,
				)
				t.Fail()
			}
		})
	}
}
