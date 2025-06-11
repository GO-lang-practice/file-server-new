package utils

import (
	"context"
	"example/evolza/database"
	"example/evolza/models"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Logger struct {
	filePath string
}

func NewLogger(filePath string) *Logger {
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Printf("Error creating log directory: %v\n", err)
	}

	return &Logger{
		filePath: filePath,
	}
}

// LogOperation logs operation to database and file
func (l *Logger) LogOperation(userID primitive.ObjectID, username, operation, resource, details, ipAddress string, success bool) {
	timestamp := time.Now()

	logMessage := fmt.Sprintf("[%s] OPERATION | User: %s (%s) | Operation: %s | Resource: %s | Details: %s | IP: %s | Success: %t",
		timestamp.Format("2006-01-02 15:04:05"),
		username, userID.Hex(), operation, resource, details, ipAddress, success)

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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := collection.InsertOne(ctx, logEntry)
	if err != nil {
		fmt.Printf("Error inserting operation log to database: %v\n", err)
	}

	// Log to file
	l.writeToFile(logMessage)
}

// LogAuthentication logs authentication attempts
func (l *Logger) LogAuthentication(username, ipAddress string, success bool) {
	timestamp := time.Now()
	action := "LOGIN_SUCCESS"
	if !success {
		action = "LOGIN_FAILED"
	}

	logMessage := fmt.Sprintf("[%s] AUTH | Action: %s | User: %s | IP: %s",
		timestamp.Format("2006-01-02 15:04:05"),
		action, username, ipAddress)

	l.writeToFile(logMessage)
}

// LogError logs errors
func (l *Logger) LogError(userID primitive.ObjectID, username, operation, errorMsg, ipAddress string) {
	timestamp := time.Now()

	logMessage := fmt.Sprintf("[%s] ERROR | User: %s (%s) | Operation: %s | Error: %s | IP: %s",
		timestamp.Format("2006-01-02 15:04:05"),
		username, userID.Hex(), operation, errorMsg, ipAddress)

	// Log to database
	errorLogEntry := &models.OperationLog{
		UserID:    userID,
		Username:  username,
		Operation: operation,
		Resource:  "error",
		Details:   errorMsg,
		IPAddress: ipAddress,
		Timestamp: timestamp,
		Success:   false,
	}
	collection := database.Database.Collection("operation_logs")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := collection.InsertOne(ctx, errorLogEntry)
	if err != nil {
		fmt.Printf("Error inserting error log to database: %v\n", err)
	}

	l.writeToFile(logMessage)
}

// LogInfo logs general info messages
func (l *Logger) LogInfo(message string) {
	timestamp := time.Now()
	logMessage := fmt.Sprintf("[%s] INFO | %s", timestamp.Format("2006-01-02 15:04:05"), message)
	l.writeToFile(logMessage)
}

// LogWarning logs warning messages
func (l *Logger) LogWarning(userID primitive.ObjectID, username, operation, warning, ipAddress string) {
	timestamp := time.Now()
	logMessage := fmt.Sprintf("[%s] WARNING | User: %s (%s) | Operation: %s | Warning: %s | IP: %s",
		timestamp.Format("2006-01-02 15:04:05"),
		username, userID.Hex(), operation, warning, ipAddress)
	l.writeToFile(logMessage)
}

// writeToFile handles writing any log message to the log file
func (l *Logger) writeToFile(message string) {
	file, err := os.OpenFile(l.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error opening log file: %v\n", err)
		return
	}
	defer file.Close()

	if _, err := file.WriteString(message + "\n"); err != nil {
		fmt.Printf("Error writing to log file: %v\n", err)
	}
}
