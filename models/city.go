package models

import (
	"errors"

	"gorm.io/gorm"
)

// City: struct
type City struct {
	ID   uint   `form:"id"`
	City string `form:"city"`
}

// TableName name: string
func (City) TableName() string {
	return "cities"
}

// City: Save
func (city City) Save(db *gorm.DB) (City, error) {
	if err := db.Create(&city).Error; err != nil {
		return City{}, err
	}

	return city, nil
}

// City: FindAll
func (City) FindAll(db *gorm.DB) ([]City, error) {
	cities := []City{}

	if err := db.Find(&cities).Error; err != nil {
		return []City{}, err
	}

	return cities, nil
}

// City: FirstByID
func (city City) FirstByID(db *gorm.DB, id int) (City, error) {
	if err := db.First(&city, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return City{}, errors.New("City Not Found")
		}
		return City{}, err
	}

	return city, nil
}

// City: Delete
func (city City) Delete(db *gorm.DB, id int) error {
	tx := db.Begin()
	var count int64
	// if tx.Select("id").First(&city).Error != nil {}
	if tx.Select("id").First(&city).Count(&count); count != 1 {
		tx.Rollback()
		return errors.New("City Not Found")
	}
	// if tx.Delete(&city, id).Error != nil {}
	if err := tx.Delete(&city, id).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}
