package service

import (
	"gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/models"
	// "github.com/gofiber/fiber/v3/log"
)

type UserRepo interface {
	GetByID(user_id int) (models.User, error)
	// GetByEmail(user_email string) (models.User, error)
	// Create(*models.User) error
	// Update(*models.User) error
	// Delete(user_id int) error
}

type OfficeRepo interface {
	GetByID(office_id int) (models.Office, error)
}

type RoleRepo interface {
	GetByID(role_id int) (models.Role, error)
}

type UserService struct {
	userRepo UserRepo
}

func NewUserService(ur UserRepo) UserService {
	return UserService{ur}
}

func (us UserService) GetByID(user_id int) (*models.User, error) {
	user, err := us.userRepo.GetByID(user_id)
	if err != nil {
		return &models.User{}, err
	}

	return &user, nil
}
