package controller

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/ockibagusp/golang-website-example/app/main/middleware"
	selectTemplate "github.com/ockibagusp/golang-website-example/app/main/template"
	"github.com/ockibagusp/golang-website-example/app/main/types"

	"github.com/ockibagusp/golang-website-example/business"
	selectUser "github.com/ockibagusp/golang-website-example/business/user"
	log "github.com/ockibagusp/golang-website-example/logger"
	locationModules "github.com/ockibagusp/golang-website-example/modules/location"
)

var uclogger = log.NewPackage("user_controller")

func init() {
	// Templates: userController
	templates := selectTemplate.AppendTemplates
	templates["users/user-all.html"] = selectTemplate.ParseFilesBase("views/users/user-all.html")
	templates["users/user-add.html"] = selectTemplate.ParseFilesBase("views/users/user-add.html", "views/users/user-form.html")
	templates["users/user-read.html"] = selectTemplate.ParseFilesBase("views/users/user-read.html", "views/users/user-form.html")
	templates["users/user-view.html"] = selectTemplate.ParseFilesBase("views/users/user-view.html", "views/users/user-form.html")
	templates["user-view-password.html"] = selectTemplate.ParseFileHTMLOnly("views/users/user-view-password.html")
}

/*
 * Users All
 *
 * @target: Users
 * @method: GET
 * @route: /users
 */
