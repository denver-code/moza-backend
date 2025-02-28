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
	user.Get("/:id", handler.GetUser)
	// user.Post("/", handler.CreateUser)
	// user.Patch("/:id", middleware.Protected(), handler.UpdateUser)
	// user.Delete("/:id", middleware.Protected(), handler.DeleteUser)

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

	// Product
	product := api.Group("/product")
	product.Get("/", handler.GetAllProducts)
	product.Get("/:id", handler.GetProduct)
	product.Post("/", middleware.Protected(), handler.CreateProduct)
	product.Delete("/:id", middleware.Protected(), handler.DeleteProduct)
}
