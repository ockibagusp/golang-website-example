package city

import "github.com/ockibagusp/golang-website-example/business"

type Repository interface {
	FindAll(ic business.InternalContext) (selectedCities *[]City, err error)
	FirstByID(ic business.InternalContext, id int) (selectedCity *City, err error)
	Save(ic business.InternalContext) (selectedCity *City, err error)
	Delete(ic business.InternalContext, id int) (err error)
}
