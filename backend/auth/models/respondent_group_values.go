package models

type RespondentGroupValue struct {
	RespondentID int         `json:"respondent_id" gorm:"Column:RespondentID"`
	Respondent   *Respondent `json:"respondent" gorm:"foreignKey:RespondentID"`
	GroupValueID int         `json:"group_value_id" gorm:"Column:GroupValueID"`
	GroupValue   *GroupValue `json:"group_value" gorm:"foreignKey:GroupValueID"`
}

func (*RespondentGroupValue) TableName() string {
	return "respondent_group_values"
}
