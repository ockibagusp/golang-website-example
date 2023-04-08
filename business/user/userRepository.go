package user

import "github.com/ockibagusp/golang-website-example/business"

type Repository interface {
	FindAll(ic business.InternalContext, role ...string) (selectedUsers *[]User, err error)
	FindByID(ic business.InternalContext, uid uint) (selectedUser *User, err error)
	FindByEmail(ic business.InternalContext, email string) (selectedUser *User, err error)
	Create(ic business.InternalContext, newUser *User) (*User, error)
	CreatesBatch(ic business.InternalContext, newUsers *[]User) (*[]User, error)
	FirstUserByID(ic business.InternalContext, uid uint) (selectedUser *User, err error)
	FirstUserByUsername(ic business.InternalContext, username string) (selectedUser *User, err error)
	FirstByIDAndUsername(ic business.InternalContext, uid uint, username string, too ...bool) (selectedUser *User, err error)
	FirstByCityID(ic business.InternalContext, uid uint) (selectedUser *User, err error)
	UnscopedFirstUserByID(ic business.InternalContext, uid uint) (selectedUser *User, err error)
	Update(ic business.InternalContext, oldUser *User, updateUser *User) (*User, error)
	UpdateByIDandPassword(ic business.InternalContext, uid uint, password string) (err error)
	Delete(ic business.InternalContext, uid uint) (err error)
	FindDeleteAll(ic business.InternalContext, role ...string) (selectedUsers *[]User, err error)
	Restore(ic business.InternalContext, uid uint) error
	DeletePermanently(ic business.InternalContext, uid uint) error
}
