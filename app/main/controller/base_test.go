package controller_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect/v2"
	ctrl "github.com/ockibagusp/golang-website-example/app/main/controller"
	"github.com/ockibagusp/golang-website-example/app/main/router"
	"github.com/ockibagusp/golang-website-example/business/auth"
	"github.com/ockibagusp/golang-website-example/business/user"
	"github.com/ockibagusp/golang-website-example/config"
	userModule "github.com/ockibagusp/golang-website-example/modules/user"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

const (
	ADMIN      string = "admin"
	SUGRIWA           = "sugriwa"
	SUBALI            = "subali"
	OCKIBAGUSP        = "ockibagusp"
)

var conf *config.Config = config.GetAPPConfig()

func TestController(t *testing.T) {
	/*
		assert := assert.New(t)
		assert.NotNil(setupTestController())

		or,
	*/
	assert.NotNil(t, setupTestController())
}

// setup test Handler
func setupTestHandler() http.Handler {
	return router.RegisterPath(
		conf,
		setupTestController(),
	)
}

func newUserService(db *gorm.DB) user.Service {
	userRepo := userModule.NewGormRepository(db)

	// userService
	return user.NewService(userRepo)
}

// Controller test
func setupTestController() *ctrl.Controller {
	db := conf.GetDatabaseConnection()

	userService := newUserService(db)
	authService := auth.NewService(userService)
	return ctrl.NewController(
		conf,
		authService,
		userService,
	)
}

/*
Setup test sever

repository: .env
1. function conf.GetSessionTest()
@SESSION_TEST: {true} or {false}

2. function conf.GetDebug()
@DEBUG: {true} or {false}
*/
func SetupTestServer(t *testing.T, debug ...bool) (no_auth *httpexpect.Expect) {
	conf.GetSessionTest()
	conf.GetDebug()

	handler := setupTestHandler()

	server := httptest.NewServer(handler)
	defer server.Close()

	new_config := httpexpect.Config{
		BaseURL: server.URL,
		Client: &http.Client{
			Transport: httpexpect.NewBinder(handler),
			Jar:       httpexpect.NewJar(),
		},
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewCompactPrinter(t),
		},
	}

	if conf.GetDebugAsTrue(debug) {
		new_config.Printers = []httpexpect.Printer{
			httpexpect.NewDebugPrinter(t, true),
		}
	}

	no_auth = httpexpect.WithConfig(new_config)
	return
}
