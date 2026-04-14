package middleware

import (
	"strings"

	"aprilpollo/internal/utils"

	"github.com/gofiber/fiber/v2"
)

const (
	LocalsUserID = "user_id"
	LocalsEmail  = "email"
	LocalsOrgID  = "organization_id"
)

func JWTProtected(secretKey string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"code":    fiber.StatusUnauthorized,
				"message": "unauthorized",
				"error":   "missing or invalid authorization header",
				"payload": nil,
			})
		}

		tokenStr := strings.TrimPrefix(auth, "Bearer ")
		claims, err := utils.ParseToken(tokenStr, secretKey)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"code":    fiber.StatusUnauthorized,
				"message": "unauthorized",
				"error":   "invalid or expired token",
				"payload": nil,
			})
		}

		c.Locals(LocalsUserID, claims.UserID)
		c.Locals(LocalsEmail, claims.Email)

		return c.Next()
	}
}

func OrganizationProtected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		org_id := c.Get("Organization-ID")
		if org_id == "" || !utils.IsValidInt64(org_id) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"code":    fiber.StatusUnauthorized,
				"message": "unauthorized",
				"error":   "missing organization header",
				"payload": nil,
			})
		}
		c.Locals(LocalsOrgID, org_id)
		return c.Next()
	}
}
