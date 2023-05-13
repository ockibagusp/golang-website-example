package middleware

import (
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/ockibagusp/golang-website-example/business/auth"
	selectUser "github.com/ockibagusp/golang-website-example/business/user"
	"github.com/ockibagusp/golang-website-example/config"
)

// SetCookieNoAuth: set cookie from User: anonymous
func SetCookieNoAuth(c echo.Context) (err error) {
	JWTAuthSign := config.GetAPPConfig().AppJWTAuthSign

	// Finally, we set the client cookie for "token" as the anonymous just generated
	// Declare the expiration time of the token
	// here, we have kept it as 24 hour
	expirationTime := time.Now().Add(24 * time.Hour)

	// Create claims with multiple fields populated
	claims := auth.JwtClaims{
		UserID:   0,
		Username: "anonymous",
		Role:     "anonymous",
		RegisteredClaims: jwt.RegisteredClaims{
			// A usual scenario is to set the expiration time relative to the current time
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString([]byte(JWTAuthSign))
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		return errors.New("server error")
	}

	// Finally, we set the client cookie for "token" as the JWT we just generated
	// we also set an expiry time which is the same as the token itself
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = tokenString
	cookie.Expires = expirationTime
	c.SetCookie(cookie)

	return
}

// SetCookie: set cookie from User
func SetCookie(c echo.Context, user *selectUser.User, JWTAuthSign string) (err error) {
	// Declare the expiration time of the token
	// here, we have kept it as 24 hour
	expirationTime := time.Now().Add(24 * time.Hour)

	// Create claims with multiple fields populated
	claims := auth.JwtClaims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			// A usual scenario is to set the expiration time relative to the current time
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString([]byte(JWTAuthSign))
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		return errors.New("server error")
	}

	// Finally, we set the client cookie for "token" as the JWT we just generated
	// we also set an expiry time which is the same as the token itself
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = tokenString
	cookie.Expires = expirationTime
	c.SetCookie(cookie)

	return
}

// ClearCookie: delete cookie from User
func ClearCookie(c echo.Context) {
	c.SetCookie(&http.Cookie{
		Name:   "token",
		Value:  "anonymous",
		MaxAge: -1,
	})
}
