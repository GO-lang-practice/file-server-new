package routes

import (
	handlers "example/evolza/handlers"
	"example/evolza/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes configures all application routes
func SetupRoutes(app *fiber.App, adminHandler *handlers.AdminHandler, fileHandler *handlers.FileHandler) {
	// API routes
	api := app.Group("/api")

	// Public routes (no authentication required)
	setupPublicRoutes(api)

	// Protected routes (authentication required)
	protected := api.Group("/", middleware.AuthMiddleware())

	// Admin routes
	setupAdminRoutes(protected, adminHandler)

	// File routes
	setupFileRoutes(protected, fileHandler)
}

// setupPublicRoutes configures public routes that don't require authentication
func setupPublicRoutes(api fiber.Router) {
	// Health check endpoint
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "File Server is running",
		})
	})

	// API version endpoint
	api.Get("/version", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"version": "1.0.0",
			"api":     "File Server API",
		})
	})
}

// setupAdminRoutes configures admin-only routes
func setupAdminRoutes(protected fiber.Router, adminHandler *handlers.AdminHandler) {
	admin := protected.Group("/admin")
	admin.Get("/users", adminHandler.GetAllUsers)
}

// setupFileRoutes configures file management routes
func setupFileRoutes(protected fiber.Router, fileHandler *handlers.FileHandler) {
	files := protected.Group("/files")
	files.Post("/upload", fileHandler.UploadFile)
	files.Get("/", fileHandler.GetUserFiles)
	files.Get("/:id", fileHandler.DownloadFile)
	files.Delete("/:id", fileHandler.DeleteFile)
}
