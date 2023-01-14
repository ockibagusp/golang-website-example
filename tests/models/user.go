package models

import "github.com/ockibagusp/golang-website-example/models"

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
