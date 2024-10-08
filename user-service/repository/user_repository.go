package repository

import (
	"errors"
	"user-service/persistance"
	"user-service/presentation/web-schemas"
	"user-service/repository/dto-schemas"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user web_schemas.NewUserIn) error
	GetByUsername(username string) (dto_schemas.UserDtoSchema, error)
	Update(user web_schemas.NewUserIn) error
	GetRequiredById(id uint) (dto_schemas.UserDtoSchema, error)
}

type GORMUserRepository struct {
	db *gorm.DB
}

func NewGORMUserRepository(db *gorm.DB) *GORMUserRepository {
	return &GORMUserRepository{
		db: db,
	}
}

func (r *GORMUserRepository) Create(user web_schemas.NewUserIn) error {
	var existingUser persistance.UserModel
	if err := r.db.Where("username = ?", user.Username).First(&existingUser).Error; err == nil {
		return errors.New("user already exists")
	}

	newUser := persistance.UserModel{
		Username: user.Username,
		Password: user.Password,
	}

	return r.db.Create(&newUser).Error
}

func (r *GORMUserRepository) GetByUsername(username string) (dto_schemas.UserDtoSchema, error) {
	var user persistance.UserModel
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return dto_schemas.UserDtoSchema{}, errors.New("user not found")
	}

	userDto := dto_schemas.UserDtoSchema{
		ID:       user.ID,
		Username: user.Username,
		Password: user.Password,
	}

	return userDto, nil
}

func (r *GORMUserRepository) GetRequiredById(id uint) (dto_schemas.UserDtoSchema, error) {
	var user persistance.UserModel
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return dto_schemas.UserDtoSchema{}, errors.New("user not found")
	}

	userDto := dto_schemas.UserDtoSchema{
		ID:       user.ID,
		Username: user.Username,
		Password: user.Password,
	}

	return userDto, nil
}

func (r *GORMUserRepository) Update(user web_schemas.NewUserIn) error {
	var existingUser dto_schemas.UserDtoSchema
	if err := r.db.Where("username = ?", user.Username).First(&existingUser).Error; err != nil {
		return errors.New("user not found")
	}
	existingUser.Password = user.Password
	return r.db.Save(&existingUser).Error
}
