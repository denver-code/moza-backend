package router

import (
	"github.com/denver-code/moza-backend/handler"
	"github.com/denver-code/moza-backend/handler/banking"
	"github.com/denver-code/moza-backend/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// SetupRoutes setup router api
func SetupRoutes(app *fiber.App) {
	// Middleware
	api := app.Group("/api", logger.New())
	api.Get("/", handler.Hello)

	// Auth
	auth := api.Group("/auth")
	auth.Post("/login", handler.Login)
	auth.Post("/register", handler.Register)

	// Private
	private := api.Group("/private")
	private.Get("/", middleware.Protected(), handler.ProtectedTest)

	// User
	user := api.Group("/user")
	user.Get("/profile", middleware.Protected(), handler.GetProfile)

	// Banking
	banking_group := api.Group("/banking")
	banking_group.Use(middleware.Protected()) // All banking routes require authentication

	// Bank Accounts
	banking_group.Post("/accounts", banking.CreateBankAccount)
	banking_group.Get("/accounts", banking.GetUserAccounts)
	banking_group.Get("/accounts/:id/transactions", banking.GetAccountTransactions)

	// Cards
	banking_group.Post("/cards", banking.CreateCard)
	banking_group.Get("/accounts/:id/cards", banking.GetCards)

	// Transactions
	banking_group.Post("/transfer", banking.Transfer)

}