func (ctrl *Controller) Users(c echo.Context) error {
	trackerID, log := uclogger.StartTrackerID(c)
	defer log.End()

	ic := business.NewInternalContext(trackerID)
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
		log.Info(`for GET to users admin ctrl.userService.FindAll(ic, "admin")`)
		typing = "Admin"
		users, err = ctrl.userService.FindAll(ic, "admin")
	} else if c.QueryParam("user") == "all" {
		log.Info(`for GET to users user ctrl.userService.FindAll(ic, "user")`)
		typing = "User"
		users, err = ctrl.userService.FindAll(ic, "user")
	} else {
		log.Info(`for GET to users ctrl.userService.FindAll(ic) or ctrl.userService.FindAll(ic, "all")`)
		typing = "All"
		users, err = ctrl.userService.FindAll(ic)
	}

	if err != nil {
		log.Warnf("for GET to users without ctrl.userService.FindAll errors: `%v`", err)
		log.Warn("END request method GET for users: [-]failure")
		// HTTP response status: 404 Not Found
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": err.Error(),
		})
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
	trackerID, log := uclogger.StartTrackerID(c)
	defer log.End()

	ic := business.NewInternalContext(trackerID)
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
		user, err = ctrl.userService.FirstUserByID(ic, uid)
		if err != nil {
			log.Warnf(`for GET for create user without select "id" where "username" errors: "%v"`, err)
			log.Warn("END request method GET for create user: [-]failure")
			return err
		}

		middleware.SetFlashError(c, "403 Forbidden")
		log.Infof("END request method GET for create user to users/read/%v: [-]failure", user.ID)
		return c.Redirect(http.StatusFound, fmt.Sprintf("/users/read/%v", user.ID))
	}

	locations, _ := locationModules.NewDB().FindAll(ic)
	if c.Request().Method == "POST" {
		log.Info("START request method POST for create user")

		// var location uint
		// if c.FormValue("location") != "" {
		// 	location64, err := strconv.ParseUint(c.FormValue("location"), 10, 32)
		// 	if err != nil {
		// 		log.Warnf("for POST to create user without location64 strconv.ParseUint() to error `%v`", err)
		// 		log.Warn("END request method POST for create user: [-]failure")
		// 		// HTTP response status: 400 Bad Request
		// 		return c.JSON(http.StatusBadRequest, echo.Map{
		// 			"message": err.Error(),
		// 		})
		// 	}
		// 	// Location or District?
		// 	// location = uint(location64)
		// }

		// var (
		// 	validPhoto         bool = true
		// 	errs               []error
		// 	fileType, fileName string
		// 	fileByte           []byte
		// 	userForm           types.UserForm
		// )

		// request on c: parse input, type multipart/form-data
		c.Request().ParseMultipartForm(1024)

		var (
			// validPhoto bool = false
			// fileType, fileName string
			// fileByte           []byte
			err error
		)

		// source
		file, err := c.FormFile("photo")
		if err != nil {
			return err
		}
		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()
		// > upload

		// fileByte, _ = ioutil.ReadAll(src)
		// fileType = http.DetectContentType(fileByte)

		// fileName = fmt.Sprint("members/", uuid.New().Time())
		// if fileType == "image/jpeg" {
		// 	fileName += ".jpeg"
		// } else if fileType == "image/png" {
		// 	fileName += ".png"
		// } else {
		// 	return errors.New("no image jpeg, jpg and png")
		// }

		file.Filename = fmt.Sprint("members/", uuid.New().Time(), ".jpeg")

		// Destination
		dst, err := os.Create(file.Filename)
		if err != nil {
			return err
		}
		defer dst.Close()

		// Copy
		if _, err = io.Copy(dst, src); err != nil {
			return err
		}

		return c.JSON(http.StatusOK, "ok")

		// userForm: type of a user
		// userForm = types.UserForm{
		// 	Role:            c.FormValue("role"),
		// 	Username:        c.FormValue("username"),
		// 	Email:           c.FormValue("email"),
		// 	Password:        c.FormValue("password"),
		// 	ConfirmPassword: c.FormValue("confirm_password"),
		// 	Name:            c.FormValue("name"),
		// 	Location:        location,
		// }

		return c.JSON(http.StatusOK, "errs")

		// // userForm: Validate of a validate user
		// err = validation.Errors{
		// 	// "username": validation.Validate(
		// 	// 	userForm.Username, validation.Required, validation.Length(4, 15),
		// 	// ),
		// 	// "email": validation.Validate(userForm.Email, validation.Required, validation.Length(5, 30), is.EmailFormat),
		// 	// "password": validation.Validate(
		// 	// 	userForm.Password, validation.Required, validation.Length(6, 18),
		// 	// 	validation.By(types.PasswordEquals(userForm.ConfirmPassword)),
		// 	// ),
		// 	// "name":     validation.Validate(userForm.Name, validation.Required, validation.Length(3, 30)),
		// 	// "location": validation.Validate(userForm.Location),
		// 	"photo": validation.Validate(userForm.Photo),
		// }.Filter()
		// /* if err = validation.Errors{...}.Filter(); err != nil {
		// 	...
		// } why?
		// */
		// if err != nil {
		// 	log.Warnf("for POST to create user without validation.Errors: `%v`", err)
		// 	middleware.SetFlashError(c, err.Error())

		// 	log.Warn("END request method POST for create user: [-]failure")
		// 	// HTTP response status: 400 Bad Request
		// 	return c.Render(http.StatusBadRequest, "users/user-add.html", echo.Map{
		// 		"name":             "User Add",
		// 		"nav":              "user Add", // (?)
		// 		"session_username": username,
		// 		"session_role":     role,
		// 		"flash_error 	":    middleware.GetFlashError(c),
		// 		"csrf":             c.Get("csrf"),
		// 		"locations":        locations,
		// 		"is_new":           true,
		// 	})
		// }

		// // Password Hash
		// var hash string
		// hash, err = ctrl.authService.PasswordHash(userForm.Password)
		// if err != nil {
		// 	log.Warnf("for POST to create user without middleware.PasswordHash error: `%v`", err)
		// 	log.Warn("END request method POST for create user: [-]failure")
		// 	return err
		// }

		// user = &selectUser.User{
		// 	Role:     userForm.Role,
		// 	Username: userForm.Username,
		// 	Email:    userForm.Email,
		// 	Password: hash,
		// 	Name:     userForm.Name,
		// 	Location: userForm.Location,
		// 	Photo:    userForm.Photo,
		// }

		// if _, err := ctrl.userService.Create(ic, user); err != nil {
		// 	log.WithField("user_failure", user).Warn("for POST to create user without models.User: nil")
		// 	middleware.SetFlashError(c, err.Error())

		// 	log.Warn("END request method POST for create user: [-]failure")
		// 	// HTTP response status: 400 Bad Request
		// 	return c.Render(http.StatusBadRequest, "users/user-add.html", echo.Map{
		// 		"name":             "User Add",
		// 		"nav":              "user Add", // (?)
		// 		"session_username": username,
		// 		"session_role":     role,
		// 		"csrf":             c.Get("csrf"),
		// 		"flash_error":      middleware.GetFlashError(c),
		// 		"locations":        locations,
		// 		"is_new":           true,
		// 	})
		// }

		// log.WithField("user_success", user).Info("models.User: [+]success")
		// middleware.SetFlashSuccess(c, fmt.Sprintf("success new user: %s!", user.Username))
		// // create user
		// if role == "anonymous" {
		// 	if _, err := middleware.SetSession(user, c); err != nil {
		// 		middleware.SetFlashError(c, err.Error())
		// 		log.Warn("to middleware.SetSession session not found for create user")
		// 		log.Warn("END request method POST for create user: [-]failure")
		// 		// err: session not found
		// 		return c.JSON(http.StatusForbidden, echo.Map{
		// 			"message": err.Error(),
		// 		})
		// 	}
		// 	log.Info("END request method POST for create user: [+]success")
		// 	return c.Redirect(http.StatusMovedPermanently, "/")
		// }
		// log.Info("END request method POST for create user: [+]success")
		// // create admin
		// return c.Redirect(http.StatusMovedPermanently, "/users")
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
	trackerID, log := uclogger.StartTrackerID(c)
	defer log.End()

	ic := business.NewInternalContext(trackerID)
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

	user, err := ctrl.userService.FirstUserByID(ic, uid)
	if err != nil {
		log.Warnf(
			"for GET to read user without models.User{}.FirstByID() errors: `%v`", err,
		)
		log.Warn("END request method GET for read user: [-]failure")
		// HTTP response status: 406 Method Not Acceptable
		return c.JSON(http.StatusNotAcceptable, echo.Map{
			"message": err.Error(),
		})
	}

	locations, err := locationModules.NewDB().FindAll(ic)
	if err != nil {
		log.Warnf("for GET to read user without models.location{}.FindAll() errors: `%v`", err)
		log.Warn("END request method GET for read user: [-]failure")
		// HTTP response status: 406 Not Acceptable
		return c.JSON(http.StatusNotAcceptable, echo.Map{
			"message": err.Error(),
		})
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
	trackerID, log := uclogger.StartTrackerID(c)
	defer log.End()

	ic := business.NewInternalContext(trackerID)
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
	user, err = ctrl.userService.FirstUserByID(ic, uid)
	if err != nil {
		log.Info("START request method GET/POST for update user")
		log.Warnf(
			"for GET to update user without models.User{}.FirstByID() errors: `%v`", err,
		)
		log.Warn("END request method GET for update user: [-]failure")
		// HTTP response status: 404 Not Found
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": err.Error(),
		})
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
		return c.JSON(http.StatusForbidden, echo.Map{
			"message": "Forbidden",
		})
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
				return c.JSON(http.StatusBadRequest, echo.Map{
					"message": err.Error(),
				})
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

		// newUser, err = ctrl.userService.Update(ic, user, updateUser); err != nil: equal
		if user, err = ctrl.userService.Update(ic, user, updateUser); err != nil {
			log.Warnf(
				"for POST to update user without models.User{}.Update() errors: `%v`", err,
			)
			middleware.SetFlashError(c, err.Error())
			log.Warn("END request method POST for update user: [-]failure")

			locations, _ := locationModules.NewDB().FindAll(ic)
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

		log.WithField("user_update", user).Info("models.User: [+]success", "user_update", user)
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

	locations, _ := locationModules.NewDB().FindAll(ic)
	if err != nil {
		log.Warnf("for GET to update user without models.location{}.FindAll() errors: `%v`", err)
		log.Warn("END request method GET for update user: [-]failure")
		// HTTP response status: 405 Method Not Allowed
		return c.JSON(http.StatusNotAcceptable, echo.Map{
			"message": err.Error(),
		})
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

/*
 * Update User ID by Password
 *
 * @target: Users
 * @method: GET or POST
 * @route: /users/view/:id/password
 */
func (ctrl *Controller) UpdateUserByPassword(c echo.Context) error {
	trackerID, log := uclogger.StartTrackerID(c)
	defer log.End()

	ic := business.NewInternalContext(trackerID)
	id, _ := strconv.Atoi(c.Param("id"))
	uid := uint(id)
	username, _ := c.Get("username").(string)
	role, _ := c.Get("role").(string)

	if role == "anonymous" {
		log.Warn("for GET to update user by password without no-session [@route: /login]")
		middleware.SetFlashError(c, "login process failed!")
		log.Warn("END request method GET for update user by password: [-]failure")
		return c.Redirect(http.StatusFound, "/login")
	}

	log.Info("START request method GET or POST for update user by password")

	var (
		user *selectUser.User
		err  error
	)
	user, err = ctrl.userService.FirstUserByID(ic, uid)
	if err != nil {
		log.Warnf(
			"for GET to update user by password without models.User{}.FirstByID() errors: `%v`", err,
		)
		log.Warn("END request method GET for update user by password: [-]failure")
		// HTTP response status: 404 Not Found
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": err.Error(),
		})
	}

	/*
		for example:
		username ockibagusp update by password 'ockibagusp': ok
		username ockibagusp update by password 'sugriwa': no
	*/
	_, err = ctrl.userService.FirstByIDAndUsername(
		ic, uid, username,
	)

	if role == "user" && err != nil {
		log.Warnf(
			"for GET to update user by password without models.User{}.FirstByIDAndUsername() errors: `%v`", err,
		)
		log.Warn("END request method GET for update user by password: [-]failure")
		// HTTP response status: 403 Forbidden
		return c.JSON(http.StatusForbidden, echo.Map{
			"message": err.Error(),
		})
	}

	if c.Request().Method == "POST" {
		// newPasswordForm: type of a password user
		_newPasswordForm := types.NewPasswordForm{
			OldPassword:        c.FormValue("old_password"),
			NewPassword:        c.FormValue("new_password"),
			ConfirmNewPassword: c.FormValue("confirm_new_password"),
		}

		if !ctrl.authService.CheckHashPassword(user.Password, _newPasswordForm.OldPassword) {
			log.Warn("for POST to update user by password without not check hash password: 403 Forbidden")
			middleware.SetFlashError(c, "check hash password is wrong!")
			log.Warn("END request method POST for update user by password: [-]failure")
			return c.Render(http.StatusForbidden, "user-view-password.html", echo.Map{
				"name":             fmt.Sprintf("User: %s", user.Name),
				"session_username": username,
				"session_role":     role,
				"flash_error":      middleware.GetFlashError(c),
				"user":             user,
				"is_html_only":     true,
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
				"name":             fmt.Sprintf("User: %s", user.Name),
				"session_username": username,
				"session_role":     role,
				"flash_error":      middleware.GetFlashError(c),
				"user":             user,
				"is_html_only":     true,
			})
		}

		// Password Hash
		hash, err := ctrl.authService.PasswordHash(_newPasswordForm.NewPassword)
		if err != nil {
			log.Warnf("for POST to update user by password without middleware.PasswordHash() errors: `%v`", err)
			log.Warn("END request method POST for update user by password: [-]failure")
			return err
		}

		// err := ctrl.userService.UpdateByIDandPassword(...): be able
		if err := ctrl.userService.UpdateByIDandPassword(ic, uid, hash); err != nil {
			log.Warnf("for POST to update user by password without models.User{}.UpdateByIDandPassword() errors: `%v`", err)
			log.Warn("END request method POST for update user by password: [-]failure")
			// HTTP response status: 405 Method Not Allowed
			return c.JSON(http.StatusNotAcceptable, echo.Map{
				"message": err.Error(),
			})
		}

		log.WithField("user_update_password", user).Info("models.User: [+]success")
		middleware.SetFlashSuccess(c, fmt.Sprintf("success update user by password: %s!", user.Username))
		if role == "user" {
			log.Info("END [user] request method POST for update user by password: [+]success")
			// update user by password
			return c.Redirect(http.StatusMovedPermanently, "/")
		}
		log.Info("END [admin] request method POST for update user by password: [+]success")
		// update user by password [admin]
		return c.Redirect(http.StatusMovedPermanently, "/users")
	}

	// admin
	if user == nil {
		user, _ = ctrl.userService.FirstUserByID(ic, uid)
	}

	log.Info("END request method GET for update user by password: [+]success")
	/*
		name (string): "users/user-view-password.html" -> no
			{..,"status":500,"error":"html/template: \"users/user-view-password.html\" is undefined",..}
			why?
		name (string): "user-view-password.html" -> yes
	*/
	return c.Render(http.StatusOK, "user-view-password.html", echo.Map{
		"session_username": username,
		"session_role":     role,
		"csrf":             c.Get("csrf"),
		"name":             fmt.Sprintf("User: %s", user.Name),
		"user":             user,
		"is_html_only":     true,
	})
}

/*
 * Delete User ID
 *
 * @target: Users
 * @method: GET
 * @route: /users/delete/:id
 */
func (ctrl *Controller) DeleteUser(c echo.Context) error {
	trackerID, log := uclogger.StartTrackerID(c)
	defer log.End()

	ic := business.NewInternalContext(trackerID)
	role, _ := c.Get("role").(string)
	if role == "anonymous" {
		log.Warn("for GET to delete user without no-session [@route: /login]")
		middleware.SetFlashError(c, "login process failed!")
		log.Warn("END request method GET for delete user: [-]failure")
		return c.Redirect(http.StatusFound, "/login")
	}

	log.Info("START request method GET for delete user")
	id, _ := strconv.Atoi(c.Param("id"))
	uid := uint(id)

	// why?
	// delete not for admin
	if uid == 1 {
		log.Warn("END request method GET for delete user [admin]: [-]failure")
		// HTTP response status: 403 Forbidden
		return c.JSON(http.StatusForbidden, echo.Map{
			"message": "Forbidden",
		})
	}

	var (
		user *selectUser.User
		err  error
	)
	user, err = ctrl.userService.FirstUserByID(ic, uid)
	if err != nil {
		log.Warnf("for GET to delete user without models.User{}.FirstByID() errors: `%v`", err)
		log.Warn("END request method GET for delete user: [-]failure")
		// HTTP response status: 404 Not Found
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": err.Error(),
		})
	}

	/*
		TODO:
		for example:
		username ockibagusp delete 'ockibagusp': ok
		username ockibagusp delete 'sugriwa': no
		insyaallah
	*/
	oldUsername := user.Username
	_, err = ctrl.userService.FirstByIDAndUsername(
		ic, uid, oldUsername,
	)

	if !(role == "admin") {
		if err != nil {
			log.Warnf(
				"for GET to delete without models.User{}.FirstByIDAndUsername() errors: `%v`", err,
			)
			log.Warn("END request method GET for delete user: [-]failure")
			// HTTP response status: 403 Forbidden
			return c.JSON(http.StatusForbidden, echo.Map{
				"message": err.Error(),
			})
		}
	}

	if err := ctrl.userService.Delete(ic, uid); err != nil {
		log.Warnf("for GET to delete user without models.User{}.Delete() errors: `%v`", err)
		log.Warn("END request method GET for delete user: [-]failure")
		// HTTP response status: 403 Forbidden
		return c.JSON(http.StatusForbidden, echo.Map{
			"message": err.Error(),
		})
	}

	middleware.SetFlashSuccess(c, fmt.Sprintf("success delete user: %s!", user.Username))
	if role == "user" {
		log.Info("END [user] request method GET for delete user: [+]success")
		if err := middleware.ClearSession(c); err != nil {
			log.Warn("to middleware.ClearSession session not found")
			// err: session not found
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": err.Error(),
			})
		}
		// delete user
		return c.Redirect(http.StatusSeeOther, "/")
	}
	log.Info("END request method GET for delete user: [+]success")
	// delete admin
	return c.Redirect(http.StatusMovedPermanently, "/users")
}
