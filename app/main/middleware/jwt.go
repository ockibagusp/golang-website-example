package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type (
	Response struct {
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}
)

var responseForbidden = echo.Map{
	"message": http.StatusText(http.StatusForbidden),
}

func JwtAuthMiddleware(secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if strings.Contains(c.Request().URL.Path, "/login") {
				return next(c)
			}

			signature := strings.Split(c.Request().Header.Get("Authorization"), " ")
			if len(signature) < 2 {
				return c.JSON(http.StatusForbidden, responseForbidden)
			}
			if signature[0] != "Bearer" {
				return c.JSON(http.StatusForbidden, responseForbidden)
			}

			claim := jwt.MapClaims{}
			token, err := jwt.ParseWithClaims(signature[1], claim, func(token *jwt.Token) (interface{}, error) {
				_, ok := token.Method.(*jwt.SigningMethodHMAC)
				if !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}

				return []byte(secret), nil
			})
			if err != nil {
				return c.JSON(http.StatusForbidden, responseForbidden)
			}

			method, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok || method != jwt.SigningMethodHS256 {
				return c.JSON(http.StatusForbidden, responseForbidden)
			}

			if claim.Valid() != nil {
				return c.JSON(http.StatusForbidden, responseForbidden)
			}

			userID, _ := claim["user_id"].(float64)
			username, _ := claim["username"].(string)
			role, _ := claim["role"].(string)
			c.Set("user_id", int(userID))
			c.Set("username", username)
			c.Set("role", role)

			return next(c)
		}
	}
}

// https://medium.com/monstar-lab-bangladesh-engineering/jwt-auth-in-go-dde432440924
func isAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		isAdmin := claims["admin"].(bool)

		if isAdmin == false {
			return echo.ErrUnauthorized
		}
		return next(c)
	}
}
