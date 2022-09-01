package controllers

import (
	"github.com/ockibagusp/golang-website-example/db"
	"gorm.io/gorm"
)

// Controller is a controller for this application
type Controller struct {
	DB *gorm.DB
}

// New Controller
func New() *Controller {
	// PROD or DEV
	db_manager := db.Init("PROD")

	return &Controller{
		DB: db_manager,
	}
}
