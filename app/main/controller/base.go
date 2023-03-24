package controller

import (
	"github.com/ockibagusp/golang-website-example/business/auth"
	"github.com/ockibagusp/golang-website-example/business/user"
	"github.com/ockibagusp/golang-website-example/config"
	"github.com/ockibagusp/golang-website-example/logger"
)

type Controller struct {
	appConfig   *config.Config
	authService auth.Service
	userService user.Service
	logger      *logger.Logger
}

func NewController(
	appConfig *config.Config,
	authService auth.Service,
	userService user.Service,
) *Controller {
	logger := logger.New()
	return &Controller{
		appConfig,
		authService,
		userService,
		logger,
	}
}
