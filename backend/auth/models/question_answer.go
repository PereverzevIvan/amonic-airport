package models

type QuestionAnswer struct {
	ID         int       `json:"id" gorm:"Column:ID"`
	QuestionID int       `json:"question_id" gorm:"Column:QuestionID"`
	Question   *Question `json:"question" gorm:"foreignKey:QuestionID"`
	Value      int       `json:"value" gorm:"Column:Value"`
	Text       string    `json:"text" gorm:"Column:Text"`
}

func (*QuestionAnswer) TableName() string {
	return "question_answers"
}

type QuestionWithAnswers struct {
	Question
	Answers []QuestionAnswer `json:"answers"`
}
