package user

import "github.com/ockibagusp/golang-website-example/business"

type (
	service struct {
		repository Repository
	}

	Service interface {
		FindAll(ic business.InternalContext, role ...string) (selectedUsers *[]User, err error)
		FindByID(ic business.InternalContext, id int) (selectedUser *User, err error)
		FindByEmail(ic business.InternalContext, email string) (selectedUser *User, err error)
		Save(ic business.InternalContext, newUser *User) (*User, error)
		FirstUserByID(ic business.InternalContext, id int) (selectedUser *User, err error)
		FirstUserByUsername(ic business.InternalContext, username string) (selectedUser *User, err error)
		FirstByIDAndUsername(ic business.InternalContext, id int, username string, too ...bool) (selectedUser *User, err error)
		FirstByCityID(ic business.InternalContext, id int) (selectedUser *User, err error)
		Update(ic business.InternalContext, id int) (selectedUser *User, err error)
		UpdateByIDandPassword(ic business.InternalContext, id int, password string) (err error)
		Delete(ic business.InternalContext, id int) (err error)
		FindDeleteAll(ic business.InternalContext, role ...string) (selectedUsers *[]User, err error)
		Restore(ic business.InternalContext, id int) error
		DeletePermanently(ic business.InternalContext, id int) error
	}
)

func NewService(repository Repository) Service {
	return &service{
		repository,
	}
}

func (s *service) FindAll(ic business.InternalContext, role ...string) (selectedUsers *[]User, err error) {
	return s.repository.FindAll(ic, role...)
}

func (s *service) FindByID(ic business.InternalContext, id int) (selectedUser *User, err error) {
	return s.repository.FindByID(ic, id)
}

// login
func (s *service) FirstUserByUsername(ic business.InternalContext, username string) (selectedUser *User, err error) {
	return s.repository.FirstUserByUsername(ic, username)
}

func (s *service) FindByEmail(ic business.InternalContext, email string) (selectedUser *User, err error) {
	return s.repository.FindByEmail(ic, email)
}

func (s *service) Save(ic business.InternalContext, user *User) (selectedUser *User, err error) {
	return s.repository.Save(ic, user)
}

func (s *service) FirstUserByID(ic business.InternalContext, id int) (selectedUser *User, err error) {
	return s.repository.FirstUserByID(ic, id)
}

func (s *service) FirstByIDAndUsername(ic business.InternalContext, id int, username string, too ...bool) (selectedUser *User, err error) {
	return s.repository.FirstByIDAndUsername(ic, id, username, too...)
}

func (s *service) FirstByCityID(ic business.InternalContext, id int) (selectedUser *User, err error) {
	return s.repository.FirstByCityID(ic, id)
}

func (s *service) Update(ic business.InternalContext, id int) (selectedUser *User, err error) {
	return s.repository.Update(ic, id)
}

func (s *service) UpdateByIDandPassword(ic business.InternalContext, id int, password string) (err error) {
	return s.repository.UpdateByIDandPassword(ic, id, password)
}

func (s *service) Delete(ic business.InternalContext, id int) (err error) {
	return s.repository.Delete(ic, id)
}

func (s *service) FindDeleteAll(ic business.InternalContext, role ...string) (selectedUsers *[]User, err error) {
	return s.repository.FindDeleteAll(ic, role...)
}

func (s *service) Restore(ic business.InternalContext, id int) error {
	return s.repository.Restore(ic, id)
}

func (s *service) DeletePermanently(ic business.InternalContext, id int) error {
	return s.repository.DeletePermanently(ic, id)
}
