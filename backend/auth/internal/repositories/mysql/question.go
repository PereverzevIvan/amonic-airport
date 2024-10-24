package mysql_repo

import (
	"gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/models"
	"gorm.io/gorm"
)

type QuestionRepo struct {
	Conn *gorm.DB
}

func NewQuestionRepo(conn *gorm.DB) QuestionRepo {
	return QuestionRepo{
		Conn: conn,
	}
}

func (r QuestionRepo) GetByID(id int) (*models.Question, error) {
	var question models.Question
	err := r.Conn.
		First(&question, id).
		Error
	if err != nil {
		return nil, err
	}

	return &question, nil
}

func (r QuestionRepo) GetAll() ([]models.Question, error) {
	var questions []models.Question
	err := r.Conn.
		Find(&questions).
		Error
	if err != nil {
		return nil, err
	}

	return questions, nil
}
