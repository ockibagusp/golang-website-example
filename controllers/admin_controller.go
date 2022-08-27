package controllers

import (
	"fmt"
	"net/http"

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
		log.Warn("for GET to users without no-session [@route: /login]")
		middleware.SetFlashError(c, "login process failed!")
		log.Warn("END request method GET for users: [-]failure")
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
		log.Infof(`for GET to admin delete permanently: admin models.User{}.FindAll(db, "admin")`)
		typing = "Admin"
		users, err = models.User{}.FindAll(controller.DB, "admin")
	} else if c.QueryParam("user") == "all" {
		log.Infof(`for GET to admin delete permanently: user models.User{}.FindAll(db, "user")`)
		typing = "User"
		users, err = models.User{}.FindAll(controller.DB, "user")
	} else {
		log.Infof(`for GET to admin delete permanently: models.User{}.FindAll(db) or models.User{}.FindAll(db, "all")`)
		typing = "All"
		// models.User{} or (models.User{}) or var user models.User or user := models.User{}
		users, err = models.User{}.FindAll(controller.DB)
	}

	if err != nil {
		log.Warnf("for GET to admin delete permanently without models.User{}.FindAll() errors: `%v`", err)
		log.Warn("END request method GET for admin delete permanently: [-]failure")
		// HTTP response status: 404 Not Found
		return c.HTML(http.StatusNotFound, err.Error())
	}

	log.Info("END request method GET for users: [+]success")
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
