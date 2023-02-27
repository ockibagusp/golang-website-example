package user

import (
	"github.com/ockibagusp/golang-website-example/business"
)

type (
	User struct {
		business.Model
		// enum: admin and user
		Role string
		// database: just `username` varchar 15
		Username string `gorm:"unique;not null;type:varchar(15)" form:"username"`
		Email    string `gorm:"unique;not null;type:varchar(30)" form:"email"`
		Password string `gorm:"not null" form:"password"`
		Name     string `gorm:"not null;type:varchar(30)" form:"name"`
		Location uint   `form:"location"`
		Photo    string `form:"photo"`

		business.ObjectMetadata
	}
)
