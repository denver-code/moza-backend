package handler

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/denver-code/moza-backend/database"
	"github.com/denver-code/moza-backend/model"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
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
	accountNumber := generateAccountNumber()

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
	cardNumber := generateCardNumber()
	cvv := generateCVV()
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
		Reference:     generateTransactionReference(),
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

// Helper functions for generating numbers and references
func generateAccountNumber() string {
	return fmt.Sprintf("%d", rand.Intn(900000000)+100000000)
}

func generateCardNumber() string {
	return fmt.Sprintf("%d", rand.Intn(9000000000000000)+1000000000000000)
}

func generateCVV() string {
	return fmt.Sprintf("%d", rand.Intn(900)+100)
}

func generateTransactionReference() string {
	return fmt.Sprintf("TXN%d", time.Now().UnixNano())
}
