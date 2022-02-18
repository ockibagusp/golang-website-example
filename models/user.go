package models

import (
	"errors"

	"gorm.io/gorm"
)

// User: struct
type User struct {
	Model
	// database: just `username` varchar 15
	Username string `gorm:"unique;not null;type:varchar(15)" form:"username"`
	Email    string `gorm:"unique;not null" form:"email"`
	Password string `gorm:"not null" form:"password"`
	Name     string `gorm:"not null" form:"name"`
	City     uint   `form:"city"`
	Photo    string `form:"photo"`
	IsAdmin  uint16 `form:"is_admin"`
}

// UserCity: struct
type UserCity struct {
	User
	CityMassage string `gorm:"index:city_massage" form:"city_massage"`
}

/*
 * the tx *DB: exists apparetly
 */

// User: Save
func (user User) Save(db *gorm.DB) (*User, error) {
	if err := db.Create(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// User: FindAll
func (User) FindAll(db *gorm.DB) ([]User, error) {
	users := []User{}

	// Limit: 25 ?
	err := db.Limit(25).Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}

// User: FirstByID
func (user User) FirstByID(db *gorm.DB, id int) (*User, error) {
	err := db.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("User Not Found")
		}
		return nil, err
	}

	return &user, nil
}

// User: FirstByIDAndUsername
//
// example:
// user, err := models.User{}.FirstByIDAndUsername(controllers.DB, 1, "ockibagusp")
//
// or,
//
// user, err := models.User{}.FirstByIDAndUsername(controllers.DB, 1, "ockibagusp", true)
func (user User) FirstByIDAndUsername(db *gorm.DB, id int, username string, too ...bool) (*User, error) {
	var err error
	if len(too) == 0 {
		err = db.Select("id", "username", "password").
			Where("username = ?", username).First(&user, id).Error
	} else if len(too) == 1 {
		err = db.Where("username = ?", username).First(&user, id).Error
	} else { // too agrs [2,..]=bool
		return nil, errors.New("models.User{}.FirstByIDAndUsername: too agrs [0, 1]=bool")
	}

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("User Not Found")
		}
		return nil, err
	}

	return &user, nil
}

// User: FirstByCityID
func (user User) FirstByCityID(db *gorm.DB, id int) (*User, error) {
	err := db.Select("users.*, cities.id as city_id, cities.city as city_massage").
		Joins("left join cities on users.city = cities.id").First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("User Not Found")
		}
		return nil, err
	}

	return &user, nil
}

// User: Update
func (user User) Update(db *gorm.DB, id int) (*User, error) {
	err := db.Where("id = ?", id).Updates(&User{
		Username: user.Username,
		Email:    user.Email,
		Name:     user.Name,
		City:     user.City,
		Photo:    user.Photo,
	}).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// User: Update By ID and Password
func (user User) UpdateByIDandPassword(db *gorm.DB, id int, password string) (err error) {
	if err = db.Model(&user).Where("id = ?", id).Update("password", password).First(&user).Error; err != nil {
		return err
	}

	return
}

// User: Delete
func (user User) Delete(db *gorm.DB, id int) error {
	tx := db.Begin()
	var count int64
	// if tx.Select("id").First(&user).Error != nil {}
	if tx.Select("id").First(&user).Count(&count); count != 1 {
		tx.Rollback()
		return errors.New("User Not Found")
	}

	// if tx.Delete(&user, id).Error != nil {}
	if err := tx.Delete(&user, id).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}
