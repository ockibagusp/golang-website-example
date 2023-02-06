package controller

import (
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	selectTemplate "github.com/ockibagusp/golang-website-example/app/main/template"

	"github.com/ockibagusp/golang-website-example/business"
	selectUser "github.com/ockibagusp/golang-website-example/business/user"
)

func init() {
	// Templates: userController
	templates := selectTemplate.AppendTemplates
	templates["users/user-all.html"] = selectTemplate.ParseFilesBase("views/users/user-all.html")
}

func (ctrl *Controller) Users(c echo.Context) error {
	session := sessions.Session{}

	var (
		users []selectUser.User
		// err   error

		// typing: all, admin and user
		typing string
	)

	users = append(users, selectUser.User{
		Model:    business.Model{ID: 1},
		Role:     "user",
		Username: "ockibagusp",
		Name:     "Ocki Bagus Pratama",
	})

	typing = "all"

	return c.Render(http.StatusOK, "users/user-all.html", echo.Map{
		"name":    fmt.Sprintf("Users: %v", typing),
		"nav":     "users", // (?)
		"session": session,
		/*
			"flash": echo.Map{"success": ..., "error": ...}
			or,
			"flash_success": ....
			"flash_error": ....
		*/

		"flash": echo.Map{
			"success": []string{},
			"error":   []string{},
		},
		"users": users,
	})
}
