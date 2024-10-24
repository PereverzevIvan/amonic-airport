package mysql_repo

import (
	"gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/models"
	"gorm.io/gorm"
)

type QuestionAnswerRepo struct {
	Conn *gorm.DB
}

func NewQuestionAnswerRepo(conn *gorm.DB) QuestionAnswerRepo {
	return QuestionAnswerRepo{
		Conn: conn,
	}
}

func (r QuestionAnswerRepo) GetAllByQuestionID(question_id int) ([]models.QuestionAnswer, error) {
	var question_answers []models.QuestionAnswer
	err := r.Conn.
		Where("QuestionID = ?", question_id).
		Find(&question_answers).
		Error

	if err != nil {
		return nil, err
	}

	return question_answers, nil
}

func (r QuestionAnswerRepo) GetAll() ([]models.QuestionAnswer, error) {
	var question_answers []models.QuestionAnswer
	err := r.Conn.
		Find(&question_answers).
		Error

	if err != nil {
		return nil, err
	}

	return question_answers, nil
}
