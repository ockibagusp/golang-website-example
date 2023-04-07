package user

import (
	"errors"

	"github.com/ockibagusp/golang-website-example/business"
	"github.com/ockibagusp/golang-website-example/business/user"

	"gorm.io/gorm"
)

type (
	GormRepository struct {
		*gorm.DB
	}
)

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{
		db.Table("users"),
	}
}

// User: FindAll(ic, role={admin, user})
func (repo *GormRepository) FindAll(ic business.InternalContext, role ...string) (selectedUsers *[]user.User, err error) {
	query := repo.DB.WithContext(ic.ToContext())

	// equal,
	// if len(role) == 0 || len(role) == 1 && role[0] == "all" {...}
	if isAll(&role) {
		// Limit: 50 ?
		err = query.Limit(50).Find(&selectedUsers).Error
	} else if isAdmin(&role) {
		err = query.Limit(50).Where(`role = "admin"`).Find(&selectedUsers).Error
	} else if isUser(&role) {
		err = query.Limit(50).Where(`role = "user"`).Find(&selectedUsers).Error
	} else { // role agrs [2,..]=string
		return nil, errors.New(`models.User{}.FirstAll: role agrs [2]{"admin", "user"}=string`)
	}

	if err != nil {
		return selectedUsers, err
	}

	return
}

// User: FirstByID
func (repo *GormRepository) FindByID(ic business.InternalContext, uid uint) (selectedUser *user.User, err error) {
	query := repo.DB.WithContext(ic.ToContext())

	err = query.First(&selectedUser, uid).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New("User Not Found")
			return selectedUser, err
		}
		return selectedUser, err
	}

	return
}

func (repo *GormRepository) FindByEmail(ic business.InternalContext, email string) (selectedUser *user.User, err error) {
	query := repo.DB.WithContext(ic.ToContext())
	if err := query.Where("email = ?", email).Find(&selectedUser).Error; err != nil {
		return selectedUser, err
	}

	return
}

// User: Create
func (repo *GormRepository) Create(ic business.InternalContext, newUser *user.User) (*user.User, error) {
	query := repo.DB.WithContext(ic.ToContext())
	if err := query.Create(&newUser).Error; err != nil {
		return nil, err
	}

	return newUser, nil
}

// User: CreatesBatch
func (repo *GormRepository) CreatesBatch(ic business.InternalContext, newUsers *[]user.User) (*[]user.User, error) {
	query := repo.DB.WithContext(ic.ToContext())

	tx := query.Begin()
	// if tx.Create(&newUsers).Error != nil {}
	if err := tx.Create(&newUsers).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()

	return newUsers, nil
}

// User: FirstUserByID
func (repo *GormRepository) FirstUserByID(ic business.InternalContext, uid uint) (selectedUser *user.User, err error) {
	query := repo.DB.WithContext(ic.ToContext())

	err = query.First(&selectedUser, uid).Error
	return isFirstUserByID(selectedUser, err)
}

// User: UnscopedFirstUserByID
func (repo *GormRepository) UnscopedFirstUserByID(ic business.InternalContext, uid uint) (selectedUser *user.User, err error) {
	query := repo.DB.WithContext(ic.ToContext())

	err = query.Unscoped().First(&selectedUser, uid).Error

	return isFirstUserByID(selectedUser, err)
}

// User: isFirstUserByID
func isFirstUserByID(user *user.User, err error) (*user.User, error) {
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, errors.New("User Not Found")
		}
		return user, err
	}

	return user, nil
}

// User: FirstUserByUsername -> login
func (repo *GormRepository) FirstUserByUsername(ic business.InternalContext, username string) (selectedUser *user.User, err error) {
	query := repo.DB.WithContext(ic.ToContext())

	if err = query.Where(
		"username = ?", username,
	).First(&selectedUser).Error; err != nil {
		return nil, err
	}
	return
}

// User: FirstByIDAndUsername
//
// example:
// user, err := models.User{}.FirstByIDAndUsername(1, "ockibagusp")
//
// or,
//
// user, err := models.User{}.FirstByIDAndUsername(1, "ockibagusp", true)
func (repo *GormRepository) FirstByIDAndUsername(ic business.InternalContext, uid uint, username string, too ...bool) (selectedUser *user.User, err error) {
	query := repo.DB.WithContext(ic.ToContext())

	if len(too) == 0 {
		err = query.Select("id", "username", "password").
			Where("username = ?", username).First(&selectedUser, uid).Error
	} else if len(too) == 1 {
		err = query.Where("username = ?", username).First(&selectedUser, uid).Error
	} else { // too agrs [2,..]=bool
		return selectedUser, errors.New("models.User{}.FirstByIDAndUsername: too agrs [0, 1]=bool")
	}

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return selectedUser, errors.New("User Not Found")
		}
		return selectedUser, err
	}

	return
}

