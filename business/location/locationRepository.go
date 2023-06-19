package location

import "golang-website-example/business"

type Repository interface {
	FindAll(ic business.InternalContext) (selectedLocations *[]Location, err error)
	FirstByID(ic business.InternalContext, id int) (selectedLocation *Location, err error)
	Save(ic business.InternalContext) (selectedLocation *Location, err error)
	Delete(ic business.InternalContext, id int) (err error)
}
