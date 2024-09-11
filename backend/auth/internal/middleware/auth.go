package middleware

import (
	"github.com/gofiber/fiber/v3"
)

type authService interface {
	IsAdmin(user_id int) (bool, error)
}

type authMiddleware struct {
	authService authService
}

func NewAuthMiddleware(authService authService) *authMiddleware {
	return &authMiddleware{authService: authService}
}

func (am *authMiddleware) IsAdmin(ctx fiber.Ctx) error {
	// user_id := fiber.Query[int](ctx, "user_id")
	// is_admin, err := am.authService.IsAdmin(user_id)
	// if err != nil {
	// 	return ctx.SendString(err.Error())
	// }
	// if !is_admin {
	// 	return ctx.SendStatus(403)
	// }
	return ctx.Next()
}
