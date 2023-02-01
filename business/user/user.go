package user

import (
	"github.com/ockibagusp/golang-website-example/business"
	"gorm.io/gorm"
)

type (
	User struct {
		gorm.Model
		Role     string
		Name     string
		Email    string
		Password string

		business.ObjectMetadata
	}
)
