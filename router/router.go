package router

import (
	"github.com/denver-code/moza-backend/handler"
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
	banking := api.Group("/banking")
	banking.Use(middleware.Protected()) // All banking routes require authentication

	// Bank Accounts
	banking.Post("/accounts", handler.CreateBankAccount)
	banking.Get("/accounts", handler.GetUserAccounts)
	banking.Get("/accounts/:id/transactions", handler.GetAccountTransactions)

	// Cards
	banking.Post("/cards", handler.CreateCard)
	banking.Get("/:id/cards", handler.GetCards)

	// Transactions
	banking.Post("/transfer", handler.Transfer)

}
