package controller_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect/v2"
	ctrl "github.com/ockibagusp/golang-website-example/app/main/controller"
	"github.com/ockibagusp/golang-website-example/app/main/router"
	"github.com/ockibagusp/golang-website-example/business"
	"github.com/ockibagusp/golang-website-example/business/auth"
	"github.com/ockibagusp/golang-website-example/business/user"
	"github.com/ockibagusp/golang-website-example/config"
	userModule "github.com/ockibagusp/golang-website-example/modules/user"
	"gorm.io/gorm"
)

var conf *config.Config = config.GetAPPConfig()
var db *gorm.DB = conf.GetDatabaseConnection()

// truncate Users
func truncateUsers() {
	db.Exec("TRUNCATE users;")

	// database: just `users.username` varchar 15
	users := []user.User{
		{
			Model:    business.Model{ID: 1},
			Username: "admin",
			Email:    "admin@website.com",
			Password: "$2a$10$XJAj65HZ2c.n1iium4qUEeGarW0PJsqVcedBh.PDGMXdjqfOdN1hW",
			Name:     "Admin",
			Role:     "admin",
			Photo:    "members/admin3981.png",
		},
		{
			Model:    business.Model{ID: 2},
			Username: "sugriwa",
			Email:    "sugriwa@wanara.com",
			Password: "$2a$10$bVVMuFHe/iaydX9yO2AttOPT8WyhMPe9F8nDflEqEyJbGRD5.guFu",
			Name:     "Sugriwa",
			Role:     "user",
			Photo:    "members/sugriwa2492.png",
		},
		{
			Model:    business.Model{ID: 3},
			Username: "subali",
			Email:    "subali@wanara.com",
			Password: "$2a$10$eO8wPLSfBU.8KLUh/T9kDeBm0vIRjiCvsmWe8ou5fZHJ3cYAUcg6y",
			Name:     "Subali",
			Role:     "user",
			Photo:    "members/subali453.png",
		},
		{
			Model:    business.Model{ID: 4},
			Username: "ockibagusp",
			Email:    "ocki.bagus.p@gmail.com",
			Password: "$2a$10$Y3UewQkjw808Ig90OPjuq.zFYIUGgFkWBuYiKzwLK8n3t9S8RYuYa",
			Name:     "Ocki Bagus Pratama",
			Role:     "user",
			Photo:    "members/ockibagusp981495792267526.jpg",
		},
	}

	for _, user := range users {
		_, err := newUserService(db).Create(business.InternalContext{}, &user)
		if err != nil {
			panic("Username not already: " + err.Error())
		}
	}
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
func setupTestServer(t *testing.T, debug ...bool) (noAuth *httpexpect.Expect) {
	conf.GetSessionTest()
	conf.GetDebug()

	handler := setupTestHandler()

	server := httptest.NewServer(handler)
	defer server.Close()

	newConfig := httpexpect.Config{
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
		newConfig.Printers = []httpexpect.Printer{
			httpexpect.NewDebugPrinter(t, true),
		}
	}

	noAuth = httpexpect.WithConfig(newConfig)
	return
}
