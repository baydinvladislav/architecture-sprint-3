package service

import (
	"errors"
	"github.com/google/uuid"
	"user-service/repository"
	"user-service/schemas/dto"
	"user-service/schemas/web"
)

type UserService struct {
	repository repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{
		repository: repo,
	}
}

func (s *UserService) SignUp(user web.NewUserIn) error {
	existingUser, _ := s.repository.GetByUsername(user.Username)
	if existingUser.Username != "" {
		return errors.New("user already exists")
	}

	return s.repository.Create(user)
}

func (s *UserService) Login(username, password string) error {
	user, err := s.repository.GetByUsername(username)
	if err != nil || user.Password != password {
		return errors.New("invalid username or password")
	}

	return nil
}

func (s *UserService) Update(user web.NewUserIn) error {
	return s.repository.Update(user)
}

func (s *UserService) GetCurrent(username string) (dto.UserDtoSchema, error) {
	return s.repository.GetByUsername(username)
}

func (s *UserService) GetByUsername(username string) (dto.UserDtoSchema, error) {
	return s.repository.GetByUsername(username)
}

func (s *UserService) GetRequiredById(id uuid.UUID) (dto.UserDtoSchema, error) {
	return s.repository.GetRequiredById(id)
}
