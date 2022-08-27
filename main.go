package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/ockibagusp/golang-website-example/controllers"
	"github.com/ockibagusp/golang-website-example/router"
)

func main() {
	// controllers init
	controllers := controllers.New()

	// Echo: router
	e := router.New(controllers)

	// start the Echo server
	go e.Logger.Fatal(e.Start(":8000"))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	// defines a timeout context that will be canceled after 10 seconds
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// shutdown the Echo server
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(fmt.Sprintf("failed to shutting down Echo server: %v", err))
	} else {
		e.Logger.Info("successfully shutting down Echo server")
	}
}
