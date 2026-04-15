package middleware

import (
	"strconv"
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
		orgID := c.Get("Organization-ID")
		if orgID == "" || !utils.IsValidInt64(orgID) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"code":    fiber.StatusUnauthorized,
				"message": "unauthorized",
				"error":   "missing organization header",
				"payload": nil,
			})
		}

		parsedOrgID, _ := strconv.ParseInt(orgID, 10, 64)
		c.Locals(LocalsOrgID, parsedOrgID)
		return c.Next()
	}
}
