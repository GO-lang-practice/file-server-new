package handlers

import (
	services "example/evolza/service"
	"example/evolza/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FileHandler struct {
	fileService *services.FileService
	logger      *utils.Logger
}

func NewFileHandler(fileService *services.FileService, logger *utils.Logger) *FileHandler {
	return &FileHandler{
		fileService: fileService,
		logger:      logger,
	}
}

func (h *FileHandler) UploadFile(c *fiber.Ctx) error {
	//h=go receiver
	userID := c.Locals("userID").(primitive.ObjectID)
	username := c.Locals("username").(string)

	// Get file from form
	file, err := c.FormFile("file")
	if err != nil {
		h.logger.LogError(userID, username, "UPLOAD_FILE", err.Error(), c.IP())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "File upload failed: " + err.Error(),
		})
	}

	// Parse other form fields
	isPublic := c.FormValue("isPublic") == "true"
	tags := c.FormValue("tags", "")

	// Process tags (simple comma-separated handling)
	var tagSlice []string
	if tags != "" {
		tagSlice = strings.Split(strings.TrimSpace(tags), ",")
		// Trim whitespace from each tag
		for i, tag := range tagSlice {
			tagSlice[i] = strings.TrimSpace(tag)
		}
	}

	// Save the file
	metadata, err := h.fileService.SaveFile(file, userID, isPublic, tagSlice)
	if err != nil {
		h.logger.LogError(userID, username, "UPLOAD_FILE", err.Error(), c.IP())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save file: " + err.Error(),
		})
	}

	h.logger.LogOperation(userID, username, "UPLOAD_FILE", "file", file.Filename, c.IP(), true)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "File uploaded successfully",
		"file":    metadata,
	})
}

func (h *FileHandler) DownloadFile(c *fiber.Ctx) error {
	userID := c.Locals("userID").(primitive.ObjectID)
	username := c.Locals("username").(string)

	fileID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid file ID",
		})
	}

	// Get file metadata
	file, err := h.fileService.GetFile(fileID, userID)
	if err != nil {
		h.logger.LogError(userID, username, "DOWNLOAD_FILE", err.Error(), c.IP())
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "File not found or access denied",
		})
	}

	h.logger.LogOperation(userID, username, "DOWNLOAD_FILE", "file", file.OriginalName, c.IP(), true)

	return c.Download(file.FilePath, file.OriginalName)
}

func (h *FileHandler) GetUserFiles(c *fiber.Ctx) error {
	userID := c.Locals("userID").(primitive.ObjectID)

	files, err := h.fileService.GetUserFiles(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve files: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"files": files,
	})
}

func (h *FileHandler) DeleteFile(c *fiber.Ctx) error {
	userID := c.Locals("userID").(primitive.ObjectID)
	username := c.Locals("username").(string)

	fileID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid file ID",
		})
	}

	if err := h.fileService.DeleteFile(fileID, userID); err != nil {
		h.logger.LogError(userID, username, "DELETE_FILE", err.Error(), c.IP())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete file: " + err.Error(),
		})
	}

	h.logger.LogOperation(userID, username, "DELETE_FILE", "file", fileID.Hex(), c.IP(), true)

	return c.JSON(fiber.Map{
		"message": "File deleted successfully",
	})
}
