package controller

import (
	"github.com/ockibagusp/golang-website-example/business/auth"
	"github.com/ockibagusp/golang-website-example/business/location"
	"github.com/ockibagusp/golang-website-example/business/user"
	"github.com/ockibagusp/golang-website-example/config"
)

type Controller struct {
	appConfig       *config.Config
	authService     auth.Service
	userService     user.Service
	locationService location.Service
}

func NewController(
	appConfig *config.Config,
	authService auth.Service,
	userService user.Service,
	locationService location.Service,
) *Controller {
	return &Controller{
		appConfig,
		authService,
		userService,
		locationService,
	}
}
