package city

import "github.com/ockibagusp/golang-website-example/business"

type (
	service struct {
		repository Repository
	}

	Service interface {
		FindAll(ic business.InternalContext) (selectedCities *[]City, err error)
		FirstByID(ic business.InternalContext, id int) (selectedCity *City, err error)
		Save(ic business.InternalContext) (selectedCity *City, err error)
		Delete(ic business.InternalContext, id int) (err error)
	}
)

func NewService(repository Repository) Service {
	return &service{
		repository,
	}
}

func (s *service) FindAll(ic business.InternalContext) (selectedCities *[]City, err error) {
	return s.repository.FindAll(ic)
}

func (s *service) FirstByID(ic business.InternalContext, id int) (selectedCity *City, err error) {
	return s.repository.FirstByID(ic, id)
}

func (s *service) Save(ic business.InternalContext) (selectedCity *City, err error) {
	return s.repository.Save(ic)
}

func (s *service) Delete(ic business.InternalContext, id int) (err error) {
	return s.repository.Delete(ic, id)
}
