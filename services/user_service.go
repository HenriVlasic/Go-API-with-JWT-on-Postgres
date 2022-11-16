package services

import (
	"cyberzell.com/seguros/models"
	"cyberzell.com/seguros/models/apperrors"
	"log"
)

type userService struct {
	UserRepository models.IUserRepository
}

type USConfig struct {
	UserRepository models.IUserRepository
}

func NewUserService(c *USConfig) models.IUserService {
	return &userService{
		UserRepository: c.UserRepository,
	}
}

const SecretKey = "secret"

func (s *userService) Login(email string, password string) (*models.User, error) {
	user, err := s.UserRepository.FindUserByEmail(email)

	if err != nil {
		return nil, apperrors.NewAuthorization(apperrors.InvalidCredentials)
	}

	match, err := comparePasswords(user.Password, password)

	if err != nil {
		return nil, apperrors.NewInternal()
	}

	if !match {
		return nil, apperrors.NewAuthorization(apperrors.InvalidCredentials)
	}

	return user, nil
}

func (s *userService) Register(user *models.User) (*models.User, error) {
	hashedPassword, err := hashPassword(user.Password)

	if err != nil {
		log.Printf("Unable to signup user for email: %v\n", user.Email)
		return nil, apperrors.NewInternal()
	}

	user.Password = hashedPassword

	return s.UserRepository.CreateUser(user)
}

func (s *userService) GetUserById(id int) (*models.User, error) {
	return s.UserRepository.GetUserById(id)
}