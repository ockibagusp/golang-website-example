package user

import "github.com/ockibagusp/golang-website-example/business"

type Repository interface {
	FindByIDandVersion(ic business.InternalContext, id, version int) (selectedUser User, err error)

	FindByEmail(ic business.InternalContext, email string) (selectedUser User, err error)
}
