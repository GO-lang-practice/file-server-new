package handelers

import (
	"example/evolza/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AdminHandler struct {
	userRepo *repository.UserRepository
}

func NewAdminHandler(userRepo *repository.UserRepository) *AdminHandler {
	return &AdminHandler{
		userRepo: userRepo,
	}
}

func (h *AdminHandler) GetAllUsers(c *gin.Context) {
	userID := c.MustGet("userID").(primitive.ObjectID)
}
