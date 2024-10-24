package mysql_repo

import (
	"gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/models"
	"gorm.io/gorm"
)

type GroupRepo struct {
	Conn *gorm.DB
}

func NewGroupRepo(conn *gorm.DB) GroupRepo {
	return GroupRepo{
		Conn: conn,
	}
}

func (r GroupRepo) GetByName(name string) (*models.Group, error) {
	var group *models.Group

	err := r.Conn.
		Where(&models.Group{Name: name}).
		First(&group).
		Error
	if err != nil {
		return nil, err
	}

	return group, nil
}

func (r GroupRepo) GetAll() ([]models.Group, error) {
	var groups []models.Group

	err := r.Conn.
		Find(&groups).
		Error
	if err != nil {
		return nil, err
	}

	return groups, nil
}
