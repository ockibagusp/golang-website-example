package middleware

import (
	"net/http"
	"strings"

	"golang-website-example/business/auth"

	"golang-website-example/app/main/helpers"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func JwtAuthMiddleware(secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Initialize a new instance of `Claims` and ok
			var (
				claims *auth.JwtClaims = &auth.JwtClaims{}
				ok     bool
			)

			cookie, err := c.Request().Cookie("token")
			if err != nil || cookie.Value == "anonymous" {
				// http.ErrNoCookie or for any other type of error
				// for JWT claims anonymous
				claims.UserID = 0
				claims.Username = "anonymous"
				claims.Role = "anonymous"

				if err == http.ErrNoCookie {
					// If the cookie is not set, new the cookie
					SetCookieNoAuth(c)
				}
			} else if cookie.Value != "anonymous" {
				// Get the JWT string from the cookie
				tokenString := cookie.Value

				// Parse the JWT string and store the result in `claims`.
				// Note that we are passing the key in this method as well. This method will return an error
				// if the token is invalid (if it has expired according to the expiry time we set on sign in),
				// or if the signature does not match
				token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
					return []byte(secret), nil
				})
				if err != nil {
					if err == jwt.ErrSignatureInvalid {
						return c.JSON(http.StatusUnauthorized, helpers.ResponseUnauthorized)
					}
					return c.JSON(http.StatusBadRequest, helpers.ResponseBadRequest)
				}
				if !token.Valid {
					return c.JSON(http.StatusUnauthorized, helpers.ResponseUnauthorized)
				}

				claims, ok = token.Claims.(*auth.JwtClaims)
				if !ok && !token.Valid {
					return c.JSON(http.StatusForbidden, helpers.ResponseForbidden)
				}
			} else {
				// For any other type of error, return a bad request status
				return c.JSON(http.StatusBadRequest, helpers.ResponseBadRequest)
			}

			path := c.Request().URL.Path
			// -> role = "anonymous"
			if strings.Contains(path, "/login") || strings.Contains(path, "/logout") {
				return next(c)
			}

			c.Set("uid", claims.UserID)
			c.Set("username", claims.Username)
			c.Set("role", claims.Role)

			return next(c)
		}
	}
}

// https://medium.com/monstar-lab-bangladesh-engineering/jwt-auth-in-go-dde432440924
func isAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(*auth.JwtClaims)
		isAdmin := claims.Role == "admin"

		if isAdmin == false {
			return echo.ErrUnauthorized
		}
		return next(c)
	}
}
