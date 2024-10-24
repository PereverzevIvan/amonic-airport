package models

type RespondentAnswer struct {
	RespondentID     int             `json:"respondent_id" gorm:"Column:RespondentID"`
	Respondent       *Respondent     `json:"respondent" gorm:"foreignKey:RespondentID"`
	QuestionAnswerID int             `json:"question_answer_id" gorm:"Column:QuestionAnswerID"`
	QuestionAnswer   *QuestionAnswer `json:"question_answer" gorm:"foreignKey:QuestionAnswerID"`
}

func (*RespondentAnswer) TableName() string {
	return "respondent_answers"
}
