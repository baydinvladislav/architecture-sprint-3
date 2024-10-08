package service

import (
	"errors"
	"user-service/presentation/web-schemas"
	"user-service/repository"
	"user-service/repository/dto-schemas"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) SignUp(user web_schemas.NewUserIn) error {
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

func (s *UserService) Update(user web_schemas.NewUserIn) error {
	return s.repo.Update(user)
}

func (s *UserService) GetCurrent(username string) (dto_schemas.UserDtoSchema, error) {
	return s.repo.GetByUsername(username)
}

func (s *UserService) GetByUsername(username string) (dto_schemas.UserDtoSchema, error) {
	return s.repo.GetByUsername(username)
}

func (s *UserService) GetRequiredById(id uint) (dto_schemas.UserDtoSchema, error) {
	return s.repo.GetRequiredById(id)
}
