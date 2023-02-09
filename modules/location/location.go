package location

import (
	"errors"

	"github.com/ockibagusp/golang-website-example/business"
	selectedLocation "github.com/ockibagusp/golang-website-example/business/location"
	"github.com/ockibagusp/golang-website-example/config"
	"gorm.io/gorm"
)

type (
	GormRepository struct {
		*gorm.DB
	}
)

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{
		db.Table("locations"),
	}
}

func NewServiceDB() selectedLocation.Service {
	conf := config.GetAPPConfig()
	db := conf.GetDatabaseConnection()
	locationRepo := NewGormRepository(db)

	return selectedLocation.NewService(locationRepo)
}

// Location: FindAll(ic)
func (repo *GormRepository) FindAll(ic business.InternalContext) (selectedLocations *[]selectedLocation.Location, err error) {
	query := repo.DB.WithContext(ic.ToContext())

	if err := query.Find(&selectedLocations).Error; err != nil {
		return nil, err
	}

	return
}

// Location: FirstByID(ic, id)
func (repo *GormRepository) FirstByID(ic business.InternalContext, id int) (selectedLocation *selectedLocation.Location, err error) {
	query := repo.DB.WithContext(ic.ToContext())

	if err := query.First(&selectedLocation, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Location Not Found")
		}
		return nil, err
	}

	return
}

// Location: Save(ic)
func (repo *GormRepository) Save(ic business.InternalContext) (selectedLocation *selectedLocation.Location, err error) {
	query := repo.DB.WithContext(ic.ToContext())

	if err := query.Create(&selectedLocation).Error; err != nil {
		return nil, err
	}

	return
}

// Location: Delete(ic, id)
func (repo *GormRepository) Delete(ic business.InternalContext, id int) (err error) {
	query := repo.DB.WithContext(ic.ToContext())

	selectedLocation := selectedLocation.Location{}

	tx := query.Begin()
	var count int64
	// if tx.Select("id").First(&location).Error != nil {}
	if tx.Select("id").First(&selectedLocation).Count(&count); count != 1 {
		tx.Rollback()
		return errors.New("Location Not Found")
	}
	// if tx.Delete(&location, id).Error != nil {}
	if err := tx.Delete(&selectedLocation, id).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return
}
