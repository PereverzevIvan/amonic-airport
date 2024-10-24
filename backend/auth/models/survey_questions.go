package models

type SurveyQuestion struct {
	SurveyID   int `json:"survey_id" gorm:"Column:SurveyID"`
	QuestionID int `json:"question_id" gorm:"Column:QuestionID"`
}

func (*SurveyQuestion) TableName() string {
	return "survey_questions"
}
