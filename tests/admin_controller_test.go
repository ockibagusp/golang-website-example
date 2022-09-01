package tests

import (
	"net/http"
	"testing"

	"github.com/gavv/httpexpect/v2"
)

func TestAdminDeletePermanently(t *testing.T) {
	no_auth := setupTestServer(t)
	auth_admin := setupTestServerAuth(no_auth, 1)
	// auth_sugriwa := setupTestServerAuth(no_auth, 2)

	// test for db users
	truncateUsers(db)

	test_cases := []struct {
		name      string
		expect    *httpexpect.Expect // auth or no-auth
		url_query string
		status    int
	}{
		/*
			delete permanently [admin]
		*/
		{
			name:   "delete permanently [admin] to GET it success: all",
			expect: auth_admin,
			// HTTP response status: 200 OK
			status: http.StatusOK,
		},

		// /*
		// 	delete permanently [sugriwa]
		// */
		// {
		// 	name:   "delete permanently [sugriwa] to GET it success: all",
		// 	expect: auth_sugriwa,
		// 	// HTTP response status: 404 Not Found
		// 	status: http.StatusNotFound,
		// },

		// /*
		// 	No Auth
		// */
		// {
		// 	name:   "delete permanently [no-auth] to GET it failure",
		// 	expect: no_auth,
		// 	// HTTP response status: 404 Not Found
		// 	status: http.StatusNotFound,
		// },
	}

	for _, test := range test_cases {
		var result *httpexpect.Response
		/*
			expect := test.expect

			or,

			var expect = test.expect
		*/
		var expect *httpexpect.Expect = test.expect

		t.Run(test.name, func(t *testing.T) {
			// @route: exemple "/admin/delete-permanently?admin=all"
			if test.url_query != "" {
				result = expect.GET("/admin/delete-permanently").
					WithQuery(test.url_query, "all").
					Expect().
					Status(test.status)
			} else {
				// @route: "/admin/delete-permanently"
				result = expect.GET("/admin/delete-permanently").
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
