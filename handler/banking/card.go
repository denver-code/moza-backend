package banking

import (
	"errors"
	"time"

	"github.com/denver-code/moza-backend/database"
	"github.com/denver-code/moza-backend/model"
	"github.com/denver-code/moza-backend/util"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

// CreateCard creates a new card for a bank account
func CreateCard(c *fiber.Ctx) error {
	type CardInput struct {
		BankAccountID uint    `json:"bank_account_id"`
		CardType      string  `json:"card_type"`
		DailyLimit    float64 `json:"daily_limit"`
	}

	input := new(CardInput)
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

	// Verify bank account ownership
	var account model.BankAccount
	if err := database.DB.Where("id = ? AND user_id = ?", input.BankAccountID, userID).First(&account).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Bank account not found or unauthorized",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Could not verify bank account",
			"data":    err.Error(),
		})
	}

	// Generate card details
	cardNumber := util.GenerateCardNumber()
	cvv := util.GenerateCVV()
	expiryDate := time.Now().AddDate(4, 0, 0) // 4 years validity

	card := &model.Card{
		UserID:        userID,
		BankAccountID: input.BankAccountID,
		CardNumber:    cardNumber,
		ExpiryDate:    expiryDate,
		CVV:           cvv,
		IsActive:      true,
		DailyLimit:    input.DailyLimit,
		CardType:      input.CardType,
	}

	if err := database.DB.Create(&card).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Could not create card",
			"data":    err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Card created successfully",
		"data":    card,
	})
}

// GetCards retrieves all cards for the user bank account
func GetCards(c *fiber.Ctx) error {
	accountID := c.Params("id")

	// Get user ID from JWT token
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := uint(claims["user_id"].(float64))

	// Verify bank account ownership
	var account model.BankAccount
	if err := database.DB.Where("id = ? AND user_id = ?", accountID, userID).First(&account).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Bank account not found or unauthorized",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Could not verify bank account",
			"data":    err.Error(),
		})
	}

	var cards []model.Card
	if err := database.DB.Where("bank_account_id = ?", accountID).Find(&cards).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Could not retrieve cards",
			"data":    err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Cards retrieved successfully",
		"data":    cards,
	})
}
