package middlewares

import (
	"strings"

	"github.com/Dmitro44/relper/internal/services/auth"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func RequireAuth(authService *auth.AuthService) fiber.Handler {
	return func(c fiber.Ctx) error {
		authField := c.Get("Authorization")
		if authField == "" {
			return fiber.ErrUnauthorized
		}
		token := strings.TrimPrefix(authField, "Bearer ")

		claims, err := authService.ValidateAccessToken(c.Context(), token)
		if err != nil {
			return fiber.ErrUnauthorized
		}

		userID, err := uuid.Parse(claims.Subject)
		if err != nil {
			return fiber.ErrUnauthorized
		}
		c.Locals("userID", userID)

		return c.Next()
	}
}

func UserID(c fiber.Ctx) (uuid.UUID, bool) {
	userID, ok := c.Locals("userID").(uuid.UUID)
	return userID, ok
}
