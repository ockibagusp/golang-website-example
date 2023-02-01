package auth

import (
	"github.com/ockibagusp/golang-website-example/business"
	"github.com/ockibagusp/golang-website-example/business/user"
	"golang.org/x/crypto/bcrypt"
)

type (
	service struct {
		userService user.Service
	}

	Service interface {
		VerifyLogin(ic business.InternalContext, email string, plainPassword string) (getUser user.User, validPassword bool)
		CheckHashPassword(hash, password string) bool
	}
)

func NewService(userService user.Service) Service {
	return &service{
		userService,
	}
}

// VerifyLogin: get-user and valid-password
func (s *service) VerifyLogin(ic business.InternalContext, email string, plainPassword string) (getUser user.User, validPassword bool) {
	getUser, err := s.userService.FindByEmail(ic, email)
	if err != nil {
		return
	}

	passwordHash, err := func(plainPassword string) (string, error) {
		// GenerateFromPassword(..., bcrypt.DefaultCost{=10})
		hash, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
		return string(hash), err
	}(plainPassword)
	if err != nil {
		return
	}

	if passwordHash != getUser.Password {
		return
	}

	validPassword = true

	return
}

// CheckHashPassword: hashes for passwords to bool
func (s *service) CheckHashPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
