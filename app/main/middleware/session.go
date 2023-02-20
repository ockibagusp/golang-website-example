package middleware

import (
	"net/http"
	"strings"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/ockibagusp/golang-website-example/config"
)

func SessionNewCookieStore() echo.MiddlewareFunc {
	sessionsCookieStore := config.GetAPPConfig().SessionsCookieStore

	return session.Middleware(sessions.NewCookieStore(
		[]byte(sessionsCookieStore),
	))
}

func SessionMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			path := c.Request().URL.Path
			if strings.Contains(path, "/login") || strings.Contains(path, "/logout") {
				return next(c)
			}

			session_gorilla, err := session.Get("session", c)
			if err != nil {
				return c.HTML(http.StatusForbidden, "no session")
			}

			role := session_gorilla.Values["role"]
			if role != "admin" && role != "user" {
				role = "anonymous"
			}

			username := session_gorilla.Values["username"]
			if username != "" {
				username = "anonymous"
			}

			c.Set("username", username)
			c.Set("role", role)

			return next(c)
		}
	}
}
