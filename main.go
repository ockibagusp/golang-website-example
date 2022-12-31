package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/ockibagusp/golang-website-example/controllers"
	"github.com/ockibagusp/golang-website-example/router"
)

func main() {
	os.Setenv("session_test", "0")

	// controllers init
	controllers := controllers.New()

	// Echo: router
	e := router.New(controllers)

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
