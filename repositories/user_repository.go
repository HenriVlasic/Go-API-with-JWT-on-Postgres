package repository

import (
	"cyberzell.com/seguros/models"
	"cyberzell.com/seguros/models/apperrors"
	"errors"
	"log"
	"regexp"
	"strconv"

	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) models.IUserRepository {
	return &UserRepository{
		DB:db,
	}
}

func (r *UserRepository) FindUserByEmail(email string) (*models.User, error) {
	user := &models.User{}

	if err := r.DB.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, apperrors.NewNotFound("email", email)
		}
		return user, apperrors.NewInternal()
	}

	return user, nil

}

func (r *UserRepository) GetUserById(id int) (*models.User, error) {
	user := &models.User{}
	userId := strconv.Itoa(id)

	if err := r.DB.Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, apperrors.NewNotFound("id", userId)
		}
		return user, apperrors.NewInternal()
	}

	return user, nil
}

func (r *UserRepository) CreateUser(user *models.User) (*models.User, error) {
	if result := r.DB.Create(&user); result.Error != nil {
		if isDuplicateKeyError(result.Error) {
			return nil, apperrors.NewBadRequest(apperrors.DuplicateEmail)
		}

		log.Printf("Could not create a user with email: %v. Reason: %v\n", user.Email, result.Error)
		return nil, apperrors.NewInternal()
	}

	return user, nil
}

func isDuplicateKeyError(err error) bool {
	duplicate := regexp.MustCompile(`\(SQLSTATE 23505\)$`)
	return duplicate.MatchString(err.Error())
}