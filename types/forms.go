package types

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

/*
 * type userForm: of a user
 *
 * @method: POST
 * @controller: CreateUser
 *				(user_controller.go)
 * @route: /users/add
 */
type UserForm struct {
	Username        string
	Email           string
	Password        string
	ConfirmPassword string
	Name            string
	City            uint
	Photo           string
}

/*
 * type PasswordForm: of a username and password
 *
 * @method: POST
 * @controller: Login
 * 				(session_controller.go)
 * @route: /login
 */
type PasswordForm struct {
	Username string
	Password string
}

// (type PasswordForm) Validate: of a validate username and password
func (lf PasswordForm) Validate() error {
	return validation.ValidateStruct(&lf,
		validation.Field(&lf.Username, validation.Required, validation.Length(4, 15)),
		validation.Field(&lf.Password, validation.Required, validation.Length(6, 18)),
	)
}

/*
 * type NewPasswordForm: of a password user
 *
 * @method: POST
 * @controller: UpdateUserByPassword
 * 				(user_controller.go)
 * @route: /login
 */
type NewPasswordForm struct {
	OldPassword        string // nothing
	NewPassword        string
	ConfirmNewPassword string
}

/* function PasswordEquals: of password equals confirm password
-----
var PasswordEqual = func(...) ... {
	...
}

equals,

func PasswordEquals(...) ... {
	...
}
*/
var PasswordEquals = func(confirm_password string) validation.RuleFunc {
	return func(value interface{}) error {
		password, _ := value.(string)
		if password != confirm_password {
			return errors.New("passwords don't match")
		}
		return nil
	}
}
