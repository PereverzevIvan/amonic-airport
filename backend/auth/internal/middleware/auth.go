package middleware

import (
	"github.com/gofiber/fiber/v3"
)

type authUseCase interface {
	IsAdmin(user_id int) (bool, error)
}

type authMiddleware struct {
	authUseCase authUseCase
}

func NewAuthMiddleware(authUseCase authUseCase) *authMiddleware {
	return &authMiddleware{authUseCase: authUseCase}
}

func (am *authMiddleware) IsAdmin(ctx fiber.Ctx) error {
	// user_id := fiber.Query[int](ctx, "user_id")
	// is_admin, err := am.authUseCase.IsAdmin(user_id)
	// if err != nil {
	// 	return ctx.SendString(err.Error())
	// }
	// if !is_admin {
	// 	return ctx.SendStatus(403)
	// }
	return ctx.Next()
}
