package banking

import (
	"time"

	"github.com/denver-code/moza-backend/database"
	"github.com/denver-code/moza-backend/database/model"
	"github.com/denver-code/moza-backend/util"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// CreateBankAccount creates a new bank account for the user
func CreateBankAccount(c *fiber.Ctx) error {
	type AccountInput struct {
		AccountType string `json:"account_type"`
		Currency    string `json:"currency"`
	}

	input := new(AccountInput)
	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input",
			"data":    err.Error(),
		})
	}

	// Get user ID from JWT token
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := uint(claims["user_id"].(float64))

	// Generate unique account number
	accountNumber := util.GenerateAccountNumber()

	account := &model.BankAccount{
		UserID:        userID,
		AccountType:   model.AccountType(input.AccountType),
		Currency:      model.Currency(input.Currency),
		AccountNumber: accountNumber,
		Balance:       0,
		IsActive:      true,
		LastActivity:  time.Now(),
	}

	if err := database.DB.Create(&account).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Could not create bank account",
			"data":    err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Bank account created successfully",
		"data":    account,
	})
}

// GetUserAccounts retrieves all bank accounts for the user
func GetUserAccounts(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := uint(claims["user_id"].(float64))

	var accounts []model.BankAccount
	if err := database.DB.Where("user_id = ?", userID).Find(&accounts).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Could not retrieve accounts",
			"data":    err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Accounts retrieved successfully",
		"data":    accounts,
	})
}
