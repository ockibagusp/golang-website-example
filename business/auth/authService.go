package auth

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/ockibagusp/golang-website-example/business"
	"github.com/ockibagusp/golang-website-example/business/user"
	"golang.org/x/crypto/bcrypt"
)

type (
	JwtClaims struct {
		UserID   uint   `json:"user_id"`
		Username string `json:"username"`
		Role     string `json:"role"`

		jwt.StandardClaims
	}

	service struct {
		userService user.Service
	}

	Service interface {
		VerifyLogin(ic business.InternalContext, email string, plainPassword string) (getUser *user.User, validPassword bool)
		GenerateToken(jwtSign string, userID uint, userName string, userRole string) (signedToken string, err error)
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
func (s *service) VerifyLogin(ic business.InternalContext, email string, plainPassword string) (getUser *user.User, validPassword bool) {
	getUser, err := s.userService.FindByEmail(ic, email)
	if err != nil {
		return
	}

	// // or,
	// passwordHash, err := s.PasswordHash(plainPassword)
	passwordHash, err := func(plainPassword string) (string, error) {
		// GenerateFromPassword(..., bcrypt.DefaultCost{=10})
		hash, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
		return string(hash), err
	}(plainPassword)
	if err != nil {
		return
	}

	if passwordHash != getUser.Password {
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

func newJWTClaims(userID uint, username string, role string, issuedAt int64, expiredAt int64) JwtClaims {
	return JwtClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  issuedAt,
			ExpiresAt: expiredAt,
		},
	}
}

func (s *service) GenerateToken(jwtSign string, userID uint, userName string, userRole string) (signedToken string, err error) {
	timeNow := time.Now()
	claims := newJWTClaims(userID, userName, userRole, timeNow.Unix(), timeNow.Add(time.Hour*24).Unix())
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err = token.SignedString([]byte(jwtSign))
	if err != nil {
		return
	}

	return signedToken, nil
}
