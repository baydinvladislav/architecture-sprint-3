package service

import (
	"errors"
	"github.com/google/uuid"
	"user-service/repository"
	"user-service/schemas/dto"
	"user-service/schemas/web"
)

type UserService struct {
	userRepository repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{
		userRepository: repo,
	}
}

func (s *UserService) SignUp(user web.NewUserIn) error {
	existingUser, _ := s.userRepository.GetByUsername(user.Username)
	if existingUser.Username != "" {
		return errors.New("user already exists")
	}

	return s.userRepository.Create(user)
}

func (s *UserService) Login(username, password string) error {
	user, err := s.userRepository.GetByUsername(username)
	if err != nil || user.Password != password {
		return errors.New("invalid username or password")
	}

	return nil
}

func (s *UserService) Update(user web.NewUserIn) error {
	return s.userRepository.Update(user)
}

func (s *UserService) GetCurrent(username string) (dto.UserDtoSchema, error) {
	return s.userRepository.GetByUsername(username)
}

func (s *UserService) GetByUsername(username string) (dto.UserDtoSchema, error) {
	return s.userRepository.GetByUsername(username)
}

func (s *UserService) GetRequiredById(id uuid.UUID) (dto.UserDtoSchema, error) {
	return s.userRepository.GetRequiredById(id)
}
