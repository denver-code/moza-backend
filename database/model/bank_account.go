package model

import (
	"time"

	"gorm.io/gorm"
)

// Currency type for handling different currencies
type Currency string

const (
	USD Currency = "USD"
	EUR Currency = "EUR"
	GBP Currency = "GBP"
)

// AccountType represents the type of bank account
type AccountType string

const (
	CHECKING AccountType = "CHECKING"
	SAVINGS  AccountType = "SAVINGS"
	BUSINESS AccountType = "BUSINESS"
)

// BankAccount represents a user's bank account
type BankAccount struct {
	gorm.Model
	UserID      uint        `gorm:"not null" json:"user_id"`
	AccountType AccountType `gorm:"not null" json:"account_type"`
	Currency    Currency    `gorm:"not null" json:"currency"`
	Balance     float64     `gorm:"type:decimal(20,2);not null;default:0.00" json:"balance"`
	AccountNumber string    `gorm:"uniqueIndex;not null" json:"account_number"`
	IsActive    bool       `gorm:"not null;default:true" json:"is_active"`
	LastActivity time.Time  `json:"last_activity"`
}

// Card represents a payment card associated with a bank account
type Card struct {
	gorm.Model
	UserID        uint      `gorm:"not null" json:"user_id"`
	BankAccountID uint      `gorm:"not null" json:"bank_account_id"`
	CardNumber    string    `gorm:"uniqueIndex;not null" json:"card_number"`
	ExpiryDate    time.Time `gorm:"not null" json:"expiry_date"`
	CVV          string    `gorm:"not null" json:"-"` // CVV is hidden in JSON responses
	IsActive     bool      `gorm:"not null;default:true" json:"is_active"`
	DailyLimit   float64   `gorm:"type:decimal(20,2);not null" json:"daily_limit"`
	CardType     string    `gorm:"not null" json:"card_type"` // VISA, MASTERCARD, etc.
}

// Transaction represents a financial transaction
type Transaction struct {
	gorm.Model
	FromAccountID uint    `gorm:"not null" json:"from_account_id"`
	ToAccountID   uint    `gorm:"not null" json:"to_account_id"`
	Amount        float64 `gorm:"type:decimal(20,2);not null" json:"amount"`
	Currency      Currency `gorm:"not null" json:"currency"`
	Description   string  `json:"description"`
	Type         string  `gorm:"not null" json:"type"` // TRANSFER, DEPOSIT, WITHDRAWAL
	Status       string  `gorm:"not null" json:"status"` // PENDING, COMPLETED, FAILED
	Reference    string  `gorm:"uniqueIndex;not null" json:"reference"`
} 