package controller_test

import (
	"net/http"
	"regexp"
	"testing"

	"github.com/gavv/httpexpect/v2"
	methodTest "github.com/ockibagusp/golang-website-example/app/main/controller/mock/method"
	modelsTest "github.com/ockibagusp/golang-website-example/app/main/controller/mock/models"
	"github.com/ockibagusp/golang-website-example/business"
	"github.com/stretchr/testify/assert"
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

func TestAdminRestoreByID(t *testing.T) {
	noAuth := setupTestServer(t)

	// test for SetSession = false
	methodTest.SetSession = false
	// test for db users
	truncateUsers()

	// delete user sugriwa (id=2)
	err := newUserService(db).Delete(business.InternalContext{}, 2)
	if err != nil {
		panic("sugriwa: username not already: " + err.Error())
	}

	testCases := []struct {
		name             string
		expect           string // auth or no-auth
		path             string
		status           int
		htmlNavbar       regex
		htmlHeading      regex
		htmlFlashSuccess regex
		jsonMessageError regex
	}{
		/*
			delete restore by id [admin]
		*/
		{
			name:   "delete restore by id [admin] to GET it failure: id=1",
			expect: ADMIN,
			path:   "1",
			// HTTP response status: 403 Forbidden
			status: http.StatusForbidden,
			jsonMessageError: regex{
				mustCompile: `{"message":"(.*)"}`,
				actual:      `{"message":"Forbidden"}`,
			},
		},
		{
			name:   "delete restore by id [admin] to GET it success: id=2",
			expect: ADMIN,
			path:   "2",
			// HTTP response status: 200 OK
			status: http.StatusOK,
			// body navbar
			htmlNavbar: regex{
				mustCompile: `<a class="btn">(.*)</a>`,
				actual:      `<a class="btn">ADMIN</a>`,
			},
			// body heading
			htmlHeading: regex{
				mustCompile: `<h2 class="mt-4">(.*)</h2>`,
				actual:      `<h2 class="mt-4">Delete Permanently?</h2>`,
			},
			// flash message success
			htmlFlashSuccess: regex{
				mustCompile: `<strong>success:</strong> (.*)`,
				actual:      `<strong>success:</strong> success restore user: sugriwa!`,
			},
		},
		{
			name:   "delete restore by id [admin] to GET it failure: id=-1",
			expect: ADMIN,
			path:   "-1",
			// HTTP response status: 404 Not Found
			status: http.StatusNotFound,
			jsonMessageError: regex{
				mustCompile: `{"message":"(.*)"}`,
				actual:      `{"message":"User Not Found"}`,
			},
		},

		/*
			delete restore by id [sugriwa]
		*/
		{
			name:   "delete restore by id [sugriwa] to GET it success: id=1",
			expect: SUGRIWA,
			path:   "1",
			// HTTP response status: 404 Not Found
			status: http.StatusNotFound,
			jsonMessageError: regex{
				mustCompile: `{"message":"(.*)"}`,
				actual:      `{"message":"Not Found"}`,
			},
		},

		/*
			No Auth
		*/
		{
			name:   "delete restore by id [no-auth] to GET it failure: id=1",
			expect: ANONYMOUS,
			// HTTP response status: 404 Not Found
			status: http.StatusNotFound,
			jsonMessageError: regex{
				mustCompile: `{"message":"(.*)"}`,
				actual:      `{"message":"Not Found"}`,
			},
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
			// @route: exemple "/admin/restore/3"
			result = noAuth.GET("/admin/restore/{id}", test.path).
				Expect().
				Status(test.status)

			resultBody := result.Body().Raw()

			var (
				mustCompile, actual, match string
				regex                      *regexp.Regexp
			)

			if test.htmlNavbar.mustCompile != "" {
				mustCompile = test.htmlNavbar.mustCompile
				actual = test.htmlNavbar.actual

				regex = regexp.MustCompile(mustCompile)
				match = regex.FindString(resultBody)

				// assert.Equal(t, match, actual)
				//
				// or,
				//
				// assert := assert.New(t)
				// ...
				// assert.Equal(match, actual)
				assert.Equal(t, match, actual)
			}

			if test.htmlHeading.mustCompile != "" {
				mustCompile = test.htmlHeading.mustCompile
				actual = test.htmlHeading.actual

				regex = regexp.MustCompile(mustCompile)
				match = regex.FindString(resultBody)

				assert.Equal(t, match, actual)
			}

			if test.htmlFlashSuccess.mustCompile != "" {
				mustCompile = test.htmlFlashSuccess.mustCompile
				actual = test.htmlFlashSuccess.actual

				regex = regexp.MustCompile(mustCompile)
				match = regex.FindString(resultBody)

				assert.Equal(t, match, actual)
			}

			if test.jsonMessageError.mustCompile != "" {
				mustCompile = test.jsonMessageError.mustCompile
				actual = test.jsonMessageError.actual

				regex = regexp.MustCompile(mustCompile)
				match = regex.FindString(resultBody)

				assert.Equal(t, match, actual)
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
