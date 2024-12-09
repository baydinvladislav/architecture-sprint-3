package service

import (
	"errors"
	"github.com/google/uuid"
	"user-service/repository"
	"user-service/schemas/dto"
	"user-service/schemas/web"
)

type UserService struct {
	AuthService    *AuthService
	userRepository repository.UserRepository
}

func NewUserService(authService *AuthService, repository repository.UserRepository) *UserService {
	return &UserService{
		AuthService:    authService,
		userRepository: repository,
	}
}

func (s *UserService) SignUp(userData *web.UserIn) (*web.LoginResponse, error) {
	existingUser, _ := s.userRepository.GetByUsername(userData.Username)
	if existingUser.Username != "" {
		return &web.LoginResponse{}, errors.New("username already exists")
	}

	userDto := &dto.UserDtoSchema{
		Username: userData.Username,
		Password: userData.Password,
	}

	newUser, err := s.userRepository.Create(userDto)
	if err != nil {
		return &web.LoginResponse{}, errors.New("error during userData creation")
	}

	accessToken, err := s.AuthService.GenerateAccessToken(newUser.Username)
	if err != nil {
		return &web.LoginResponse{}, errors.New("error during access token generation")
	}

	refreshToken, err := s.AuthService.GenerateRefreshToken(newUser.Username)
	if err != nil {
		return &web.LoginResponse{}, errors.New("error during refresh token generation")
	}

	return &web.LoginResponse{
		ID:           newUser.ID,
		Username:     newUser.Username,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *UserService) Login(username, password string) (*web.LoginResponse, error) {
	user, err := s.userRepository.GetByUsername(username)
	// TODO: а пароль что не захэширован что ли?
	if err != nil || user.Password != password {
		return &web.LoginResponse{}, errors.New("invalid username or password")
	}

	accessToken, err := s.AuthService.GenerateAccessToken(user.Username)
	if err != nil {
		return &web.LoginResponse{}, errors.New("error during access token generation")
	}

	refreshToken, err := s.AuthService.GenerateRefreshToken(user.Username)
	if err != nil {
		return &web.LoginResponse{}, errors.New("error during refresh token generation")
	}

	return &web.LoginResponse{
		ID:           user.ID,
		Username:     user.Username,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *UserService) Update(user web.UserIn) (*web.UserOut, error) {
	updatedUser, err := s.userRepository.Update(&dto.UserDtoSchema{
		Username: user.Username,
		Password: user.Password,
	})
	if err != nil {
		return &web.UserOut{}, nil
	}

	return &web.UserOut{
		ID:       updatedUser.ID,
		Username: updatedUser.Username,
	}, nil
}

func (s *UserService) GetCurrent(username string) (*web.UserOut, error) {
	currentUser, err := s.userRepository.GetByUsername(username)
	if err != nil {
		return &web.UserOut{}, err
	}

	return &web.UserOut{
		ID:       currentUser.ID,
		Username: currentUser.Username,
	}, nil
}

func (s *UserService) GetByUsername(username string) (*web.UserOut, error) {
	user, err := s.userRepository.GetByUsername(username)
	if err != nil {
		return &web.UserOut{}, err
	}

	return &web.UserOut{
		ID:       user.ID,
		Username: user.Password,
	}, nil
}

func (s *UserService) GetRequiredById(id uuid.UUID) (*web.UserOut, error) {
	user, err := s.userRepository.GetRequiredById(id)
	if err != nil {
		return &web.UserOut{}, err
	}

	return &web.UserOut{
		ID:       user.ID,
		Username: user.Password,
	}, nil
}

func (s *UserService) UpdateUserRefreshToken(refreshToken string) (string, string, error) {
	claims, err := s.AuthService.ValidateRefreshToken(refreshToken)
	if err != nil {
		return "", "", err
	}

	accessToken, err := s.AuthService.GenerateAccessToken(claims.Username)
	if err != nil {
		return "", "", err
	}

	newRefreshToken, err := s.AuthService.GenerateRefreshToken(claims.Username)
	if err != nil {
		return "", "", err
	}

	return accessToken, newRefreshToken, nil
}
