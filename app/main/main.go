package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"golang-website-example/app/main/controller"
	"golang-website-example/app/main/router"
	"golang-website-example/business/auth"
	"golang-website-example/business/user"
	"golang-website-example/config"
	userModule "golang-website-example/modules/user"

	"gorm.io/gorm"
)

func newUserService(db *gorm.DB) user.Service {
	userRepo := userModule.NewGormRepository(db)

	// userService
	return user.NewService(userRepo)
}

func main() {
	conf := config.GetAPPConfig()
	db := conf.GetDatabaseConnection()

	userService := newUserService(db)
	authService := auth.NewService(userService)

	controllerAPP := controller.NewController(
		conf,
		authService,
		userService,
	)

	e := router.RegisterPath(
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
