package models

import "github.com/ockibagusp/golang-website-example/models"

/*
 * Test Users
 */
var TestUsers []models.User = []models.User{
	{
		Username: "admin",
		IsAdmin:  1,
	},
	{
		Username: "sugriwa",
		IsAdmin:  0,
	},
	{
		Username: "subali",
		IsAdmin:  0,
	},
	{
		Username: "ockibagusp",
		IsAdmin:  0,
	},
}
