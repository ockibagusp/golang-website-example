package user

import "github.com/ockibagusp/golang-website-example/business"

type (
	service struct {
		repository Repository
	}

	Service interface {
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
)

func NewService(repository Repository) Service {
	return &service{
		repository,
	}
}

func (s *service) FindAll(ic business.InternalContext, role ...string) (selectedUsers *[]User, err error) {
	return s.repository.FindAll(ic, role...)
}

func (s *service) FindByID(ic business.InternalContext, uid uint) (selectedUser *User, err error) {
	return s.repository.FindByID(ic, uid)
}

func (s *service) UnscopedFirstUserByID(ic business.InternalContext, uid uint) (*User, error) {
	return s.repository.UnscopedFirstUserByID(ic, uid)
}

// login
func (s *service) FirstUserByUsername(ic business.InternalContext, username string) (selectedUser *User, err error) {
	return s.repository.FirstUserByUsername(ic, username)
}

func (s *service) FindByEmail(ic business.InternalContext, email string) (selectedUser *User, err error) {
	return s.repository.FindByEmail(ic, email)
}

func (s *service) Create(ic business.InternalContext, newUser *User) (selectedUser *User, err error) {
	return s.repository.Create(ic, newUser)
}

func (s *service) CreatesBatch(ic business.InternalContext, newUsers *[]User) (*[]User, error) {
	return s.repository.CreatesBatch(ic, newUsers)
}

func (s *service) FirstUserByID(ic business.InternalContext, uid uint) (selectedUser *User, err error) {
	return s.repository.FirstUserByID(ic, uid)
}

func (s *service) FirstByIDAndUsername(ic business.InternalContext, uid uint, username string, too ...bool) (selectedUser *User, err error) {
	return s.repository.FirstByIDAndUsername(ic, uid, username, too...)
}

func (s *service) FirstByCityID(ic business.InternalContext, uid uint) (selectedUser *User, err error) {
	return s.repository.FirstByCityID(ic, uid)
}

func (s *service) Update(ic business.InternalContext, user *User, updateUser *User) (*User, error) {
	return s.repository.Update(ic, user, updateUser)
}

func (s *service) UpdateByIDandPassword(ic business.InternalContext, uid uint, password string) (err error) {
	return s.repository.UpdateByIDandPassword(ic, uid, password)
}

func (s *service) Delete(ic business.InternalContext, uid uint) (err error) {
	return s.repository.Delete(ic, uid)
}

func (s *service) FindDeleteAll(ic business.InternalContext, role ...string) (selectedUsers *[]User, err error) {
	return s.repository.FindDeleteAll(ic, role...)
}

func (s *service) Restore(ic business.InternalContext, uid uint) error {
	return s.repository.Restore(ic, uid)
}

func (s *service) DeletePermanently(ic business.InternalContext, uid uint) error {
	return s.repository.DeletePermanently(ic, uid)
}
