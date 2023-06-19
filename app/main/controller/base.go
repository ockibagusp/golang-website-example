package controller

import (
	"golang-website-example/business/auth"
	"golang-website-example/business/user"
	"golang-website-example/config"
)

type Controller struct {
	appConfig   *config.Config
	authService auth.Service
	userService user.Service
}

func NewController(
	appConfig *config.Config,
	authService auth.Service,
	userService user.Service,
) *Controller {
	return &Controller{
		appConfig,
		authService,
		userService,
	}
}
