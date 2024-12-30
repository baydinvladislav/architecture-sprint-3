package repository

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"user-service/persistance"
	"user-service/schemas/dto"
)

type UserRepository interface {
	Create(user *dto.UserDtoSchema) (*dto.UserDtoSchema, error)
	GetByUsername(username string) (*dto.UserDtoSchema, error)
	Update(user *dto.UserDtoSchema) (*dto.UserDtoSchema, error)
	GetRequiredById(id uuid.UUID) (*dto.UserDtoSchema, error)
}

type GORMUserRepository struct {
	db *gorm.DB
}

func NewGORMUserRepository(db *gorm.DB) *GORMUserRepository {
	return &GORMUserRepository{
		db: db,
	}
}

func (r *GORMUserRepository) Create(user *dto.UserDtoSchema) (*dto.UserDtoSchema, error) {
	var existingUser persistance.UserModel
	if err := r.db.Where("username = ?", user.Username).First(&existingUser).Error; err == nil {
		return &dto.UserDtoSchema{}, errors.New("user already exists")
	}

	newUser := persistance.UserModel{
		Username: user.Username,
		Password: user.Password,
	}

	err := r.db.Create(&newUser).Error
	if err != nil {
		return &dto.UserDtoSchema{}, errors.New("error during database insertion")
	}

	createdUser, err := r.GetByUsername(user.Username)
	if err != nil {
		return &dto.UserDtoSchema{}, err
	}

	return &dto.UserDtoSchema{
		ID:       createdUser.ID,
		Username: createdUser.Username,
		Password: createdUser.Password,
		Email:    createdUser.Email,
	}, nil
}

func (r *GORMUserRepository) GetByUsername(username string) (*dto.UserDtoSchema, error) {
	var user persistance.UserModel
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return &dto.UserDtoSchema{}, errors.New("user not found")
	}

	return &dto.UserDtoSchema{
		ID:       user.ID,
		Username: user.Username,
		Password: user.Password,
	}, nil
}

func (r *GORMUserRepository) GetRequiredById(id uuid.UUID) (*dto.UserDtoSchema, error) {
	var user persistance.UserModel
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return &dto.UserDtoSchema{}, errors.New("user not found")
	}

	return &dto.UserDtoSchema{
		ID:       user.ID,
		Username: user.Username,
		Password: user.Password,
	}, nil
}

func (r *GORMUserRepository) Update(user *dto.UserDtoSchema) (*dto.UserDtoSchema, error) {
	var existingUser dto.UserDtoSchema
	if err := r.db.Where("username = ?", user.Username).First(&existingUser).Error; err != nil {
		return &dto.UserDtoSchema{}, nil
	}

	existingUser.Password = user.Password
	err := r.db.Save(&existingUser).Error
	if err != nil {
		return &dto.UserDtoSchema{}, errors.New("error during updating user data")
	}

	return r.GetByUsername(existingUser.Username)
}
