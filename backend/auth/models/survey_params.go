package models

import (
	"fmt"
	"time"
)

type SurveyAnswersParams struct {
	BeginDate         string `query:"begin_date"`
	EndDate           string `query:"end_date"`
	FilterGroupValues []int  `query:"filter_group_values"`
}

func (s *SurveyAnswersParams) Validate() error {
	_, err := time.Parse("2006-01-02", s.BeginDate)
	if err != nil {
		return fmt.Errorf("invalid begin_date: %v", err.Error())
	}

	_, err = time.Parse("2006-01-02", s.EndDate)
	if err != nil {
		return fmt.Errorf("invalid end_date: %v", err.Error())
	}
	return nil
}

type SurveyAnswerResult struct {
	QuestionID   int `json:"question_id"`
	AnswerID     int `json:"answer_id"`
	GroupID      int `json:"group_id"`
	GroupValueID int `json:"group_value_id"`
	CountAnswers int `json:"count_answers"`
}
