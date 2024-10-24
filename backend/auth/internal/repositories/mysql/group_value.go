package mysql_repo

import (
	"gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/models"
	"gorm.io/gorm"
)

type GroupValueRepo struct {
	Conn *gorm.DB
}

func NewGroupValueRepo(conn *gorm.DB) GroupValueRepo {
	return GroupValueRepo{
		Conn: conn,
	}
}

func (r GroupValueRepo) GetAllByGroupID(group_id int) ([]models.GroupValue, error) {
	var group_values []models.GroupValue

	err := r.Conn.
		Where("GroupID = ?", group_id).
		Find(&group_values).
		Error

	if err != nil {
		return nil, err
	}

	return group_values, nil
}

func (r GroupValueRepo) GetAll() ([]models.GroupValue, error) {
	var group_values []models.GroupValue

	err := r.Conn.
		Find(&group_values).
		Error

	if err != nil {
		return nil, err
	}

	return group_values, nil
}
