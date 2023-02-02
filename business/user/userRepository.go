package user

import "github.com/ockibagusp/golang-website-example/business"

type Repository interface {
	FindAll(ic business.InternalContext, role ...string) (selectedUsers []User, err error)
	FindByID(ic business.InternalContext, id int) (selectedUser User, err error)
	FindByEmail(ic business.InternalContext, email string) (selectedUser User, err error)
	Save(ic business.InternalContext) (selectedUser User, err error)
	FirstUserByID(ic business.InternalContext, id int) (selectedUser User, err error)
	FirstByIDAndUsername(ic business.InternalContext, id int, username string, too ...bool) (selectedUser User, err error)
	FirstByCityID(ic business.InternalContext, id int) (selectedUser User, err error)
	Update(ic business.InternalContext, id int) (selectedUser User, err error)
	UpdateByIDandPassword(ic business.InternalContext, id int, password string) (err error)
	Delete(ic business.InternalContext, id int) (err error)
	FindDeleteAll(ic business.InternalContext, role ...string) (selectedUsers []User, err error)
	Restore(ic business.InternalContext, id int) error
	DeletePermanently(ic business.InternalContext, id int) error
}
