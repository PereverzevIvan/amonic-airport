package service

import (
	"encoding/csv"
	"fmt"
	"io"
	"mime/multipart"
	"regexp"
	"strconv"
	"time"

	"gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/models"
	"github.com/gofiber/fiber/v2/log"
)

type SurveyRepo interface {
	GetAnswers(params *models.SurveyAnswersParams) ([]models.SurveyAnswerResult, error)
	CreateWithRespondentsAndAnswers(
		survey *models.Survey,
		survey_questions []models.SurveyQuestion,
		respondents []models.Respondent,
		respondents_group_values [][]models.RespondentGroupValue,
		respondents_answers [][]models.RespondentAnswer,
	) error
}

type GroupRepo interface {
	GetByName(name string) (*models.Group, error)
}

type GroupValueRepo interface {
	GetAllByGroupID(group_id int) ([]models.GroupValue, error)
}

type QuestionRepo interface {
	GetByID(id int) (*models.Question, error)
	GetAll() ([]models.Question, error)
}
type QuestionAnswerRepo interface {
	GetAllByQuestionID(question_id int) ([]models.QuestionAnswer, error)
	GetAll() ([]models.QuestionAnswer, error)
}

type surveyService struct {
	surveyRepo         SurveyRepo
	groupRepo          GroupRepo
	groupValueRepo     GroupValueRepo
	questionRepo       QuestionRepo
	questionAnswerRepo QuestionAnswerRepo
}

func NewSurveyService(
	surveyRepo SurveyRepo,
	groupRepo GroupRepo,
	groupValueRepo GroupValueRepo,
	questionRepo QuestionRepo,
	questionAnswerRepo QuestionAnswerRepo,
) surveyService {
	return surveyService{
		surveyRepo:         surveyRepo,
		groupRepo:          groupRepo,
		groupValueRepo:     groupValueRepo,
		questionRepo:       questionRepo,
		questionAnswerRepo: questionAnswerRepo,
	}
}

func (s surveyService) GetAnswers(params *models.SurveyAnswersParams) ([]models.SurveyAnswerResult, error) {
	return s.surveyRepo.GetAnswers(params)
}

/*
Алгоритм добавления нового опроса и статистики голосов для него:
1. Получить дату из названия файла
2. Получить заголовок csv файла (1-ая строка)
3. Соотнести колонки с функциями, которые будем вызывать
  - колонки групп клиентов
    (в ячейках - категория группы клиента)
  - колонка вопроса и ответы клиентов
    (в ячейках - ответ клиента на вопрос)

4. Парсим строки со значениями

5. Создаем опрос
6. Добавляем surveyd_clients
7. Добавляем surveyd_clients_answers

Нужно:

	count_groups int -- кол-во групп
	group_ids []int -- id групп в порядке следования в заголовке

	count_questions int -- кол-во вопросов
	question_ids []int -- id вопросов в порядке следования в заголовке

	groups_map  map[string][]string -- map[название_группы]категории
	categories_map map[string]SuveyedGroupCategory -- map[название_категории]категория

	surveyed_clients_map map[int]SurveydClients -- map[id_категории]объект_клиентов
	surveyed_clients_answers_map map[int]SurveydClientsAnswers -- map[id_категории]объект_ответов
*/
func (s surveyService) AddSurveyFromCSV(file *multipart.FileHeader) error {
	// Parse the date from the filename
	file_date, err := parseDateFromName(file.Filename)
	if err != nil {
		return err
	}

	// Create a new survey
	survey := &models.Survey{
		Date:             *file_date,
		RespondentsCount: 0,
	}

	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		return models.ErrFailedToOpenFile
	}
	defer src.Close()

	// Create a buffer to store the file contents
	reader := csv.NewReader(src)
	// Set FieldsPerRecord to -1 to allow variable number of fields
	reader.FieldsPerRecord = -1

	var header_groups []models.Group
	var header_groups_values [][]models.GroupValue

	var header_questions []models.Question
	var header_questions_answers [][]models.QuestionAnswer

	// Read the CSV headers and create a new survey
	err = s.readSurveyCSVHeader(
		reader,
		&header_groups, &header_groups_values,
		&header_questions, &header_questions_answers,
	)
	if err != nil {
		return err
	}

	// log.Info(header_groups)
	// log.Info(header_questions)

	respondents := []models.Respondent{}
	respondents_group_values := [][]models.RespondentGroupValue{}
	respondents_answers := [][]models.RespondentAnswer{}

	// Read remaining rows
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break // End of file reached
		}
		if err != nil { // Log the error and continue to the next record
			log.Info("Error reading record: %v\n", err)
			continue
		}

		// log.Info(record)
		// Check if the record has the correct number of fields
		if len(record) != len(header_groups)+len(header_questions) {
			continue
		}

		respondent := models.Respondent{}
		respondent_group_values := []models.RespondentGroupValue{}
		respondent_answers := []models.RespondentAnswer{}

		var i int = 0
		// parse group_values
		for ; i < len(header_groups); i++ {
			if record[i] == "" {
				continue
			}
			group := &header_groups[i]
			group_values := header_groups_values[i]

			switch group.Name {
			case "Gender":
				switch record[i] {
				case "M":
					respondent_group_values = append(respondent_group_values, models.RespondentGroupValue{GroupValueID: group_values[0].ID})
				case "F":
					respondent_group_values = append(respondent_group_values, models.RespondentGroupValue{GroupValueID: group_values[1].ID})
				}
			case "Age":
				respondent_age, err := strconv.Atoi(record[i])
				if err != nil || respondent_age < 18 {
					continue
				}

				switch {
				case respondent_age < 25:
					respondent_group_values = append(respondent_group_values, models.RespondentGroupValue{GroupValueID: group_values[0].ID})
				case respondent_age < 40:
					respondent_group_values = append(respondent_group_values, models.RespondentGroupValue{GroupValueID: group_values[1].ID})
				case respondent_age < 60:
					respondent_group_values = append(respondent_group_values, models.RespondentGroupValue{GroupValueID: group_values[2].ID})
				default:
					respondent_group_values = append(respondent_group_values, models.RespondentGroupValue{GroupValueID: group_values[3].ID})
				}
			default:
				for _, group_value := range group_values {
					if group_value.Name == record[i] {
						respondent_group_values = append(respondent_group_values, models.RespondentGroupValue{GroupValueID: group_value.ID})
						break
					}
				}
			}
		}

		// parse question answers
		for ; i < len(header_groups)+len(header_questions); i++ {
			if record[i] == "" {
				continue
			}

			answer_value, err := strconv.Atoi(record[i])
			if err != nil {
				continue
			}

			for _, answer := range header_questions_answers[i-len(header_groups)] {
				if answer.Value == answer_value {
					respondent_answers = append(respondent_answers,
						models.RespondentAnswer{
							QuestionAnswerID: answer.ID,
						})
					break
				}
			}
		}

		// Добавляем данные
		respondents = append(respondents, respondent)
		respondents_group_values = append(respondents_group_values, respondent_group_values)
		respondents_answers = append(respondents_answers, respondent_answers)

		survey.RespondentsCount++
	}

	if survey.RespondentsCount == 0 {
		return fmt.Errorf("empty survey, no respondents")
	}

	survey_questions := []models.SurveyQuestion{}
	for i := range header_questions {
		survey_question := models.SurveyQuestion{
			QuestionID: header_questions[i].ID,
		}
		survey_questions = append(survey_questions, survey_question)
	}

	err = s.surveyRepo.CreateWithRespondentsAndAnswers(
		survey,
		survey_questions,
		respondents,
		respondents_group_values,
		respondents_answers,
	)
	if err != nil {
		return err
	}

	return nil
}

