package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/labstack/echo/v4"
	"github.com/ockibagusp/golang-website-example/middleware"
	"github.com/ockibagusp/golang-website-example/models"
	"github.com/ockibagusp/golang-website-example/types"
	log "github.com/sirupsen/logrus"
)

/*
 * Users All
 *
 * @target: Users
 * @method: GET
 * @route: /users
 */
func (controller *Controller) Users(c echo.Context) error {
	session, _ := middleware.GetAuth(c)
	log := log.WithFields(log.Fields{
		"username": session.Values["username"],
		"route":    c.Path(),
	})
	log.Info("START request method GET for users")

	if session.Values["is_auth_type"] == -1 {
		log.Warn("for GET to users without no-session [@route: /login]")
		middleware.SetFlashError(c, "login!")
		log.Warn("END request method GET for users: [-]failure")
		return c.Redirect(http.StatusFound, "/login")
	}

	// models.User{} or (models.User{}) or var user models.User or user := models.User{}
	users, err := models.User{}.FindAll(controller.DB)
	if err != nil {
		log.Warnf("for GET to users without models.User{}.FindAll() errors: `%v`", err)
		log.Warn("END request method GET for users: [-]failure")
		// HTTP response status: 404 Not Found
		return c.HTML(http.StatusNotFound, err.Error())
	}

	log.Info("END request method GET for users: [+]success")
	return c.Render(http.StatusOK, "users/user-all.html", echo.Map{
		"name":    "Users",
		"nav":     "users", // (?)
		"session": session,
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
func (controller *Controller) CreateUser(c echo.Context) error {
	session, _ := middleware.GetAuth(c)
	log := log.WithFields(log.Fields{
		"username": session.Values["username"],
		"route":    c.Path(),
	})

	if c.Request().Method == "POST" {
		log.Info("START request method POST for create user")

		var city uint
		if c.FormValue("city") != "" {
			city64, err := strconv.ParseUint(c.FormValue("city"), 10, 32)
			if err != nil {
				log.Warnf("for POST to create user without city64 strconv.ParseUint() to error `%v`", err)
				log.Warn("END request method POST for create user: [-]failure")
				// HTTP response status: 400 Bad Request
				return c.HTML(http.StatusBadRequest, err.Error())
			}
			// City and District ?
			city = uint(city64)
		}

		// userForm: type of a user
		_userForm := types.UserForm{
			Username:        c.FormValue("username"),
			Email:           c.FormValue("email"),
			Password:        c.FormValue("password"),
			ConfirmPassword: c.FormValue("confirm_password"),
			Name:            c.FormValue("name"),
			City:            city,
			Photo:           c.FormValue("photo"),
		}

		// _userForm: Validate of a validate user
		err := validation.Errors{
			"username": validation.Validate(
				_userForm.Username, validation.Required, validation.Length(4, 15),
			),
			"email": validation.Validate(_userForm.Email, validation.Required, is.Email),
			"password": validation.Validate(
				_userForm.Password, validation.Required, validation.Length(6, 18),
				validation.By(types.PasswordEquals(_userForm.ConfirmPassword)),
			),
			"name":  validation.Validate(_userForm.Name, validation.Required),
			"city":  validation.Validate(_userForm.City),
			"photo": validation.Validate(_userForm.Photo),
		}.Filter()
		/* if err = validation.Errors{...}.Filter(); err != nil {
			...
		} why?
		*/
		if err != nil {
			log.Warnf("for POST to create user without validation.Errors: `%v`", err)
			middleware.SetFlashError(c, err.Error())

			cities, _ := models.City{}.FindAll(controller.DB)
			log.Warn("END request method POST for create user: [-]failure")
			// HTTP response status: 400 Bad Request
			return c.Render(http.StatusBadRequest, "users/user-add.html", echo.Map{
				"name":        "User Add",
				"nav":         "user Add", // (?)
				"session":     session,
				"flash_error": middleware.GetFlashError(c),
				"csrf":        c.Get("csrf"),
				"cities":      cities,
				"is_new":      true,
			})
		}

		// Password Hash
		hash, err := middleware.PasswordHash(_userForm.Password)
		if err != nil {
			log.Warnf("for POST to create user without middleware.PasswordHash error: `%v`", err)
			log.Warn("END request method POST for create user: [-]failure")
			return err
		}

		user := models.User{
			Username: _userForm.Username,
			Email:    _userForm.Email,
			Password: hash,
			Name:     _userForm.Name,
			City:     _userForm.City,
			Photo:    _userForm.Photo,
		}

		// _, err := user.Save(...): be able
		if _, err := user.Save(controller.DB); err != nil {
			log.WithField("user_failure", user).
				Warn("for POST to create user without models.User: nil")
			middleware.SetFlashError(c, err.Error())

			cities, _ := models.City{}.FindAll(controller.DB)
			log.Warn("END request method POST for create user: [-]failure")
			// HTTP response status: 400 Bad Request
			return c.Render(http.StatusBadRequest, "users/user-add.html", echo.Map{
				"name":        "User Add",
				"nav":         "user Add", // (?)
				"session":     session,
				"csrf":        c.Get("csrf"),
				"flash_error": middleware.GetFlashError(c),
				"cities":      cities,
				"is_new":      true,
			})
		}

		log.WithField("user_success", user).Info("models.User: [+]success")
		middleware.SetFlashSuccess(c, fmt.Sprintf("success new user: %s!", user.Username))
		log.Info("END request method POST for create user: [+]success")
		return c.Redirect(http.StatusMovedPermanently, "/users")
	}

	log.Info("START request method GET for create user")

	// models.City{} or (models.City{}) or var city models.City or city := models.City{}
	cities, err := models.City{}.FindAll(controller.DB)
	if err != nil {
		log.Warnf("for GET to create user without models.City{}.FindAll() errors: `%v`", err)
		log.Warn("END request method GET for create user: [-]failure")
		// HTTP response status: 405 Method Not Allowed
		return c.HTML(http.StatusNotAcceptable, err.Error())
	}

	log.Info("END request method GET for create user: [+]success")
	return c.Render(http.StatusOK, "users/user-add.html", echo.Map{
		"name":        "User Add",
		"nav":         "user Add", // (?)
		"session":     session,
		"csrf":        c.Get("csrf"),
		"flash_error": middleware.GetFlashError(c),
		"cities":      cities,
		"is_new":      true,
	})
}

/*
 * Read User ID
 *
 * @target: Users
 * @method: GET
 * @route: /users/read/:id
 */
func (controller *Controller) ReadUser(c echo.Context) error {
	session, _ := middleware.GetAuth(c)
	log := log.WithFields(log.Fields{
		"username": session.Values["username"],
		"route":    fmt.Sprintf("%v -> id:%v", c.Path(), c.Param("id")),
	})
	if session.Values["is_auth_type"] == -1 {
		log.Warn("for GET to read user without no-session [@route: /login]")
		middleware.SetFlashError(c, "login!")
		log.Warn("END request method GET for read user: [-]failure")
		return c.Redirect(http.StatusFound, "/login")
	}

	log.Info("START request method GET for read user")

	id, _ := strconv.Atoi(c.Param("id"))

	// var user models.User
	// ...
	// _user, err := user.FirstByID(...): be able
	user, err := models.User{}.FirstByID(controller.DB, id)
	if err != nil {
		log.Warnf(
			"for GET to read user without models.User{}.FirstByID() errors: `%v`", err,
		)
		log.Warn("END request method GET for read user: [-]failure")
		// HTTP response status: 405 Method Not Allowed
		return c.HTML(http.StatusNotAcceptable, err.Error())
	}

	// models.City{} or (models.City{}) or var city models.City or city := models.City{}
	cities, err := models.City{}.FindAll(controller.DB)
	if err != nil {
		log.Warnf("for GET to read user without models.City{}.FindAll() errors: `%v`", err)
		log.Warn("END request method GET for read user: [-]failure")
		// HTTP response status: 406 Not Acceptable
		return c.HTML(http.StatusNotAcceptable, err.Error())
	}

	log.Info("END request method GET for read user: [+]success")
	return c.Render(http.StatusOK, "users/user-read.html", echo.Map{
		"name":    fmt.Sprintf("User: %s", user.Name),
		"nav":     fmt.Sprintf("User: %s", user.Name), // (?)
		"session": session,
		"user":    user,
		"cities":  cities,
		"is_read": true,
	})
}

/*
 * Update User ID
 *
 * @target: Users
 * @method: GET or POST
 * @route: /users/view/:id
 */
func (controller *Controller) UpdateUser(c echo.Context) error {
	session, _ := middleware.GetAuth(c)
	log := log.WithFields(log.Fields{
		"username": session.Values["username"],
		"route":    fmt.Sprintf("%v -> id:%v", c.Path(), c.Param("id")),
	})
	if session.Values["is_auth_type"] == -1 {
		log.Warn("for GET to update user without no-session [@route: /login]")
		middleware.SetFlashError(c, "login!")
		log.Warn("END request method GET for update user: [-]failure")
		return c.Redirect(http.StatusFound, "/login")
	}

	id, _ := strconv.Atoi(c.Param("id"))

	if c.Request().Method == "POST" {
		log.Info("START request method POST for update user")

		var user models.User
		cities, _ := models.City{}.FindAll(controller.DB)

		// HTTP response status: 400 Bad Request
		if err := c.Bind(&user); err != nil {
			log.Warnf(
				"for POST to update user without c.Bind() errors: `%v`", err,
			)
			middleware.SetFlashError(c, err.Error())
			log.Warn("END request method POST for update user: [-]failure")
			return c.Render(http.StatusBadRequest, "users/user-view.html", echo.Map{
				"name":        fmt.Sprintf("User: %s", user.Name),
				"nav":         fmt.Sprintf("User: %s", user.Name), // (?)
				"session":     session,
				"flash_error": middleware.GetFlashError(c),
				"csrf":        c.Get("csrf"),
				"user":        user,
				"cities":      cities,
			})
		}

		// _, err = user.Update(...): be able
		if _, err := user.Update(controller.DB, id); err != nil {
			log.Warnf(
				"for POST to update user without models.User{}.Update() errors: `%v`", err,
			)
			middleware.SetFlashError(c, err.Error())
			log.Warn("END request method POST for update user: [-]failure")
			// HTTP response status: 405 Method Not Allowed
			return c.Render(http.StatusNotAcceptable, "users/user-view.html", echo.Map{
				"name":        fmt.Sprintf("User: %s", user.Name),
				"nav":         fmt.Sprintf("User: %s", user.Name), // (?)
				"session":     session,
				"flash_error": middleware.GetFlashError(c),
				"csrf":        c.Get("csrf"),
				"user":        user,
				"cities":      cities,
			})
		}

		log.WithField("user_update", user).Info("models.User: [+]success")
		middleware.SetFlashSuccess(c, fmt.Sprintf("success update user: %s!", user.Username))
		log.Info("END request method POST for update user: [+]success")
		return c.Redirect(http.StatusMovedPermanently, "/users")
	}

	log.Info("START request method GET for update user")

	// var user models.User
	// ...
	// _user, err := user.FirstByID(...): be able
	user, err := models.User{}.FirstByID(controller.DB, id)
	if err != nil {
		log.Warnf(
			"for GET to update user without models.User{}.FirstByID() errors: `%v`", err,
		)
		log.Warn("END request method GET for update user: [-]failure")
		// HTTP response status: 405 Method Not Allowed
		return c.HTML(http.StatusNotAcceptable, err.Error())
	}

	// models.City{} or (models.City{}) or var city models.City or city := models.City{}
	cities, err := models.City{}.FindAll(controller.DB)
	if err != nil {
		log.Warnf("for GET to update user without models.City{}.FindAll() errors: `%v`", err)
		log.Warn("END request method GET for update user: [-]failure")
		// HTTP response status: 405 Method Not Allowed
		return c.HTML(http.StatusNotAcceptable, err.Error())
	}

	log.Info("END request method GET for update user: [+]success")
	return c.Render(http.StatusOK, "users/user-view.html", echo.Map{
		"name":        fmt.Sprintf("User: %s", user.Name),
		"nav":         fmt.Sprintf("User: %s", user.Name), // (?)
		"session":     session,
		"flash_error": middleware.GetFlashError(c),
		"csrf":        c.Get("csrf"),
		"user":        user,
		"cities":      cities,
	})
}

/*
 * Update User ID by Password
 *
 * @target: Users
 * @method: GET or POST
 * @route: /users/view/:id/password
 */
func (controller *Controller) UpdateUserByPassword(c echo.Context) error {
	session, _ := middleware.GetAuth(c)
	log := log.WithFields(log.Fields{
		"username": session.Values["username"],
		"route":    fmt.Sprintf("%v -> id:%v", c.Path(), c.Param("id")),
	})
	if session.Values["is_auth_type"] == -1 {
		log.Warn("for GET to update user by password without no-session [@route: /login]")
		middleware.SetFlashError(c, "login!")
		log.Warn("END request method GET for read user: [-]failure")
		return c.Redirect(http.StatusFound, "/login")
	}

	log.Info("START request method GET or POST for update user by password")

	id, _ := strconv.Atoi(c.Param("id"))

	/*
		TODO:
		for example:
		username ockibagusp update by password 'ockibagusp': ok
		username ockibagusp update by password 'sugriwa': no
	*/
	user, err := models.User{}.FirstByIDAndUsername(
		controller.DB, id, session.Values["username"].(string),
	)
	if err != nil {
		log.Warnf(
			"for GET to update user by password without models.User{}.FirstByIDAndUsername() errors: `%v`", err,
		)
		log.Warn("END request method GET for update user by password: [-]failure")
		// HTTP response status: 405 Method Not Allowed
		return c.HTML(http.StatusNotAcceptable, err.Error())
	}

	if c.Request().Method == "POST" {
		// newPasswordForm: type of a password user
		_newPasswordForm := types.NewPasswordForm{
			OldPassword:        c.FormValue("old_password"),
			NewPassword:        c.FormValue("new_password"),
			ConfirmNewPassword: c.FormValue("confirm_new_password"),
		}

		if !middleware.CheckHashPassword(user.Password, _newPasswordForm.OldPassword) {
			log.Warnf("for POST to update user by password without !middleware.CheckHashPassword() errors: `%v`", err)
			middleware.SetFlashError(c, "check hash password is wrong!")
			log.Warn("END request method POST for update user by password: [-]failure")
			return c.Render(http.StatusForbidden, "user-view-password.html", echo.Map{
				"name":         fmt.Sprintf("User: %s", user.Name),
				"session":      session,
				"flash_error":  middleware.GetFlashError(c),
				"user":         user,
				"is_html_only": true,
			})
		}

		// _newPasswordForm: Validate of a validate user
		err := validation.Errors{
			"password": validation.Validate(
				_newPasswordForm.NewPassword, validation.Required, validation.Length(6, 18),
				validation.By(types.PasswordEquals(_newPasswordForm.ConfirmNewPassword)),
			),
		}.Filter()
		/* if err = validation.Errors{...}.Filter(); err != nil {
			...
		} why?
		*/
		if err != nil {
			log.Warnf("for POST to update user by password without validation.Errors errors: `%v`", err)
			middleware.SetFlashError(c, err.Error())
			log.Warn("END request method POST for update user by password: [-]failure")
			// return c.JSON(http.StatusBadRequest, echo.Map{
			// 	"message": "Passwords Don't Match",
			// })
			return c.Render(http.StatusForbidden, "user-view-password.html", echo.Map{
				"name":         fmt.Sprintf("User: %s", user.Name),
				"session":      session,
				"flash_error":  middleware.GetFlashError(c),
				"user":         user,
				"is_html_only": true,
			})
		}

		// Password Hash
		hash, err := middleware.PasswordHash(_newPasswordForm.NewPassword)
		if err != nil {
			log.Warnf("for POST to update user by password without middleware.PasswordHash() errors: `%v`", err)
			log.Warn("END request method POST for update user by password: [-]failure")
			return err
		}

		// err := user.UpdateByIDandPassword(...): be able
		if err := user.UpdateByIDandPassword(controller.DB, id, hash); err != nil {
			log.Warnf("for POST to update user by password without models.User{}.UpdateByIDandPassword() errors: `%v`", err)
			log.Warn("END request method POST for update user by password: [-]failure")
			// HTTP response status: 405 Method Not Allowed
			return c.HTML(http.StatusNotAcceptable, err.Error())
		}

		log.WithField("user_update_password", user).Info("models.User: [+]success")
		middleware.SetFlashSuccess(c, fmt.Sprintf("success update user by password: %s!", user.Username))
		log.Info("END request method POST for update user by password: [+]success")
		return c.Redirect(http.StatusMovedPermanently, "/users")
	}

	log.Info("END request method GET for update user by password: [+]success")
	/*
		name (string): "users/user-view-password.html" -> no
			{..,"status":500,"error":"html/template: \"users/user-view-password.html\" is undefined",..}
			why?
		name (string): "user-view-password.html" -> yes
	*/
	return c.Render(http.StatusOK, "user-view-password.html", echo.Map{
		"session":      session,
		"csrf":         c.Get("csrf"),
		"name":         fmt.Sprintf("User: %s", user.Name),
		"user":         user,
		"is_html_only": true,
	})
}

/*
 * Delete User ID
 *
 * @target: Users
 * @method: GET
 * @route: /users/delete/:id
 */
func (controller *Controller) DeleteUser(c echo.Context) error {
	session, _ := middleware.GetAuth(c)
	log := log.WithFields(log.Fields{
		"username": session.Values["username"],
		"route":    fmt.Sprintf("%v -> id:%v", c.Path(), c.Param("id")),
	})
	if session.Values["is_auth_type"] == -1 {
		log.Warn("for GET to delete user without no-session [@route: /login]")
		return c.Redirect(http.StatusFound, "/login")
	}

	log.Info("START request method GET for delete user")
	id, _ := strconv.Atoi(c.Param("id"))

	user, err := (models.User{}).FirstByID(controller.DB, id)
	if err != nil {
		log.Warnf("for GET to delete user without models.User{}.FirstByID() errors: `%v`", err)
		log.Warn("END request method GET for delete user: [-]failure")
		// HTTP response status: 405 Method Not Allowed
		return c.HTML(http.StatusNotAcceptable, err.Error())
	}

	// (models.User{}) or var user models.User or user := models.User{}
	if err := (models.User{}).Delete(controller.DB, id); err != nil {
		log.Warnf("for GET to delete user without models.User{}.Delete() errors: `%v`", err)
		log.Warn("END request method GET for delete user: [-]failure")
		// HTTP response status: 405 Method Not Allowed
		return c.HTML(http.StatusNotAcceptable, err.Error())
	}

	middleware.SetFlashSuccess(c, fmt.Sprintf("success delete user: %s!", user.Username))
	log.Info("END request method GET for delete user: [+]success")
	return c.Redirect(http.StatusMovedPermanently, "/users")
}
