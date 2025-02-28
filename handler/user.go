package handler

import (
	"github.com/denver-code/moza-backend/database"
	"github.com/denver-code/moza-backend/database/model"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func GetProfile(c *fiber.Ctx) error {
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := uint(claims["user_id"].(float64))

	db := database.DB
	var user model.User
	if err := db.First(&user, userID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "User not found", "data": nil})
	}

	// Don't send sensitive information
	user.Password = ""
	return c.JSON(fiber.Map{"status": "success", "message": "Profile retrieved successfully", "data": user})
}
