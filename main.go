package main

import (
	"example/evolza/database"
	handlers "example/evolza/handlers"
	"example/evolza/repository"
	"example/evolza/routes"
	services "example/evolza/service"
	"example/evolza/utils"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Connect to MongoDB
	database.Connect()

	// Initialize repositories
	userRepo := repository.NewUserRepository()
	fileRepo := repository.NewFileRepository()

	// Initialize services
	fileService := services.NewFileService(fileRepo)

	// Initialize logger
	appLogger := utils.NewLogger("./logs/app.log")

	// Initialize handlers
	adminHandler := handlers.NewAdminHandler(userRepo)
	fileHandler := handlers.NewFileHandler(fileService, appLogger)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New())

	// Setup routes
	routes.SetupRoutes(app, adminHandler, fileHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
