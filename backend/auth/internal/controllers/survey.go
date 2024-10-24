package controllers

import (
	"net/http"
	"path/filepath"
	"strings"

	"gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/models"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
)

type SurveyController struct {
	jwtUseCase            jwtUseCase
	surveyService         surveyService
	groupService          GroupService
	groupValueService     GroupValueService
	questionService       QuestionService
	questionAnswerService QuestionAnswerService
}

func AddSurveyControllerRoutes(
	api *fiber.Router,
	jwtUseCase jwtUseCase,
	surveyService surveyService,
	groupService GroupService,
	groupValueService GroupValueService,
	questionService QuestionService,
	questionAnswerService QuestionAnswerService,
	authMiddleware AuthMiddleware,
) {
	controller := &SurveyController{
		jwtUseCase:            jwtUseCase,
		surveyService:         surveyService,
		groupService:          groupService,
		groupValueService:     groupValueService,
		questionService:       questionService,
		questionAnswerService: questionAnswerService,
	}

	// TODO: Add middleware
	(*api).Get("/surveys/respondents-answers", controller.RespondentsAnswers)
	(*api).Post("/surveys/upload", controller.SurveysUpload)
	// (*api).Get("/surveys", controller.Surveys, authMiddleware.IsActive)
	(*api).Get("/survey/groups-with-values", controller.GroupsWithValues)
	(*api).Get("/survey/questions-with-answers", controller.QuestionsWithAnswers)
	// (*api).Get("/survey/questions", controller.Surveys, authMiddleware.IsActive)
}

// @Summary      Загрузить файл CSV для обновления или добавления формы и результата опроса
// @Tags         Surveys
// @Accept       json
// @Produce      json
// @Success      200  {object} models.SurveysUploadResult
// @Failure      400
// @Failure      404
// @Router       /surveys/upload [post]
func (controller *SurveyController) SurveysUpload(ctx fiber.Ctx) error {
	// Get the file from the form input
	file, err := ctx.FormFile("file")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString("Failed to get file")
	}

	if file.Size > 1024*1024*15 {
		return ctx.Status(fiber.StatusBadRequest).SendString("File size must be less than 15 MB")
	}

	if ext := strings.ToLower(filepath.Ext(file.Filename)); ext != ".csv" {
		return ctx.Status(fiber.StatusBadRequest).SendString("Only CSV files are allowed")
	}

	// Add new survey from csv
	err = controller.surveyService.AddSurveyFromCSV(file)
	if err != nil {
		if err == models.ErrFailedToOpenFile ||
			err == models.ErrFailedToParseDateFromName {
			return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
		}
		log.Error(err)
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return ctx.SendStatus(http.StatusOK)
}

// Get All Groups (Список групп)
// @Summary      Get All Groups
// @Description  Get All Groups
// @Tags         Surveys
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[int]models.GroupWithValues
// @Failure      400
// @Failure      404
// @Router       /survey/groups-with-values [get]
func (controller *SurveyController) GroupsWithValues(ctx fiber.Ctx) error {

	groups, err := controller.groupService.GetAll()
	if err != nil {
		log.Error(err)
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	group_values, err := controller.groupValueService.GetAll()
	if err != nil {
		log.Error(err)
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	groupsWithValues := map[int]models.GroupWithValues{}
	for _, group := range groups {
		groupsWithValues[group.ID] = models.GroupWithValues{
			Group:  group,
			Values: []models.GroupValue{},
		}
	}
	for _, group_value := range group_values {
		currentGroupWithValues := groupsWithValues[group_value.GroupID]
		currentGroupWithValues.Values = append(currentGroupWithValues.Values, group_value)

		groupsWithValues[group_value.GroupID] = currentGroupWithValues
	}

	return ctx.Status(http.StatusOK).JSON(groupsWithValues)
}

func (controller *SurveyController) QuestionsWithAnswers(ctx fiber.Ctx) error {
	questions, err := controller.questionService.GetAll()
	if err != nil {
		log.Error(err)
		return ctx.SendStatus(http.StatusInternalServerError)
	}
	question_answers, err := controller.questionAnswerService.GetAll()
	if err != nil {
		log.Error(err)
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	questionWithAnswers := map[int]models.QuestionWithAnswers{}
	for _, question := range questions {
		questionWithAnswers[question.ID] = models.QuestionWithAnswers{
			Question: question,
			Answers:  []models.QuestionAnswer{},
		}
	}
	for _, question_answer := range question_answers {
		currentQuestionWithAnswers := questionWithAnswers[question_answer.QuestionID]
		currentQuestionWithAnswers.Answers = append(currentQuestionWithAnswers.Answers, question_answer)

		questionWithAnswers[question_answer.QuestionID] = currentQuestionWithAnswers
	}
	return ctx.Status(http.StatusOK).JSON(questionWithAnswers)
}

func (controller *SurveyController) RespondentsAnswers(ctx fiber.Ctx) error {
	var params models.SurveyAnswersParams
	if err := ctx.Bind().Query(&params); err != nil {
		log.Error(err)
		ctx.SendStatus(http.StatusBadRequest)
		return ctx.SendString(err.Error())
	}

	// Проверка параметров
	err := params.Validate()
	if err != nil {
		ctx.SendStatus(http.StatusBadRequest)
		return ctx.SendString(err.Error())
	}

	// Получаем список полетов
	survey_answers, err := controller.surveyService.GetAnswers(&params)
	if err != nil {
		log.Error(err)
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	return ctx.Status(http.StatusOK).JSON(survey_answers)
}
