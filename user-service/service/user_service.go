package service

import (
	"errors"
	"github.com/google/uuid"
	"user-service/repository"
	"user-service/schemas/dto"
	"user-service/schemas/web"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) SignUp(user web.NewUserIn) error {
	existingUser, _ := s.repo.GetByUsername(user.Username)
	if existingUser.Username != "" {
		return errors.New("user already exists")
	}

	return s.repo.Create(user)
}

func (s *UserService) Login(username, password string) error {
	user, err := s.repo.GetByUsername(username)
	if err != nil || user.Password != password {
		return errors.New("invalid username or password")
	}

	return nil
}

func (s *UserService) Update(user web.NewUserIn) error {
	return s.repo.Update(user)
}

func (s *UserService) GetCurrent(username string) (dto.UserDtoSchema, error) {
	return s.repo.GetByUsername(username)
}

func (s *UserService) GetByUsername(username string) (dto.UserDtoSchema, error) {
	return s.repo.GetByUsername(username)
}

func (s *UserService) GetRequiredById(id uuid.UUID) (dto.UserDtoSchema, error) {
	return s.repo.GetRequiredById(id)
}
