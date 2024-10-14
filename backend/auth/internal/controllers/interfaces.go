package controllers

import (
	"mime/multipart"

	"gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/models"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt"
)

type userService interface {
	GetByID(user_id int) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	IsPasswordCorrect(user *models.User, password string) bool
	GetAll(params map[string]string) (*[]models.User, error)
	IsActive(user_id int) (bool, error)
	IsAdmin(user_id int) (bool, error)
	Create(user *models.User) error
	Update(user *models.User) error
	UpdateActive(user_id int, is_active bool) error
}

type jwtUseCase interface {
	RefreshTokens(ctx fiber.Ctx, user *models.User) error
	// GetAccessToken(ctx fiber.Ctx) (*jwt.Token, error)
	// GetRefreshToken(ctx fiber.Ctx) (*jwt.Token, error)
	ClearJWTCookies(ctx fiber.Ctx) error
	GetTokenUserId(token *jwt.Token) (int, error)
	GetUserIdFromToken(ctx fiber.Ctx, use_refresh_token bool) (int, error)
}

type userSessionUseCase interface {
	ObtainUserSessionFromBody(ctx fiber.Ctx, user_id int) (*models.UserSession, error)
	CreateNewLoginSession(ctx fiber.Ctx, user_id int) error
	LogoutLastSession(ctx fiber.Ctx, user_id int) error
	UpdateNoLogoutSession(ctx fiber.Ctx, user_id int) error

	GetByUserID(ctx fiber.Ctx, user_id int, params *models.UserSessionParams) (*[]models.UserSession, error)
	GetByID(ctx fiber.Ctx, id int) (*models.UserSession, error)
	GetLastByUserId(ctx fiber.Ctx, user_id int) (*models.UserSession, error)
	Create(ctx fiber.Ctx, user_id int) (*models.UserSession, error)
	Update(ctx fiber.Ctx, user_session *models.UserSession) error
}

type AuthMiddleware interface {
	IsActive(ctx fiber.Ctx) error
	IsAdmin(ctx fiber.Ctx) error
}

type scheduleService interface {
	GetAll(*models.SchedulesParams) (*[]models.Schedule, error)
	GetByID(schedule_id int) (*models.Schedule, error)
	UpdateConfirmed(schedule_id int, set_confirmed bool) error
	UpdateByID(schedule_id int, params *models.ScheduleUpdateParams) error
	ApplyChangesFromSCV(src *multipart.File) (models.SchedulesUploadResult, error) // successful, duplicated, missing fields cnt
	SearchFlights(params *models.SearchFlightsParams) ([][]*models.Schedule, error)
}

type airportService interface {
	GetAll() (*[]models.Airport, error)
}

type ticketService interface {
	CountRemainingSeats(params *models.TicketsCountRemainingSeatsParams) (*models.TicketsRemainingSeatsCount, error)
	BookTickets(params *models.TicketsBookingParams) ([]*models.Ticket, error)
	ChangeTicketsStatus(ticket_ids []int, set_confirmed bool) error
}

type CountryService interface {
	GetByID(id int) (*models.Country, error)
	GetByName(title string) (*models.Country, error)
	GetAll() ([]*models.Country, error)
}
