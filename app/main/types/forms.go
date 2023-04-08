package types

import (
	"encoding/json"
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
	Role            string `form:"role"`
	Username        string `form:"username"`
	Email           string `form:"email"`
	Password        string `form:"password"`
	ConfirmPassword string `form:"confirm_password"`
	Name            string `form:"name"`
	Location        uint   `form:"location"`
	Photo           string `form:"photo"`
}

func (userForm UserForm) MarshalJSON() ([]byte, error) {
	type oldUF UserForm
	redactUF := oldUF(userForm)
	redactUF.Password = "[REDACTED]"

	return json.Marshal((*oldUF)(&redactUF))
}

/*
 * type LoginForm: of a username and password
 *
 * @method: POST
 * @controller: Login
 * 				(session_controller.go)
 * @route: /login
 */
type LoginForm struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

func (lf LoginForm) MarshalJSON() ([]byte, error) {
	type oldLF LoginForm
	redactLF := oldLF(lf)
	redactLF.Password = "[REDACTED]"

	return json.Marshal((*oldLF)(&redactLF))
}

// (type PasswordForm) Validate: of a validate username and password
func (lf LoginForm) Validate() error {
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
	OldPassword        string `form:"old_password"`
	NewPassword        string `form:"new_password"`
	ConfirmNewPassword string `form:"confirm_new_password"`
}

func (npf NewPasswordForm) MarshalJSON() ([]byte, error) {
	type oldNPF NewPasswordForm
	redactNPF := oldNPF(npf)
	redactNPF.OldPassword = "[REDACTED]"
	redactNPF.NewPassword = "[REDACTED]"

	return json.Marshal((*oldNPF)(&redactNPF))
}

/*
	function PasswordEquals: of password equals confirm password

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
