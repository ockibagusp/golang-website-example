package controller

import (
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/ockibagusp/golang-website-example/app/main/middleware"
	selectTemplate "github.com/ockibagusp/golang-website-example/app/main/template"

	"github.com/ockibagusp/golang-website-example/business"
	selectUser "github.com/ockibagusp/golang-website-example/business/user"
	locationModules "github.com/ockibagusp/golang-website-example/modules/location"
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
	ic := business.InternalContext{}

	var (
		users *[]selectUser.User
		err   error

		// typing: all, admin and user
		typing string
	)

	id, _ := c.Get("id").(int)
	username, _ := c.Get("username").(string)
	role, _ := c.Get("role").(string)

	// is user?
	if role != "admin" {
		user, err := ctrl.userService.FirstUserByID(ic, id)
		log.Print(user)
		log.Print(err)
		if err != nil {
			log.Warnf(`for GET for create user without select "id" where "username" errors: "%v"`, err)
			log.Warn("END request method GET for user: [-]failure")
			return err
		}
		log.Infof("END [user] request method GET for users to users/read/%v: [+]success", user.ID)
		return c.Redirect(http.StatusFound, fmt.Sprintf("/users/read/%v", user.ID))
	}

	if c.QueryParam("admin") == "all" {
		log.Infof(`for GET to users admin ctrl.userService.FindAll(ic, "admin")`)
		typing = "Admin"
		users, err = ctrl.userService.FindAll(ic, "admin")
	} else if c.QueryParam("user") == "all" {
		log.Infof(`for GET to users user ctrl.userService.FindAll(ic, "user")`)
		typing = "User"
		users, err = ctrl.userService.FindAll(ic, "user")
	} else {
		log.Infof(`for GET to users ctrl.userService.FindAll(ic) or ctrl.userService.FindAll(ic, "all")`)
		typing = "All"
		users, err = ctrl.userService.FindAll(ic)
	}

	if err != nil {
		log.Warnf("for GET to users without ctrl.userService.FindAll errors: `%v`", err)
		log.Warn("END request method GET for users: [-]failure")
		// HTTP response status: 404 Not Found
		return c.HTML(http.StatusNotFound, err.Error())
	}

	log.Info("END request method GET for users: [+]success")
	return c.Render(http.StatusOK, "users/user-all.html", echo.Map{
		"name":             fmt.Sprintf("Users: %v", typing),
		"nav":              "users", // (?)
		"session_username": username,
		"session_role":     role,
		/*
			"flash": echo.Map{"success": ..., "error": ...}
			or,
			"flash_success": ....
			"flash_error": ....
		*/

		"flash": echo.Map{
			"success": middleware.GetFlashSuccess(c),
			"error":   middleware.GetFlashError(c),
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
	ic := business.NewInternalContext("create user")

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

	locations, _ := locationModules.NewDB().FindAll(ic)
	return c.Render(http.StatusOK, "users/user-add.html", echo.Map{
		"name":        "User Add",
		"nav":         "user Add", // (?)
		"session":     session,
		"csrf":        c.Get("csrf"),
		"flash_error": []string{},
		"locations":   locations,
		"is_new":      true,
	})
}
