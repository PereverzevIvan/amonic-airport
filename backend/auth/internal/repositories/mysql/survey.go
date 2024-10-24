package mysql_repo

import (
	"fmt"

	"gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/models"
	"gorm.io/gorm"
)

type SurveyRepo struct {
	Conn *gorm.DB
}

func NewSurveyRepo(conn *gorm.DB) SurveyRepo {
	return SurveyRepo{
		Conn: conn,
	}
}

func (r SurveyRepo) CreateWithRespondentsAndAnswers(
	survey *models.Survey,
	survey_questions []models.SurveyQuestion,
	respondents []models.Respondent,
	respondents_group_values [][]models.RespondentGroupValue,
	respondents_answers [][]models.RespondentAnswer,
) error {

	tx := r.Conn.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	if err := tx.Create(survey).Error; err != nil {
		tx.Rollback()
		return err
	}

	for i := range survey_questions {
		survey_questions[i].SurveyID = survey.ID
	}
	if err := tx.Create(&survey_questions).Error; err != nil {
		tx.Rollback()
		return err
	}

	for i := range respondents {
		respondents[i].SurveyID = survey.ID
	}
	if err := tx.Create(&respondents).Error; err != nil {
		tx.Rollback()
		return err
	}

	for i, respondent := range respondents {
		// log.Info("bbb", len(respondents_group_values[i]))
		for j := range respondents_group_values[i] {
			respondents_group_values[i][j].RespondentID = respondent.ID
		}
		for j := range respondents_answers[i] {
			respondents_answers[i][j].RespondentID = respondent.ID
		}

		if err := tx.Create(&respondents_group_values[i]).Error; err != nil {
			tx.Rollback()
			return err
		}

		if err := tx.Create(&respondents_answers[i]).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func (r SurveyRepo) GetAnswers(params *models.SurveyAnswersParams) ([]models.SurveyAnswerResult, error) {
	var answers []models.SurveyAnswerResult

	rows, err := r.Conn.
		Raw(`
			SELECT
				question_answers.QuestionID AS question_id,
				question_answers.ID AS answer_id,
				group_values.GroupID AS group_id,
				rgv.GroupValueID AS group_value_id,
				COUNT(*) as count_answers
			FROM
				surveys
				
			INNER JOIN respondents ON respondents.SurveyID = surveys.ID
			`+generateFiltersForGetAnswersSQL(params)+`
			-- INNER JOIN respondent_group_values as rgv_filter1 ON rgv_filter1.RespondentID = respondents.ID AND rgv_filter1.GroupValueID = 39 
			-- INNER JOIN respondent_group_values as rgv_filter2 ON rgv_filter2.RespondentID = respondents.ID AND rgv_filter2.GroupValueID = 42 

			INNER JOIN respondent_group_values AS rgv ON rgv.RespondentID = respondents.ID
			INNER JOIN group_values ON group_values.ID = rgv.GroupValueID

			INNER JOIN respondent_answers ON respondent_answers.RespondentID = respondents.ID
			INNER JOIN question_answers ON question_answers.ID = respondent_answers.QuestionAnswerID

			WHERE
				surveys.Date >= ? 
				AND surveys.Date <= ? 

			GROUP BY 
				question_id, 
				answer_id,
				group_id,
				group_value_id`,
			// ...params.FilterGroupValues,

			params.BeginDate,
			params.EndDate,
		).
		Rows()

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var answer models.SurveyAnswerResult
		err = rows.Scan(&answer.QuestionID, &answer.AnswerID, &answer.GroupID, &answer.GroupValueID, &answer.CountAnswers)
		if err != nil {
			return nil, err
		}
		answers = append(answers, answer)
	}

	return answers, nil
}

func generateFiltersForGetAnswersSQL(params *models.SurveyAnswersParams) string {
	var filters string

	for i, group_value_id := range params.FilterGroupValues {
		inner_join_group_value_filter := fmt.Sprintf(
			`INNER JOIN respondent_group_values AS rgv_filter%v ON rgv_filter%v.RespondentID = respondents.ID AND rgv_filter%v.GroupValueID = %v`,
			i+1, i+1, i+1, group_value_id)

		filters += inner_join_group_value_filter + "\n			"
		// append(rgv_filters, "(rgv_filter"+strconv.Itoa(i)+".GroupValueID = "+strconv.Itoa(v)+")")
		// filters = append(filters, "(rgv_filter"+strconv.Itoa(i)+".GroupValueID = "+strconv.Itoa(v)+")")
	}
	return filters
}
