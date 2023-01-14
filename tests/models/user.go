package models

import (
	"errors"

	"github.com/gorilla/sessions"
	"github.com/ockibagusp/golang-website-example/models"
)

var UserSelectTest string

/*
 * Users Test
 */
var UsersTest []models.User = []models.User{
	{
		Username: "admin",
		IsAdmin:  1,
	},
	{
		Username: "sugriwa",
	},
	{
		Username: "subali",
	},
	{
		Username: "ockibagusp",
	},
}

func UserUsername(user *models.User) {
	for _, testUser := range UsersTest {
		if UserSelectTest == testUser.Username {
			user.Username = testUser.Username
		}
	}
}

func GetAuthSession() (session_gorilla *sessions.Session, err error) {
	if UserSelectTest == "" {
		session_gorilla = &sessions.Session{
			Values: map[interface{}]interface{}{
				"username":     "",
				"is_auth_type": -1,
			},
		}

		err = errors.New("no session")
		return
	}

	for _, testUser := range UsersTest {
		if UserSelectTest == testUser.Username {
			session_gorilla = &sessions.Session{
				Values: map[interface{}]interface{}{
					"username": testUser.Username,
				},
			}

			if testUser.IsAdmin == 1 {
				session_gorilla.Values["is_auth_type"] = 1 // admin: 1
			} else if testUser.IsAdmin == 0 {
				session_gorilla.Values["is_auth_type"] = 2 // user: 2
			}
		}
	}
	return
}
