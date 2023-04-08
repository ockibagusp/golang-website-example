package location

import "github.com/ockibagusp/golang-website-example/business"

type (
	service struct {
		repository Repository
	}

	Service interface {
		FindAll(ic business.InternalContext) (selectedLocations *[]Location, err error)
		FirstByID(ic business.InternalContext, id int) (selectedLocation *Location, err error)
		Save(ic business.InternalContext) (selectedLocation *Location, err error)
		Delete(ic business.InternalContext, id int) (err error)
	}
)

func NewService(repository Repository) Service {
	return &service{
		repository,
	}
}

func (s *service) FindAll(ic business.InternalContext) (selectedLocations *[]Location, err error) {
	return s.repository.FindAll(ic)
}

func (s *service) FirstByID(ic business.InternalContext, id int) (selectedLocation *Location, err error) {
	return s.repository.FirstByID(ic, id)
}

func (s *service) Save(ic business.InternalContext) (selectedLocation *Location, err error) {
	return s.repository.Save(ic)
}

func (s *service) Delete(ic business.InternalContext, id int) (err error) {
	return s.repository.Delete(ic, id)
}
