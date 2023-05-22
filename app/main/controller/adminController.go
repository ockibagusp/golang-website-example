package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/ockibagusp/golang-website-example/app/main/helpers"
	"github.com/ockibagusp/golang-website-example/app/main/middleware"
	selectTemplate "github.com/ockibagusp/golang-website-example/app/main/template"
	"github.com/ockibagusp/golang-website-example/business"
	selectUser "github.com/ockibagusp/golang-website-example/business/user"
	log "github.com/ockibagusp/golang-website-example/logger"
)

var aclogger = log.NewPackage("admin_controller")

func init() {
	// Templates: adminController
	templates := selectTemplate.AppendTemplates
	templates["admin/admin-delete-permanently.html"] = selectTemplate.ParseFilesBase("views/admin/admin-delete-permanently.html")
}

/*
 * Delete Permanently
 *
 * @target: [Admin] Delete Permanently
 * @method: GET
 * @route: /admin/delete-permanently
 */
func (ctrl *Controller) DeletePermanently(c echo.Context) error {
	log := aclogger.Start(c)
	defer log.End()
	log.Info("START request method GET for admin delete permanently")

	role, _ := c.Get("role").(string)
	if role != "admin" {
		log.Warn("for GET to admin delete permanently by id without no-session or user no admin [@route: /admin/delete/permanently/:id]")
		log.Warn("END request method GET for admin delete permanently: [-]failure")
		// HTTP response status: 404 Not Found
		return c.JSON(http.StatusNotFound, helpers.ResponseError{
			Code:   http.StatusNotFound,
			Status: "Not Found",
		})
	}

	var (
		users *[]selectUser.User
		err   error

		// typing: all, admin and user
		typing string
	)

	ic := business.InternalContext{}
	if c.QueryParam("admin") == "all" {
		log.Infof(`for GET to admin delete permanently: admin ctrl.userService.FindDeleteAll(db, "admin")`)
		typing = "Admin"
		users, err = ctrl.userService.FindDeleteAll(ic, "admin")
	} else if c.QueryParam("user") == "all" {
		log.Infof(`for GET to admin delete permanently: user ctrl.userService.FindDeleteAll(db, "user")`)
		typing = "User"
		users, err = ctrl.userService.FindDeleteAll(ic, "user")
	} else {
		log.Infof(`for GET to admin delete permanently: ctrl.userService.FindDeleteAll(db) or ctrl.userService.FindDeleteAll(db, "all")`)
		typing = "All"
		// models.User{} or (models.User{}) or var user models.User or user := models.User{}
		users, err = ctrl.userService.FindDeleteAll(ic)
	}

	if err != nil {
		log.Warnf("for GET to admin delete permanently without ctrl.userService.FindDeleteAll() errors: `%v`", err)
		log.Warn("END request method GET for admin delete permanently: [-]failure")
		// HTTP response status: 404 Not Found
		return c.JSON(http.StatusNotFound, helpers.ResponseError{
			Code:    http.StatusNotFound,
			Status:  "Not Found",
			Message: err.Error(),
		})
	}

	username, _ := c.Get("username").(string)
	log.Info("END request method GET to admin delete permanently: [+]success")
	return c.Render(http.StatusOK, "admin/admin-delete-permanently.html", echo.Map{
		"name":            fmt.Sprintf("Users: %v", typing),
		"nav":             "users", // (?)
		"claims_username": username,
		"claims_role":     role,
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
 * Restore User
 *
 * @target: [Admin] Restore User
 * @method: GET
 * @route: /admin/restore/:id
 */
func (ctrl *Controller) RestoreUser(c echo.Context) error {
	log := aclogger.Start(c)
	defer log.End()
	log.Info("START request method GET for admin restore user")

	role, _ := c.Get("role").(string)
	if role != "admin" {
		log.Warn("for GET to admin restore user without no-session or user no admin [@route: /admin/delete/permanently/:id]")
		log.Warn("END request method GET for admin restore user: [-]failure")
		// HTTP response status: 404 Not Found
		return c.JSON(http.StatusNotFound, helpers.ResponseError{
			Code:   http.StatusNotFound,
			Status: "Not Found",
		})
	}

	id, _ := strconv.Atoi(c.Param("id"))
	uid := uint(id)
	// why?
	// delete permanently not for admin
	if uid == 1 {
		log.Warn("END request method GET for admin restore [admin]: [-]failure")
		// HTTP response status: 403 Forbidden
		return c.JSON(http.StatusForbidden, helpers.ResponseError{
			Code:   http.StatusForbidden,
			Status: "Forbidden",
		})
	}

	trackerID := log.SetTrackerID()
	ic := business.NewInternalContext(trackerID)

	user, err := ctrl.userService.UnscopedFirstUserByID(ic, uid)
	if err != nil {
		log.Warnf("for GET to admin restore without ctrl.userService.FirstByID() errors: `%v`", err)
		log.Warn("END request method GET for admin restore: [-]failure")
		// HTTP response status: 404 Not Found
		return c.JSON(http.StatusNotFound, helpers.ResponseError{
			Code:    http.StatusNotFound,
			Status:  "Not Found",
			Message: err.Error(),
		})
	}

	if err := ctrl.userService.Restore(ic, uid); err != nil {
		log.Warnf("for GET to admin restore without ctrl.userService.Restore() errors: `%v`", err)
		log.Warn("END request method GET for admin restore: [-]failure")
		// HTTP response status: 403 Forbidden
		return c.JSON(http.StatusForbidden, helpers.ResponseError{
			Code:    http.StatusForbidden,
			Status:  "Forbidden",
			Message: err.Error(),
		})
	}

	middleware.SetFlashSuccess(c, fmt.Sprintf("success restore user: %s!", user.Username))

	// restore admin
	log.Info("END request method GET for admin restore: [+]success")
	return c.Redirect(http.StatusMovedPermanently, "/admin/delete-permanently")
}

/*
 * Delete Permanently By ID
 *
 * @target: [Admin] Delete Permanently By ID
 * @method: GET
 * @route: /admin/delete/permanently/:id
 */
func (ctrl *Controller) DeletePermanentlyByID(c echo.Context) error {
	log := aclogger.Start(c)
	defer log.End()
	log.Info("START request method GET for admin delete permanently by id")

	role, _ := c.Get("role").(string)
	if role != "admin" {
		log.Warn("for GET to admin delete permanently by id without no-session or user no admin [@route: /admin/delete/permanently/:id]")
		log.Warn("END request method GET for admin delete permanently by id: [-]failure")
		// HTTP response status: 404 Not Found
		return c.JSON(http.StatusNotFound, helpers.ResponseError{
			Code:   http.StatusNotFound,
			Status: "Not Found",
		})
	}

	id, _ := strconv.Atoi(c.Param("id"))
	uid := uint(id)
	// why?
	// delete permanently not for admin
	if uid == 1 {
		log.Warn("END request method GET for admin delete permanently by id [admin]: [-]failure")
		// HTTP response status: 403 Forbidden
		return c.JSON(http.StatusForbidden, helpers.ResponseError{
			Code:   http.StatusForbidden,
			Status: "Forbidden",
		})
	}

	trackerID := log.SetTrackerID()
	ic := business.NewInternalContext(trackerID)
	user, err := ctrl.userService.UnscopedFirstUserByID(ic, uid)
	if err != nil {
		log.Warnf("for GET to admin delete permanently by id without ctrl.userService.FirstByID() errors: `%v`", err)
		log.Warn("END request method GET for admin delete permanently by id: [-]failure")
		// HTTP response status: 404 Not Found
		return c.JSON(http.StatusNotFound, helpers.ResponseError{
			Code:    http.StatusNotFound,
			Status:  "Not Found",
			Message: err.Error(),
		})
	}

	if err = ctrl.userService.DeletePermanently(ic, uid); err != nil {
		log.Warnf("for GET to admin delete permanently by id without ctrl.userService.Delete() errors: `%v`", err)
		log.Warn("END request method GET for admin delete permanently by id: [-]failure")
		// HTTP response status: 403 Forbidden
		return c.JSON(http.StatusForbidden, helpers.ResponseError{
			Code:    http.StatusForbidden,
			Status:  "Forbidden",
			Message: err.Error(),
		})
	}

	middleware.SetFlashSuccess(c, fmt.Sprintf("success permanently user: %s!", user.Username))

	// delete permanently admin
	log.Info("END request method GET for admin delete permanently by id: [+]success")
	return c.Redirect(http.StatusMovedPermanently, "/admin/delete-permanently")
}
