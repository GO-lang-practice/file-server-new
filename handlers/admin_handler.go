package handlers

import (
	"example/evolza/repository"

	"github.com/gofiber/fiber/v2"
)

type AdminHandler struct {
	userRepo *repository.UserRepository
}

func NewAdminHandler(userRepo *repository.UserRepository) *AdminHandler {
	return &AdminHandler{
		userRepo: userRepo,
	}
}

// GetAllUsers returns all users (admin only)
func (h *AdminHandler) GetAllUsers(c *fiber.Ctx) error {
	users, err := h.userRepo.GetAllUsers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve users: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"users": users,
	})
}
