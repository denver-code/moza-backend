package handler

import "github.com/gofiber/fiber/v2"

// Hello handle api status
func ProtectedTest(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "success", "message": "Protection Working", "data": nil})
}
