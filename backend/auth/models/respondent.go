package models

type Respondent struct {
	ID       int     `json:"id" gorm:"Column:ID"`
	SurveyID int     `json:"survey_id" gorm:"Column:SurveyID"`
	Survey   *Survey `json:"survey" gorm:"foreignKey:SurveyID"`
}

func (*Respondent) TableName() string {
	return "respondents"
}
