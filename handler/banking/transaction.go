package banking

import (
	"github.com/denver-code/moza-backend/database"
	"github.com/denver-code/moza-backend/model"
	"github.com/denver-code/moza-backend/util"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// Transfer handles money transfer between accounts
func Transfer(c *fiber.Ctx) error {
	type TransferInput struct {
		FromAccountID uint    `json:"from_account_id"`
		ToAccountID   uint    `json:"to_account_id"`
		Amount        float64 `json:"amount"`
		Description   string  `json:"description"`
	}

	input := new(TransferInput)
	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input",
			"data":    err.Error(),
		})
	}

	if input.FromAccountID == input.ToAccountID {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Cannot transfer to the same account",
		})
	}

	// Get user ID from JWT token
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := uint(claims["user_id"].(float64))

	// Start transaction
	tx := database.DB.Begin()

	// Verify ownership and get source account
	var fromAccount model.BankAccount
	if err := tx.Where("id = ? AND user_id = ?", input.FromAccountID, userID).First(&fromAccount).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Unauthorized or account not found",
		})
	}

	// Check sufficient balance
	if fromAccount.Balance < input.Amount {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Insufficient balance",
		})
	}

	// Get destination account
	var toAccount model.BankAccount
	if err := tx.Where("id = ?", input.ToAccountID).First(&toAccount).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Destination account not found",
		})
	}

	// Update balances
	if err := tx.Model(&fromAccount).Update("balance", fromAccount.Balance-input.Amount).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Could not update source account",
		})
	}

	if err := tx.Model(&toAccount).Update("balance", toAccount.Balance+input.Amount).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Could not update destination account",
		})
	}

	// Create transaction record
	transaction := &model.Transaction{
		FromAccountID: input.FromAccountID,
		ToAccountID:   input.ToAccountID,
		Amount:        input.Amount,
		Currency:      fromAccount.Currency,
		Description:   input.Description,
		Type:          "TRANSFER",
		Status:        "COMPLETED",
		Reference:     util.GenerateTransactionReference(),
	}

	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Could not create transaction record",
		})
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Could not complete transfer",
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Transfer completed successfully",
		"data":    transaction,
	})
}

// GetAccountTransactions retrieves transactions for a specific account
func GetAccountTransactions(c *fiber.Ctx) error {
	accountID := c.Params("id")

	// Get user ID from JWT token
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := uint(claims["user_id"].(float64))

	// Verify account ownership
	var account model.BankAccount
	if err := database.DB.Where("id = ? AND user_id = ?", accountID, userID).First(&account).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Unauthorized or account not found",
		})
	}

	var transactions []model.Transaction
	if err := database.DB.Where("from_account_id = ? OR to_account_id = ?", accountID, accountID).
		Order("created_at desc").
		Find(&transactions).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Could not retrieve transactions",
			"data":    err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Transactions retrieved successfully",
		"data":    transactions,
	})
}
