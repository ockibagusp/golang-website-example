package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ockibagusp/golang-website-example/business"
	"github.com/ockibagusp/golang-website-example/business/user"
	"golang.org/x/crypto/bcrypt"
)

type (
	JwtClaims struct {
		UserID   uint   `form:"user_id" json:"user_id"`
		Username string `form:"username" json:"username"`
		Role     string `form:"role" json:"role"`

		jwt.RegisteredClaims
	}

	service struct {
		userService user.Service
	}

	Service interface {
		VerifyLogin(ic business.InternalContext, email string, plainPassword string) (getUser *user.User, validPassword bool)
		CheckHashPassword(hash, password string) bool
		PasswordHash(password string) (string, error)
	}
)

func NewService(userService user.Service) Service {
	return &service{
		userService,
	}
}

// VerifyLogin: get-user and valid-password
func (s *service) VerifyLogin(ic business.InternalContext, username string, plainPassword string) (getUser *user.User, validPassword bool) {
	getUser, err := s.userService.FirstUserByUsername(ic, username)
	if err != nil {
		return
	}

	if !s.CheckHashPassword(getUser.Password, plainPassword) {
		return
	}

	validPassword = true

	return
}

/*
password is slow: yes
----------------
password := "..."
hash, err := PasswordHash(password)
if err != nil {
	return err
}
// match = true
// match = false
if !CheckHashPassword(hash, password) {
	return ...
}
then faster: no
-----------
password := "..."
hash := sha256.New()
hash.Write([]byte(password))
sha_hash := hex.EncodeToString(hash.Sum(nil))

fmt.Println("Password -> ", password)
fmt.Println("Hash -> ", sha_hash)
*/

// PasswordHash: hash for password to string
func (s *service) PasswordHash(password string) (string, error) {
	// GenerateFromPassword(..., bcrypt.DefaultCost{=10})
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

// CheckHashPassword: hashes for passwords to bool
func (s *service) CheckHashPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func newJWTClaims(userID uint, username string, role string, issuedAt time.Time, expiredAt time.Time) JwtClaims {
	return JwtClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiredAt),
			IssuedAt:  jwt.NewNumericDate(issuedAt),
			NotBefore: jwt.NewNumericDate(issuedAt),
		},
	}
}

func GenerateToken(jwtAuhtSign string, userID uint, userName string, userRole string) (signedToken string, err error) {
	timeNow := time.Now()
	claims := newJWTClaims(userID, userName, userRole, timeNow, timeNow.Add(time.Hour*1))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err = token.SignedString([]byte(jwtAuhtSign))
	if err != nil {
		return
	}

	return signedToken, nil
}
