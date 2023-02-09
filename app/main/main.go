package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/ockibagusp/golang-website-example/app/main/controller"
	"github.com/ockibagusp/golang-website-example/app/main/router"
	"github.com/ockibagusp/golang-website-example/business/auth"
	"github.com/ockibagusp/golang-website-example/business/location"
	"github.com/ockibagusp/golang-website-example/business/user"
	"github.com/ockibagusp/golang-website-example/config"
	locationModule "github.com/ockibagusp/golang-website-example/modules/location"
	userModule "github.com/ockibagusp/golang-website-example/modules/user"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var logger = log.New()

func newUserService(db *gorm.DB) user.Service {
	userRepo := userModule.NewGormRepository(db)

	// userService
	return user.NewService(userRepo)
}

func newLocationService(db *gorm.DB) (locationService location.Service) {
	locationRepo := locationModule.NewGormRepository(db)

	locationService = location.NewService(locationRepo)
	return
}

func main() {
	conf := config.GetAPPConfig()
	db := conf.GetDatabaseConnection()

	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())

	userService := newUserService(db)
	locationService := newLocationService(db)
	authService := auth.NewService(userService)

	controllerAPP := controller.NewController(
		conf,
		authService,
		userService,
		locationService,
	)

	router.RegisterPath(
		e,
		conf,
		controllerAPP,
	)

	// start the Echo server
	go func() {
		if err := e.Start(":8000"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the Echo server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	// a timeout context after 10 seconds
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// shutdown the Echo server
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(fmt.Sprintf("failed the Echo server: %v", err))
	} else {
		e.Logger.Info("successfully the Echo server")
	}
}
