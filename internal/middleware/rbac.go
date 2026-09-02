package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// RequireRole returns a Fiber middleware that enforces role-based access control.
// It reads the "role" claim previously set by AuthenticateJWT and checks against allowedRoles.
// Usage: RequireRole("ORG_ADMIN", "DIRECTOR")
func RequireRole(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole, ok := c.Locals("role").(string)
		if !ok || userRole == "" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"error":   "MISSING_ROLE",
				"message": "Klaim role tidak ditemukan pada token. Pastikan Anda telah login.",
			})
		}

		for _, allowed := range allowedRoles {
			if allowed == userRole {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"error":   "INSUFFICIENT_PERMISSIONS",
			"message": "Anda tidak memiliki wewenang untuk tindakan ini.",
		})
	}
}
