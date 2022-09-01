package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/ockibagusp/golang-website-example/middleware"
	"github.com/ockibagusp/golang-website-example/models"
	log "github.com/sirupsen/logrus"
)

/*
 * Delete Permanently
 *
 * @target: [Admin] Delete Permanently
 * @method: GET
 * @route: /admin/delete-permanently
 */
func (controller *Controller) DeletePermanently(c echo.Context) error {
	session, _ := middleware.GetAuth(c)
	log := log.WithFields(log.Fields{
		"username": session.Values["username"],
		"route":    c.Path(),
	})
	log.Info("START request method GET for admin delete permanently")

	is_auth_type := session.Values["is_auth_type"]
	if is_auth_type == -1 {
		log.Warn("for GET to admin delete permanently without no-session [@route: /login]")
		middleware.SetFlashError(c, "login process failed!")
		log.Warn("END request method GET for admin delete permanently: [-]failure")
		return c.Redirect(http.StatusFound, "/login")
	}

	if !middleware.IsAdmin(is_auth_type) {
		log.Warn("END request method GET for admin delete permanently: [-]failure")
		// HTTP response status: 404 Not Found
		return c.HTML(http.StatusNotFound, "404 Not Found")
	}

	var (
		users []models.User
		err   error

		// typing: all, admin and user
		typing string
	)

	if c.QueryParam("admin") == "all" {
		log.Infof(`for GET to admin delete permanently: admin models.User{}.FindDeleteAll(db, "admin")`)
		typing = "Admin"
		users, err = models.User{}.FindDeleteAll(controller.DB, "admin")
	} else if c.QueryParam("user") == "all" {
		log.Infof(`for GET to admin delete permanently: user models.User{}.FindDeleteAll(db, "user")`)
		typing = "User"
		users, err = models.User{}.FindDeleteAll(controller.DB, "user")
	} else {
		log.Infof(`for GET to admin delete permanently: models.User{}.FindDeleteAll(db) or models.User{}.FindDeleteAll(db, "all")`)
		typing = "All"
		// models.User{} or (models.User{}) or var user models.User or user := models.User{}
		users, err = models.User{}.FindDeleteAll(controller.DB)
	}

	if err != nil {
		log.Warnf("for GET to admin delete permanently without models.User{}.FindDeleteAll() errors: `%v`", err)
		log.Warn("END request method GET for admin delete permanently: [-]failure")
		// HTTP response status: 404 Not Found
		return c.HTML(http.StatusNotFound, err.Error())
	}

	log.Info("END request method GET to admin delete permanently: [+]success")
	return c.Render(http.StatusOK, "admin/admin-delete-permanently.html", echo.Map{
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
			"success": middleware.GetFlashSuccess(c),
			"error":   middleware.GetFlashError(c),
		},
		"users": users,
	})
}

/*
 * Delete Permanently By ID
 *
 * @target: [Admin] Delete Permanently By ID
 * @method: GET
 * @route: /admin/delete/permanently/:id
 */
func (controller *Controller) DeletePermanentlyByID(c echo.Context) error {
	session, _ := middleware.GetAuth(c)
	log := log.WithFields(log.Fields{
		"username": session.Values["username"],
		"route":    c.Path(),
	})
	log.Info("START request method GET for admin delete permanently by id")

	is_auth_type := session.Values["is_auth_type"]
	if is_auth_type == -1 {
		log.Warn("for GET to admin delete permanently by id without no-session [@route: /login]")
		middleware.SetFlashError(c, "login process failed!")
		log.Warn("END request method GET for admin delete permanently by id: [-]failure")
		return c.Redirect(http.StatusFound, "/login")
	}

	id, _ := strconv.Atoi(c.Param("id"))

	// why?
	// delete permanently not for admin
	if id == 1 {
		log.Warn("END request method GET for admin delete permanently by id [admin]: [-]failure")
		// HTTP response status: 403 Forbidden
		return c.HTML(http.StatusForbidden, "403 Forbidden")
	}

	if !middleware.IsAdmin(is_auth_type) {
		log.Warn("END request method GET for admin delete permanently by id: [-]failure")
		// HTTP response status: 404 Not Found
		return c.HTML(http.StatusNotFound, "404 Not Found")
	}

	user, err := (models.User{}).UnscopedFirstUserByID(controller.DB, id)
	if err != nil {
		log.Warnf("for GET to admin delete permanently by id without models.User{}.FirstByID() errors: `%v`", err)
		log.Warn("END request method GET for admin delete permanently by id: [-]failure")
		// HTTP response status: 404 Not Found
		return c.HTML(http.StatusNotFound, err.Error())
	}

	if err := user.DeletePermanently(controller.DB, id); err != nil {
		log.Warnf("for GET to admin delete permanently by id without models.User{}.Delete() errors: `%v`", err)
		log.Warn("END request method GET for admin delete permanently by id: [-]failure")
		// HTTP response status: 403 Forbidden
		return c.HTML(http.StatusForbidden, err.Error())
	}

	// delete permanently admin
	log.Info("END request method GET for admin delete permanently by id: [+]success")
	return c.Redirect(http.StatusMovedPermanently, "/admin/delete-permanently")
}

/*
 * Restore User
 *
 * @target: [Admin] Restore User
 * @method: GET
 * @route: /admin/restore/:id
 */
func (controller *Controller) RestoreUser(c echo.Context) error {
	session, _ := middleware.GetAuth(c)
	log := log.WithFields(log.Fields{
		"username": session.Values["username"],
		"route":    c.Path(),
	})
	log.Info("START request method GET for admin restore")

	is_auth_type := session.Values["is_auth_type"]
	if is_auth_type == -1 {
		log.Warn("for GET to admin restore without no-session [@route: /login]")
		middleware.SetFlashError(c, "login process failed!")
		log.Warn("END request method GET for admin restore: [-]failure")
		return c.Redirect(http.StatusFound, "/login")
	}

	id, _ := strconv.Atoi(c.Param("id"))

	// why?
	// delete permanently not for admin
	if id == 1 {
		log.Warn("END request method GET for admin restore [admin]: [-]failure")
		// HTTP response status: 403 Forbidden
		return c.HTML(http.StatusForbidden, "403 Forbidden")
	}

	if !middleware.IsAdmin(is_auth_type) {
		log.Warn("END request method GET for admin restore: [-]failure")
		// HTTP response status: 404 Not Found
		return c.HTML(http.StatusNotFound, "404 Not Found")
	}

	user, err := (models.User{}).UnscopedFirstUserByID(controller.DB, id)
	if err != nil {
		log.Warnf("for GET to admin restore without models.User{}.FirstByID() errors: `%v`", err)
		log.Warn("END request method GET for admin restore: [-]failure")
		// HTTP response status: 404 Not Found
		return c.HTML(http.StatusNotFound, err.Error())
	}

	if err := user.Restore(controller.DB, id); err != nil {
		log.Warnf("for GET to admin restore without models.User{}.Restore() errors: `%v`", err)
		log.Warn("END request method GET for admin restore: [-]failure")
		// HTTP response status: 403 Forbidden
		return c.HTML(http.StatusForbidden, err.Error())
	}

	middleware.SetFlashSuccess(c, fmt.Sprintf("success restore user: %s!", user.Username))

	// restore admin
	log.Info("END request method GET for admin restore: [+]success")
	return c.Redirect(http.StatusMovedPermanently, "/admin/delete-permanently")
}