// User: FirstByCityID
func (repo *GormRepository) FirstByCityID(ic business.InternalContext, uid uint) (selectedUser *user.User, err error) {
	query := repo.DB.WithContext(ic.ToContext())

	err = query.Select("users.*, cities.id as city_id, cities.city as city_massage").
		Joins("left join cities on users.city = cities.id").First(&selectedUser, uid).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return selectedUser, errors.New("User Not Found")
		}
		return selectedUser, err
	}

	return
}

// User: Update
func (repo *GormRepository) Update(ic business.InternalContext, newUser *user.User, updateUser *user.User) (*user.User, error) {
	query := repo.DB.WithContext(ic.ToContext())

	err := query.Model(&newUser).Updates(&user.User{
		Role:     updateUser.Role,
		Username: updateUser.Username,
		Email:    updateUser.Email,
		Name:     updateUser.Name,
		Location: updateUser.Location,
		Photo:    updateUser.Photo,
	}).Error
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

// User: Update By uID and Password
func (repo *GormRepository) UpdateByIDandPassword(ic business.InternalContext, uid uint, password string) (err error) {
	query := repo.DB.WithContext(ic.ToContext())

	selectedUser := user.User{}
	if err = query.Model(&selectedUser).Where("id = ?", uid).Update("password", password).First(&selectedUser).Error; err != nil {
		return err
	}

	return
}

// User: Delete
func (repo *GormRepository) Delete(ic business.InternalContext, uid uint) (err error) {
	query := repo.DB.WithContext(ic.ToContext())

	selectedUser := user.User{}

	tx := query.Begin()
	var count int64
	// if tx.Select("id").First(&selectUser).Error != nil {}
	if tx.Select("id").First(&selectedUser, uid).Count(&count); count != 1 {
		tx.Rollback()
		return errors.New("User Not Found")
	}

	// if tx.Delete(&selectUser, uid).Error != nil {}
	if err := tx.Delete(&selectedUser, uid).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}

// User: FindDeleteAll(db, role={admin, user})
func (repo *GormRepository) FindDeleteAll(ic business.InternalContext, role ...string) (selectedUsers *[]user.User, err error) {
	query := repo.DB.WithContext(ic.ToContext())

	// equal,
	// if len(role) == 0 || len(role) == 1 && role[0] == "all" {...}
	if isAll(&role) {
		// Limit: 50 ?
		err = query.Limit(50).Unscoped().Where("deleted_at is not null").Find(&selectedUsers).Error
	} else if isAdmin(&role) {
		err = query.Limit(50).Unscoped().Where(`role = "admin" AND deleted_at is not null`).Find(&selectedUsers).Error
	} else if isUser(&role) {
		err = query.Limit(50).Unscoped().Where(`role = "user" AND deleted_at is not null`).Find(&selectedUsers).Error
	} else { // role agrs [2,..]=string
		return nil, errors.New(`models.User{}.FindDeleteAll: role agrs [2]{"admin", "user"}=string`)
	}

	if err != nil {
		return nil, err
	}

	return
}

// User: Restore
func (repo *GormRepository) Restore(ic business.InternalContext, uid uint) error {
	query := repo.DB.WithContext(ic.ToContext())

	selectedUser := user.User{}

	tx := query.Begin()
	var count int64
	if tx.Model(&selectedUser).Unscoped().Select("deleted_at").Where("id = ?", uid).First(&selectedUser).Count(&count); count != 1 {
		tx.Rollback()
		return errors.New("Recover Undeleted User")
	}

	// if tx.Model(&selectUser).Unscoped().Where("id = ?", uid).Update(...).Count; count != 1 {}
	if tx.Model(&selectedUser).Unscoped().Where("id = ?", uid).Update("deleted_at", nil).First(&selectedUser).Count(&count); count != 1 {
		tx.Rollback()
		return errors.New("Restore User Update Not deleted_at")
	}
	tx.Commit()

	return nil
}

// User: Delete Permanently
func (repo *GormRepository) DeletePermanently(ic business.InternalContext, uid uint) error {
	query := repo.DB.WithContext(ic.ToContext())

	selectedUser := user.User{}

	tx := query.Begin()
	var count int64
	// if tx.Unscoped().Select("id").First(&selectUser).Error != nil {}
	if tx.Unscoped().Select("id").First(&selectedUser).Count(&count); count != 1 {
		tx.Rollback()
		return errors.New("User Not Found")
	}

	// if tx.Unscoped().Delete(&selectUser, uid).Error != nil {}
	if err := tx.Unscoped().Delete(&selectedUser, uid).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}

// is? all, admin or user?
func isAll(role *[]string) bool {
	return len(*role) == 0 || len(*role) == 1 && (*role)[0] == "all"
}

func isAdmin(role *[]string) bool {
	return len(*role) == 1 && (*role)[0] == "admin"
}

func isUser(role *[]string) bool {
	return len(*role) == 1 && (*role)[0] == "user"
}
