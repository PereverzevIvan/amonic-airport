package usecase

import (
	"gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/models"
)

type user struct {
}

type userService interface {
	GetByID(user_id int) (*models.User, error)
	IsAdmin(user_id int) (bool, error)
}

type userUseCase struct {
	userService userService
}

func NewUserUseCase(userService userService) *userUseCase {
	return &userUseCase{userService: userService}
}

func (uu *userUseCase) GetByID(user_id int) (*models.User, error) {
	user, err := uu.userService.GetByID(user_id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (uu *userUseCase) IsAdmin(user_id int) (bool, error) {
	is_admin, err := uu.userService.IsAdmin(user_id)
	return is_admin, err
}
