package utils

import (
	"example/evolza/database"
	"example/evolza/models"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os"
	"path/filepath"
	"time"
)

type Logger struct {
	filePath string
}

func NewLogger(filePath string) *Logger {
	// Create logs directory if it doesn't exist
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Printf("Error creating log directory: %v\n", err)
	}

	return &Logger{
		filePath: filePath,
	}
}

// LogOperation logs an operation to both file and database
func (l *Logger) LogOperation(userID primitive.ObjectID, username, operation, resource, details, ipAddress string, success bool) {
	timestamp := time.Now()

	// Log to database
	logEntry := &models.OperationLog{
		UserID:    userID,
		Username:  username,
		Operation: operation,
		Resource:  resource,
		Details:   details,
		IPAddress: ipAddress,
		Timestamp: timestamp,
		Success:   success,
	}

	collection := database.Database.Collection("operation_logs")
	collection.InsertOne(nil, logEntry)

	// Log to file
	logMessage := fmt.Sprintf("[%s] User: %s (%s) | Operation: %s | Resource: %s | Details: %s | IP: %s | Success: %t\n",
		timestamp.Format("2006-01-02 15:04:05"),
		username,
		userID.Hex(),
		operation,
		resource,
		details,
		ipAddress,
		success,
	)

	l.writeToFile(logMessage)
}

// LogAuthentication logs authentication attempts
func (l *Logger) LogAuthentication(username, ipAddress string, success bool) {
	timestamp := time.Now()
	action := "LOGIN_SUCCESS"
	if !success {
		action = "LOGIN_FAILED"
	}

	logMessage := fmt.Sprintf("[%s] Authentication: %s | User: %s | IP: %s\n",
		timestamp.Format("2006-01-02 15:04:05"),
		action,
		username,
		ipAddress,
	)

	l.writeToFile(logMessage)
}

// LogError logs errors to file
func (l *Logger) LogError(userID primitive.ObjectID, username, operation, errorMsg, ipAddress string) {
	timestamp := time.Now()

	logMessage := fmt.Sprintf("[%s] ERROR | User: %s (%s) | Operation: %s | Error: %s | IP: %s\n",
		timestamp.Format("2006-01-02 15:04:05"),
		username,
		userID.Hex(),
		operation,
		errorMsg,
		ipAddress,
	)

	l.writeToFile(logMessage)
}

func (l *Logger) writeToFile(message string) {
	file, err := os.OpenFile(l.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error opening log file: %v\n", err)
		return
	}
	defer file.Close()

	if _, err := file.WriteString(message); err != nil {
		fmt.Printf("Error writing to log file: %v\n", err)
	}
}
