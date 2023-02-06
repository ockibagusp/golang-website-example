package controller

import (
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	selectTemplate "github.com/ockibagusp/golang-website-example/app/main/template"

	"github.com/ockibagusp/golang-website-example/business"
	selectedCities "github.com/ockibagusp/golang-website-example/business/city"
	selectUser "github.com/ockibagusp/golang-website-example/business/user"
)

func init() {
	// Templates: userController
	templates := selectTemplate.AppendTemplates
	templates["users/user-all.html"] = selectTemplate.ParseFilesBase("views/users/user-all.html")
	templates["users/user-add.html"] = selectTemplate.ParseFilesBase("views/users/user-add.html", "views/users/user-form.html")
}

/*
 * Users All
 *
 * @target: Users
 * @method: GET
 * @route: /users
 */
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

/*
 * User Add
 *
 * @target: Users
 * @method: GET or POST
 * @route: /users/add
 */
func (ctrl *Controller) CreateUser(c echo.Context) error {
	session := sessions.Session{}

	var (
		users []selectUser.User
		// err error
	)

	if c.Request().Method == "POST" {
		users = append(users, selectUser.User{
			Model:    business.Model{ID: 1},
			Role:     "user",
			Username: "ockibagusp",
			Name:     "Ocki Bagus Pratama",
		})

		// create admin
		return c.Redirect(http.StatusMovedPermanently, "/users")
	}

	return c.Render(http.StatusOK, "users/user-add.html", echo.Map{
		"name":        "User Add",
		"nav":         "user Add", // (?)
		"session":     session,
		"csrf":        c.Get("csrf"),
		"flash_error": []string{},
		"cities":      selectedCities.City{},
		"is_new":      true,
	})
}
