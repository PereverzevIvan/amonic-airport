package service

import (
	"gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/models"
	"golang.org/x/crypto/bcrypt"
)

type UserRepo interface {
	GetByID(user_id int) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetAll(params map[string]string) (*[]models.User, error)
	IsAdmin(user_id int) (bool, error)
	Create(*models.User) error
	IsActive(user_id int) (bool, error)
	Update(user *models.User) error
	UpdateActive(user_id int, is_active bool) error
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
	hash_byte := []byte(user.Password)

	err := bcrypt.CompareHashAndPassword(hash_byte, []byte(password))
	return err == nil
}

func (us UserService) GetByID(user_id int) (*models.User, error) {
	user, err := us.userRepo.GetByID(user_id)
	return user, err
}

func (us UserService) GetByEmail(email string) (*models.User, error) {
	user, err := us.userRepo.GetByEmail(email)
	return user, err
}

func (us UserService) GetAll(params map[string]string) (*[]models.User, error) {
	users, err := us.userRepo.GetAll(params)
	return users, err
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

func HashPassword(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return []byte{}, err
	}

	return hashedPassword, nil
}

func (us UserService) Create(user *models.User) error {
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	return us.userRepo.Create(user)
}

func (us UserService) Update(user *models.User) error {
	return us.userRepo.Update(user)
}

func (us UserService) UpdateActive(user_id int, is_active bool) error {
	return us.userRepo.UpdateActive(user_id, is_active)
}
