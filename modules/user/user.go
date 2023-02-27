package user

import (
	"errors"

	"github.com/ockibagusp/golang-website-example/business"
	selectUser "github.com/ockibagusp/golang-website-example/business/user"
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
func (repo *GormRepository) FindAll(ic business.InternalContext, role ...string) (selectedUsers *[]selectUser.User, err error) {
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
func (repo *GormRepository) FindByID(ic business.InternalContext, id int) (selectedUser *selectUser.User, err error) {
	query := repo.DB.WithContext(ic.ToContext())

	err = query.First(&selectedUser, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New("User Not Found")
			return selectedUser, err
		}
		return selectedUser, err
	}

	return
}

func (repo *GormRepository) FindByEmail(ic business.InternalContext, email string) (selectedUser *selectUser.User, err error) {
	query := repo.DB.WithContext(ic.ToContext())
	if err := query.Where("email = ?", email).Find(&selectedUser).Error; err != nil {
		return selectedUser, err
	}

	return
}

// User: Save
func (repo *GormRepository) Save(ic business.InternalContext, newUser *selectUser.User) (*selectUser.User, error) {
	query := repo.DB.WithContext(ic.ToContext())
	if err := query.Create(&newUser).Error; err != nil {
		return nil, err
	}

	return newUser, nil
}

// User: FirstUserByID
func (repo *GormRepository) FirstUserByID(ic business.InternalContext, id int) (selectedUser *selectUser.User, err error) {
	query := repo.DB.WithContext(ic.ToContext())

	err = query.First(&selectedUser, id).Error
	return isFirstUserByID(selectedUser, err)
}

// User: isFirstUserByID
func isFirstUserByID(user *selectUser.User, err error) (*selectUser.User, error) {
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, errors.New("User Not Found")
		}
		return user, err
	}

	return user, nil
}

// User: FirstUserByUsername -> login
func (repo *GormRepository) FirstUserByUsername(ic business.InternalContext, username string) (selectedUser *selectUser.User, err error) {
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
func (repo *GormRepository) FirstByIDAndUsername(ic business.InternalContext, id int, username string, too ...bool) (selectedUser *selectUser.User, err error) {
	query := repo.DB.WithContext(ic.ToContext())

	if len(too) == 0 {
		err = query.Select("id", "username", "password").
			Where("username = ?", username).First(&selectedUser, id).Error
	} else if len(too) == 1 {
		err = query.Where("username = ?", username).First(&selectedUser, id).Error
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
func (repo *GormRepository) FirstByCityID(ic business.InternalContext, id int) (selectedUser *selectUser.User, err error) {
	query := repo.DB.WithContext(ic.ToContext())

	err = query.Select("users.*, cities.id as city_id, cities.city as city_massage").
		Joins("left join cities on users.city = cities.id").First(&selectedUser, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return selectedUser, errors.New("User Not Found")
		}
		return selectedUser, err
	}

	return
}

// User: Update
func (repo *GormRepository) Update(ic business.InternalContext, id int) (selectedUser *selectUser.User, err error) {
	query := repo.DB.WithContext(ic.ToContext())

	err = query.Where("id = ?", id).Updates(&selectUser.User{
		Username: selectedUser.Username,
		Email:    selectedUser.Email,
		Name:     selectedUser.Name,
		Location: selectedUser.Location,
		Photo:    selectedUser.Photo,
	}).Error
	if err != nil {
		return selectedUser, err
	}

	return
}

// User: Update By ID and Password
func (repo *GormRepository) UpdateByIDandPassword(ic business.InternalContext, id int, password string) (err error) {
	query := repo.DB.WithContext(ic.ToContext())

	selectedUser := selectUser.User{}
	if err = query.Model(&selectedUser).Where("id = ?", id).Update("password", password).First(&selectedUser).Error; err != nil {
		return err
	}

	return
}

// User: Delete
func (repo *GormRepository) Delete(ic business.InternalContext, id int) (err error) {
	query := repo.DB.WithContext(ic.ToContext())

	selectedUser := selectUser.User{}

	tx := query.Begin()
	var count int64
	// if tx.Select("id").First(&user).Error != nil {}
	if tx.Select("id").First(&selectedUser).Count(&count); count != 1 {
		tx.Rollback()
		return errors.New("User Not Found")
	}

	// if tx.Delete(&user, id).Error != nil {}
	if err := tx.Delete(&selectedUser, id).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}

// User: FindDeleteAll(db, role={admin, user})
func (repo *GormRepository) FindDeleteAll(ic business.InternalContext, role ...string) (selectedUsers *[]selectUser.User, err error) {
	query := repo.DB.WithContext(ic.ToContext())

	// equal,
	// if len(role) == 0 || len(role) == 1 && role[0] == "all" {...}
	if isAll(&role) {
		// Limit: 50 ?
		err = query.Limit(50).Unscoped().Where("deleted_at is not null").Find(&selectedUsers).Error
	} else if isAdmin(&role) {
		err = query.Limit(50).Unscoped().Where("is_admin = 1 AND deleted_at is not null").Find(&selectedUsers).Error
	} else if isUser(&role) {
		err = query.Limit(50).Unscoped().Where("is_admin = 0 AND deleted_at is not null").Find(&selectedUsers).Error
	} else { // role agrs [2,..]=string
		return nil, errors.New(`models.User{}.FindDeleteAll: role agrs [2]{"admin", "user"}=string`)
	}

	if err != nil {
		return nil, err
	}

	return
}

// User: Restore
func (repo *GormRepository) Restore(ic business.InternalContext, id int) error {
	query := repo.DB.WithContext(ic.ToContext())

	selectedUser := selectUser.User{}

	tx := query.Begin()
	var count int64
	// if tx.Unscoped().Select("id").First(&user).Error != nil {}
	if tx.Unscoped().Select("id").First(&selectedUser).Count(&count); count != 1 {
		tx.Rollback()
		return errors.New("User Not Found")
	}

	// if tx.Model(&user).Unscoped().Where("id = ?", id).Update(...).Error; err != nil {}
	if err := tx.Model(&selectedUser).Unscoped().Where("id = ?", id).Update("deleted_at", nil).First(&selectedUser).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}

// User: Delete Permanently
func (repo *GormRepository) DeletePermanently(ic business.InternalContext, id int) error {
	query := repo.DB.WithContext(ic.ToContext())

	selectedUser := selectUser.User{}

	tx := query.Begin()
	var count int64
	// if tx.Unscoped().Select("id").First(&user).Error != nil {}
	if tx.Unscoped().Select("id").First(&selectedUser).Count(&count); count != 1 {
		tx.Rollback()
		return errors.New("User Not Found")
	}

	// if tx.Unscoped().Delete(&user, id).Error != nil {}
	if err := tx.Unscoped().Delete(&selectedUser, id).Error; err != nil {
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
