package services

import (
	"example/evolza/models"
	"example/evolza/repository"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

type FileService struct {
	fileRepo  *repository.FileRepository
	uploadDir string
}

func NewFileService(fileRepo *repository.FileRepository) *FileService {
	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = "./uploads"
	}

	// Create upload directory if it doesn't exist
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		fmt.Printf("Error creating upload directory: %v\n", err)
	}

	return &FileService{
		fileRepo:  fileRepo,
		uploadDir: uploadDir,
	}
}

func (s *FileService) SaveFile(file *multipart.FileHeader, userID primitive.ObjectID, isPublic bool, tags []string) (*models.FileMetadata, error) {
	// Generate a unique filename
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	filePath := filepath.Join(s.uploadDir, filename)

	// Save the file to disk
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return nil, err
	}

	// Create file metadata
	metadata := &models.FileMetadata{
		FileName:     filename,
		OriginalName: file.Filename,
		FilePath:     filePath,
		FileSize:     file.Size,
		ContentType:  file.Header.Get("Content-Type"),
		UploadedBy:   userID,
		UploadedAt:   time.Now(),
		IsPublic:     isPublic,
		Tags:         tags,
	}

	// Save metadata to database
	if err := s.fileRepo.SaveFileMetadata(metadata); err != nil {
		// Clean up the file if database operation fails
		os.Remove(filePath)
		return nil, err
	}

	return metadata, nil
}

func (s *FileService) GetFile(fileID primitive.ObjectID, userID primitive.ObjectID) (*models.FileMetadata, error) {
	file, err := s.fileRepo.GetFileByID(fileID)
	if err != nil {
		return nil, err
	}

	// Check if user has access to the file
	if !file.IsPublic && file.UploadedBy != userID {
		return nil, fmt.Errorf("unauthorized access to file")
	}

	return file, nil
}

func (s *FileService) DeleteFile(fileID primitive.ObjectID, userID primitive.ObjectID) error {
	file, err := s.fileRepo.GetFileByID(fileID)
	if err != nil {
		return err
	}

	// Check if user has permission to delete the file
	if file.UploadedBy != userID {
		return fmt.Errorf("unauthorized to delete this file")
	}

	// Delete file from filesystem
	if err := os.Remove(file.FilePath); err != nil {
		return err
	}

	// Delete metadata from database
	return s.fileRepo.DeleteFile(fileID)
}

func (s *FileService) GetUserFiles(userID primitive.ObjectID) ([]models.FileMetadata, error) {
	return s.fileRepo.GetFilesByUser(userID)
}