func parseDateFromName(filename string) (*time.Time, error) {
	// parse date from filename
	var file_date time.Time

	re := regexp.MustCompile(`(\d{4}-\d{2}-\d{2})`)
	matches := re.FindStringSubmatch(filename)

	if len(matches) == 0 {
		return nil, models.ErrFailedToParseDateFromName
	}
	dateStr := matches[1]
	file_date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, err
	}
	return &file_date, nil
}

func (s surveyService) readSurveyCSVHeader(
	reader *csv.Reader,
	header_groups *[]models.Group,
	header_groups_values *[][]models.GroupValue,

	header_questions *[]models.Question,
	header_questions_answers *[][]models.QuestionAnswer,
) error {
	// Read the CSV headers and create a new survey
	headers, err := reader.Read()
	if err == io.EOF {
		return fmt.Errorf("invalid headers, found EOF")
	}
	if err != nil { // Log the error and continue to the next record
		log.Info("Error reading record: %v\n", err)
		return err
	}

	for _, header := range headers {

		log.Info(header)

		re := regexp.MustCompile(`^Q(\d+)$`)
		matches := re.FindStringSubmatch(header)
		if len(matches) == 0 {
			if len(*header_questions) > 0 {
				return fmt.Errorf("group header %v in header after questions started", header)
			}

			group, err := s.groupRepo.GetByName(header)
			if err != nil {
				return err
			}
			if group == nil {
				return fmt.Errorf("group %v not found", header)
			}

			group_values, err := s.groupValueRepo.GetAllByGroupID(group.ID)
			if err != nil {
				return err
			}

			*header_groups = append(*header_groups, *group)
			*header_groups_values = append(*header_groups_values, group_values)

			continue
		}

		question_id, err := strconv.Atoi(matches[1])
		if err != nil {
			return err
		}

		question, err := s.questionRepo.GetByID(question_id)
		if err != nil {
			return err
		}
		if question == nil {
			return fmt.Errorf("question %v not found", header)
		}

		question_answers, err := s.questionAnswerRepo.GetAllByQuestionID(question_id)
		if err != nil {
			return err
		}

		*header_questions = append(*header_questions, *question)
		*header_questions_answers = append(*header_questions_answers, question_answers)
	}

	if len(*header_groups) == 0 {
		return fmt.Errorf("no groups in file")
	}
	if len(*header_questions) == 0 {
		return fmt.Errorf("no questions in file")
	}

	return nil
}
