package service

import (
	"gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/models"
	"golang.org/x/crypto/bcrypt"
	// "github.com/gofiber/fiber/v3/log"
)

type UserRepo interface {
	GetByID(user_id int) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	IsAdmin(user_id int) (bool, error)
	// GetByEmail(user_email string) (models.User, error)
	Create(*models.User) error
	IsActive(user_id int) (bool, error)
	// Create(*models.User) error
	// Update(*models.User) error
	// Delete(user_id int) error
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

func (us UserService) IsPasswordCorrect(user *models.User, password string) bool {
	return user.Password == password
}

func (us UserService) GetByID(user_id int) (*models.User, error) {
	user, err := us.userRepo.GetByID(user_id)
	return user, err
}

func (us UserService) GetByEmail(email string) (*models.User, error) {
	user, err := us.userRepo.GetByEmail(email)
	return user, err
}

func (us UserService) IsActive(user_id int) (bool, error) {
	is_active, err := us.userRepo.IsActive(user_id)
	if err != nil {
		return false, err
	}
	return is_active, err
}

func (us UserService) IsAdmin(user_id int) (bool, error) {
	is_admin, err := us.userRepo.IsAdmin(user_id)
	if err != nil {
		return false, err
	}
	return is_admin, err
}

func (us UserService) Create(user *models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	return us.userRepo.Create(user)
}
