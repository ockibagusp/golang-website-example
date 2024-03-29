package middleware

import (
	"errors"
	"net/http"
	"time"

	"golang-website-example/business/auth"
	selectUser "golang-website-example/business/user"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// SetCookieNoAuth: set cookie from User: anonymous
func SetCookieNoAuth(c echo.Context) {
	cookie := setTokenAnonymous()
	c.SetCookie(cookie)
}

// SetCookie: set cookie from User
func SetCookie(c echo.Context, user *selectUser.User, jwtAuthSign string) (err error) {
	// Declare the expiration time of the token
	// here, we have kept it as 24 hour
	expiredAt := time.Now().Add(time.Hour * 24)
	issuedAt := time.Now()

	// Create claims with multiple fields populated
	claims := auth.JwtClaims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			// A usual scenario is to set the expiration time relative to the current time
			ExpiresAt: jwt.NewNumericDate(expiredAt),
			IssuedAt:  jwt.NewNumericDate(issuedAt),
			NotBefore: jwt.NewNumericDate(issuedAt),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	signedToken, err := token.SignedString([]byte(jwtAuthSign))
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		return errors.New("server error")
	}

	// Finally, we set the client cookie for "token" as the JWT we just generated
	// we also set an expiry time which is the same as the token itself
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = signedToken
	cookie.Expires = expiredAt
	c.SetCookie(cookie)

	return
}

// ClearCookie: delete cookie from User
func ClearCookie(c echo.Context) {
	c.SetCookie(setTokenAnonymous())
}

func setTokenAnonymous() (cookie *http.Cookie) {
	cookie = new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = "anonymous"
	// Declare the expiration time of the token
	// here, we have kept it as 7 days
	cookie.Expires = time.Now().Add(7 * (24 * time.Hour))
	// // cakes are all missing
	// cookie.MaxAge = -1

	return
}
