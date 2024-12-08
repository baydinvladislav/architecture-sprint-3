package repository

import (
	"errors"
	"github.com/google/uuid"
	"user-service/persistance"
	"user-service/schemas/dto"
	"user-service/schemas/web"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user web.NewUserIn) error
	GetByUsername(username string) (dto.UserDtoSchema, error)
	Update(user web.NewUserIn) error
	GetRequiredById(id uuid.UUID) (dto.UserDtoSchema, error)
}

type GORMUserRepository struct {
	db *gorm.DB
}

func NewGORMUserRepository(db *gorm.DB) *GORMUserRepository {
	return &GORMUserRepository{
		db: db,
	}
}

func (r *GORMUserRepository) Create(user web.NewUserIn) error {
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

func (r *GORMUserRepository) GetByUsername(username string) (dto.UserDtoSchema, error) {
	var user persistance.UserModel
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return dto.UserDtoSchema{}, errors.New("user not found")
	}

	userDto := dto.UserDtoSchema{
		ID:       user.ID,
		Username: user.Username,
		Password: user.Password,
	}

	return userDto, nil
}

func (r *GORMUserRepository) GetRequiredById(id uuid.UUID) (dto.UserDtoSchema, error) {
	var user persistance.UserModel
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return dto.UserDtoSchema{}, errors.New("user not found")
	}

	userDto := dto.UserDtoSchema{
		ID:       user.ID,
		Username: user.Username,
		Password: user.Password,
	}

	return userDto, nil
}

func (r *GORMUserRepository) Update(user web.NewUserIn) error {
	var existingUser dto.UserDtoSchema
	if err := r.db.Where("username = ?", user.Username).First(&existingUser).Error; err != nil {
		return errors.New("user not found")
	}
	existingUser.Password = user.Password
	return r.db.Save(&existingUser).Error
}
