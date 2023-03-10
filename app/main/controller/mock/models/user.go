package models

import (
	"github.com/gorilla/sessions"
	"github.com/ockibagusp/golang-website-example/business"
	selectUser "github.com/ockibagusp/golang-website-example/business/user"
)

var UserSelectTest string

/*
 * Users Test
 */
var UsersTest []selectUser.User = []selectUser.User{
	{
		Model:    business.Model{ID: 1},
		Username: "admin",
		Role:     "admin",
	},
	{
		Model:    business.Model{ID: 2},
		Username: "sugriwa",
		Role:     "user",
	},
	{
		Model:    business.Model{ID: 3},
		Username: "subali",
		Role:     "user",
	},
	{
		Model:    business.Model{ID: 14},
		Username: "ockibagusp",
		Role:     "user",
	},
}

func SetAuthSession(user *selectUser.User) {
	for _, testUser := range UsersTest {
		if UserSelectTest == testUser.Username {
			user.Username = testUser.Username
		}
	}
}

func GetAuthSession() (session_gorilla *sessions.Session, err error) {
	if UserSelectTest == "anonymous" {
		session_gorilla = &sessions.Session{
			Values: map[interface{}]interface{}{
				"id":       -1,
				"username": "anonymous",
				"role":     "anonymous",
			},
		}

		err = nil
		return
	}

	for _, testUser := range UsersTest {
		if UserSelectTest == testUser.Username {
			session_gorilla = &sessions.Session{
				Values: map[interface{}]interface{}{
					"id":       testUser.Model.ID,
					"username": testUser.Username,
					"role":     testUser.Role,
				},
			}
		}
	}
	return
}
