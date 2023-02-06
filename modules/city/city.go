package city

import (
	"errors"

	"github.com/ockibagusp/golang-website-example/business"
	selectedCity "github.com/ockibagusp/golang-website-example/business/city"
	"gorm.io/gorm"
)

type (
	GormRepository struct {
		*gorm.DB
	}
)

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{
		db.Table("cities"),
	}
}

// City: FindAll(ic)
func (repo *GormRepository) FindAll(ic business.InternalContext) (selectedCities *[]selectedCity.City, err error) {
	query := repo.DB.WithContext(ic.ToContext())

	if err := query.Find(&selectedCities).Error; err != nil {
		return nil, err
	}

	return
}

// City: FirstByID(ic, id)
func (repo *GormRepository) FirstByID(ic business.InternalContext, id int) (selectedCity *selectedCity.City, err error) {
	query := repo.DB.WithContext(ic.ToContext())

	if err := query.First(&selectedCity, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("City Not Found")
		}
		return nil, err
	}

	return
}

// City: Save(ic)
func (repo *GormRepository) Save(ic business.InternalContext) (selectedCity *selectedCity.City, err error) {
	query := repo.DB.WithContext(ic.ToContext())

	if err := query.Create(&selectedCity).Error; err != nil {
		return nil, err
	}

	return
}

// City: Delete(ic, id)
func (repo *GormRepository) Delete(ic business.InternalContext, id int) (err error) {
	query := repo.DB.WithContext(ic.ToContext())

	selectedCity := selectedCity.City{}

	tx := query.Begin()
	var count int64
	// if tx.Select("id").First(&city).Error != nil {}
	if tx.Select("id").First(&selectedCity).Count(&count); count != 1 {
		tx.Rollback()
		return errors.New("City Not Found")
	}
	// if tx.Delete(&city, id).Error != nil {}
	if err := tx.Delete(&selectedCity, id).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return
}
