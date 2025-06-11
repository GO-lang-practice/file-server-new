package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
	RoleGuest Role = "guest"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username  string             `bson:"username" json:"username" validate:"required"`
	Email     string             `bson:"email" json:"email" validate:"required"`
	Password  string             `bson:"password" json:"-"`
	Role      Role               `bson:"role" json:"role"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	IsActive  bool               `bson:"is_active" json:"is_active"`
}

type FileMetadata struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FileName     string             `bson:"file_name" json:"file_name"`
	OriginalName string             `bson:"original_name" json:"original_name"`
	FilePath     string             `bson:"file_path" json:"file_path"`
	FileSize     int64              `bson:"file_size" json:"file_size"`
	ContentType  string             `bson:"content_type" json:"content_type"`
	UploadedBy   primitive.ObjectID `bson:"uploaded_by" json:"uploaded_by"`
	UploadedAt   time.Time          `bson:"uploaded_at" json:"uploaded_at"`
	IsPublic     bool               `bson:"is_public" json:"is_public"`
	Tags         []string           `bson:"tags" json:"tags"`
}

type OperationLog struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
	Username  string             `bson:"username" json:"username"`
	Operation string             `bson:"operation" json:"operation"`
	Resource  string             `bson:"resource" json:"resource"`
	Details   string             `bson:"details" json:"details"`
	IPAddress string             `bson:"ip_address" json:"ip_address"`
	Timestamp time.Time          `bson:"timestamp" json:"timestamp"`
	Success   bool               `bson:"success" json:"success"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     Role   `json:"role"`
}
