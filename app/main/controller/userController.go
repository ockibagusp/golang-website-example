package controller

import (
	"fmt"
	"net/http"
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/ockibagusp/golang-website-example/app/main/middleware"
	selectTemplate "github.com/ockibagusp/golang-website-example/app/main/template"
	"github.com/ockibagusp/golang-website-example/app/main/types"

	"github.com/ockibagusp/golang-website-example/business"
	selectUser "github.com/ockibagusp/golang-website-example/business/user"
	locationModules "github.com/ockibagusp/golang-website-example/modules/location"
)

func init() {
	// Templates: userController
	templates := selectTemplate.AppendTemplates
	templates["users/user-all.html"] = selectTemplate.ParseFilesBase("views/users/user-all.html")
	templates["users/user-add.html"] = selectTemplate.ParseFilesBase("views/users/user-add.html", "views/users/user-form.html")
	templates["users/user-read.html"] = selectTemplate.ParseFilesBase("views/users/user-read.html", "views/users/user-form.html")
	templates["users/user-view.html"] = selectTemplate.ParseFilesBase("views/users/user-view.html", "views/users/user-form.html")
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

	uid, _ := c.Get("id").(uint)
	username, _ := c.Get("username").(string)
	role, _ := c.Get("role").(string)

	if role == "anonymous" {
		log.Warn("for GET to users without no-session [@route: /login]")
		middleware.SetFlashError(c, "login process failed!")
		log.Warn("END request method GET for users: [-]failure")
		return c.Redirect(http.StatusFound, "/login")
	}

	// is user?
	if role != "admin" {
		user, err := ctrl.userService.FirstUserByID(ic, uid)
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
	var (
		user *selectUser.User
		err  error
	)

	uid, _ := c.Get("id").(uint)
	username, _ := c.Get("username").(string)
	role, _ := c.Get("role").(string)

	// is user?
	if role == "user" {
		log.Info("START request method GET for create user")
		user, err = ctrl.userService.FirstUserByID(business.InternalContext{}, uid)
		if err != nil {
			log.Warnf(`for GET for create user without select "id" where "username" errors: "%v"`, err)
			log.Warn("END request method GET for create user: [-]failure")
			return err
		}

		middleware.SetFlashError(c, "403 Forbidden")
		log.Infof("END request method GET for create user to users/read/%v: [-]failure", user.ID)
		return c.Redirect(http.StatusFound, fmt.Sprintf("/users/read/%v", user.ID))
	}

	locations, _ := locationModules.NewDB().FindAll(business.InternalContext{})
	if c.Request().Method == "POST" {
		log.Info("START request method POST for create user")

		var location uint
		if c.FormValue("location") != "" {
			location64, err := strconv.ParseUint(c.FormValue("location"), 10, 32)
			if err != nil {
				log.Warnf("for POST to create user without location64 strconv.ParseUint() to error `%v`", err)
				log.Warn("END request method POST for create user: [-]failure")
				// HTTP response status: 400 Bad Request
				return c.HTML(http.StatusBadRequest, err.Error())
			}
			// Location or District?
			location = uint(location64)
		}

		// userForm: type of a user
		userForm := types.UserForm{
			Role:            c.FormValue("role"),
			Username:        c.FormValue("username"),
			Email:           c.FormValue("email"),
			Password:        c.FormValue("password"),
			ConfirmPassword: c.FormValue("confirm_password"),
			Name:            c.FormValue("name"),
			Location:        location,
			Photo:           c.FormValue("photo"),
		}

		// userForm: Validate of a validate user
		err = validation.Errors{
			"username": validation.Validate(
				userForm.Username, validation.Required, validation.Length(4, 15),
			),
			"email": validation.Validate(userForm.Email, validation.Required, validation.Length(5, 30), is.EmailFormat),
			"password": validation.Validate(
				userForm.Password, validation.Required, validation.Length(6, 18),
				validation.By(types.PasswordEquals(userForm.ConfirmPassword)),
			),
			"name":     validation.Validate(userForm.Name, validation.Required, validation.Length(3, 30)),
			"location": validation.Validate(userForm.Location),
			"photo":    validation.Validate(userForm.Photo),
		}.Filter()
		/* if err = validation.Errors{...}.Filter(); err != nil {
			...
		} why?
		*/
		if err != nil {
			log.Warnf("for POST to create user without validation.Errors: `%v`", err)
			middleware.SetFlashError(c, err.Error())

			log.Warn("END request method POST for create user: [-]failure")
			// HTTP response status: 400 Bad Request
			return c.Render(http.StatusBadRequest, "users/user-add.html", echo.Map{
				"name":             "User Add",
				"nav":              "user Add", // (?)
				"session_username": username,
				"session_role":     role,
				"flash_error 	":    middleware.GetFlashError(c),
				"csrf":             c.Get("csrf"),
				"locations":        locations,
				"is_new":           true,
			})
		}

		// Password Hash
		var hash string
		hash, err = middleware.PasswordHash(userForm.Password)
		if err != nil {
			log.Warnf("for POST to create user without middleware.PasswordHash error: `%v`", err)
			log.Warn("END request method POST for create user: [-]failure")
			return err
		}

		user = &selectUser.User{
			Role:     userForm.Role,
			Username: userForm.Username,
			Email:    userForm.Email,
			Password: hash,
			Name:     userForm.Name,
			Location: userForm.Location,
			Photo:    userForm.Photo,
		}

		if _, err := ctrl.userService.Create(business.InternalContext{}, user); err != nil {
			log.Warn("for POST to create user without models.User: nil", "user_failure", user)
			middleware.SetFlashError(c, err.Error())

			log.Warn("END request method POST for create user: [-]failure")
			// HTTP response status: 400 Bad Request
			return c.Render(http.StatusBadRequest, "users/user-add.html", echo.Map{
				"name":             "User Add",
				"nav":              "user Add", // (?)
				"session_username": username,
				"session_role":     role,
				"csrf":             c.Get("csrf"),
				"flash_error":      middleware.GetFlashError(c),
				"locations":        locations,
				"is_new":           true,
			})
		}

		log.Info("models.User: [+]success", "user_success", user)
		middleware.SetFlashSuccess(c, fmt.Sprintf("success new user: %s!", user.Username))
		// create user
		if role == "anonymous" {
			if _, err := middleware.SetSession(user, c); err != nil {
				middleware.SetFlashError(c, err.Error())
				log.Warn("to middleware.SetSession session not found for create user")
				log.Warn("END request method POST for create user: [-]failure")
				// err: session not found
				return c.HTML(http.StatusForbidden, err.Error())
			}
			log.Info("END request method POST for create user: [+]success")
			return c.Redirect(http.StatusMovedPermanently, "/")
		}
		log.Info("END request method POST for create user: [+]success")
		// create admin
		return c.Redirect(http.StatusMovedPermanently, "/users")
	}

	log.Info("START request method GET for create user")
	log.Info("END request method GET for create user: [+]success")
	return c.Render(http.StatusOK, "users/user-add.html", echo.Map{
		"name":             "User Add",
		"nav":              "user Add", // (?)
		"session_username": username,
		"session_role":     role,
		"csrf":             c.Get("csrf"),
		"flash_error":      middleware.GetFlashError(c),
		"locations":        locations,
		"is_new":           true,
	})
}

/*
 * Read User ID
 *
 * @target: Users
 * @method: GET
 * @route: /users/read/:id
 */
func (ctrl *Controller) ReadUser(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	uid := uint(id)
	username, _ := c.Get("username").(string)
	role, _ := c.Get("role").(string)

	if role == "anonymous" {
		log.Warn("for GET to read user without no-session [@route: /login]")
		middleware.SetFlashError(c, "login process failed!")
		log.Warn("END request method GET for read user: [-]failure")
		return c.Redirect(http.StatusFound, "/login")
	}

	log.Info("START request method GET for read user")

	user, err := ctrl.userService.FirstUserByID(business.InternalContext{}, uid)
	if err != nil {
		log.Warnf(
			"for GET to read user without models.User{}.FirstByID() errors: `%v`", err,
		)
		log.Warn("END request method GET for read user: [-]failure")
		// HTTP response status: 406 Method Not Acceptable
		return c.HTML(http.StatusNotAcceptable, err.Error())
	}

	locations, err := locationModules.NewDB().FindAll(business.InternalContext{})
	if err != nil {
		log.Warnf("for GET to read user without models.location{}.FindAll() errors: `%v`", err)
		log.Warn("END request method GET for read user: [-]failure")
		// HTTP response status: 406 Not Acceptable
		return c.HTML(http.StatusNotAcceptable, err.Error())
	}

	log.Info("END request method GET for read user: [+]success")
	return c.Render(http.StatusOK, "users/user-read.html", echo.Map{
		"name":             fmt.Sprintf("User: %s", user.Name),
		"nav":              fmt.Sprintf("User: %s", user.Name), // (?)
		"session_username": username,
		"session_role":     role,
		"flash_success":    middleware.GetFlashSuccess(c),
		"flash_error":      middleware.GetFlashError(c),
		"user":             user,
		"locations":        locations,
		"is_read":          true,
	})
}

/*
 * Update User ID
 *
 * @target: Users
 * @method: GET or POST
 * @route: /users/view/:id
 */
func (ctrl *Controller) UpdateUser(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	uid := uint(id)
	username, _ := c.Get("username").(string)
	role, _ := c.Get("role").(string)

	if role == "anonymous" {
		log.Warn("for GET to update user without no-session [@route: /login]")
		middleware.SetFlashError(c, "login process failed!")
		log.Warn("END request method GET for update user: [-]failure")
		return c.Redirect(http.StatusFound, "/login")
	}

	var (
		user *selectUser.User
		err  error
	)
	user, err = ctrl.userService.FirstUserByID(business.InternalContext{}, uid)
	if err != nil {
		log.Info("START request method GET/POST for update user")
		log.Warnf(
			"for GET to update user without models.User{}.FirstByID() errors: `%v`", err,
		)
		log.Warn("END request method GET for update user: [-]failure")
		// HTTP response status: 404 Not Found
		return c.HTML(http.StatusNotFound, err.Error())
	}

	// admin: yes
	// (role is "user" and (not user.Username)): 403 Forbidden
	if role == "user" && (user.Username != username) {
		log.Info("START request method GET/POST for update user")
		log.Warnf(
			`role is "user" (%v) and [not user.Username (%v)]: 403 Forbidden`,
			role,
			(user.Username != username),
		)
		log.Warn("END request method GET for update user: [-]failure")
		return c.HTML(http.StatusForbidden, "403 Forbidden")
	}

	if c.Request().Method == "POST" {
		log.Info("START request method POST for update user")

		var location uint
		if c.FormValue("location") != "" {
			location64, err := strconv.ParseUint(c.FormValue("location"), 10, 32)
			if err != nil {
				log.Warnf("for POST to create user without location64 strconv.ParseUint() to error `%v`", err)
				log.Warn("END request method POST for create user: [-]failure")
				// HTTP response status: 400 Bad Request
				return c.HTML(http.StatusBadRequest, err.Error())
			}
			// Location or District?
			location = uint(location64)
		}

		updateUser := &selectUser.User{
			Role:     c.FormValue("role"),
			Username: c.FormValue("username"),
			Email:    c.FormValue("email"),
			Name:     c.FormValue("name"),
			Location: location,
			Photo:    c.FormValue("photo"),
		}

		// newUser, err = ctrl.userService.Update(business.InternalContext{}, user, updateUser); err != nil: equal
		if user, err = ctrl.userService.Update(business.InternalContext{}, user, updateUser); err != nil {
			log.Warnf(
				"for POST to update user without models.User{}.Update() errors: `%v`", err,
			)
			middleware.SetFlashError(c, err.Error())
			log.Warn("END request method POST for update user: [-]failure")

			locations, _ := locationModules.NewDB().FindAll(business.InternalContext{})
			// HTTP response status: 406 Method Not Acceptable
			return c.Render(http.StatusNotAcceptable, "users/user-view.html", echo.Map{
				"name":             fmt.Sprintf("User: %s", user.Name),
				"nav":              fmt.Sprintf("User: %s", user.Name), // (?)
				"session_username": username,
				"session_role":     role,
				"flash_error":      middleware.GetFlashError(c),
				"csrf":             c.Get("csrf"),
				"user":             user,
				"locations":        locations,
			})
		}

		log.Info("models.User: [+]success", "user_update", user)
		middleware.SetFlashSuccess(c, fmt.Sprintf("success update user: %s!", user.Username))

		if role == "user" {
			log.Info("END [user] request method POST for update user: [+]success")
			// update user
			return c.Redirect(http.StatusMovedPermanently, "/")
		}
		log.Info("END [admin] request method POST for update user: [+]success")
		// update admin
		return c.Redirect(http.StatusMovedPermanently, "/users")
	}

	log.Info("START request method GET for update user")

	locations, _ := locationModules.NewDB().FindAll(business.InternalContext{})
	if err != nil {
		log.Warnf("for GET to update user without models.location{}.FindAll() errors: `%v`", err)
		log.Warn("END request method GET for update user: [-]failure")
		// HTTP response status: 405 Method Not Allowed
		return c.HTML(http.StatusNotAcceptable, err.Error())
	}

	log.Info("END request method GET for update user: [+]success")
	return c.Render(http.StatusOK, "users/user-view.html", echo.Map{
		"name":             fmt.Sprintf("User: %s", user.Name),
		"nav":              fmt.Sprintf("User: %s", user.Name), // (?)
		"session_username": username,
		"session_role":     role,
		"flash_error":      middleware.GetFlashError(c),
		"csrf":             c.Get("csrf"),
		"user":             user,
		"locations":        locations,
	})
}
