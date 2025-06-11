package repository

import (
	"context"
	"example/evolza/database"
	"example/evolza/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type FileRepository struct {
	collection *mongo.Collection
}

func NewFileRepository() *FileRepository {
	return &FileRepository{
		collection: database.Database.Collection("files"),
	}
}

func (r *FileRepository) SaveFileMetadata(metadata *models.FileMetadata) error {
	metadata.UploadedAt = time.Now()

	result, err := r.collection.InsertOne(context.Background(), metadata)
	if err != nil {
		return err
	}
	metadata.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *FileRepository) GetFileByID(id primitive.ObjectID) (*models.FileMetadata, error) {
	var file models.FileMetadata
	err := r.collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&file)
	if err != nil {
		return nil, err
	}
	return &file, nil
}

func (r *FileRepository) GetFilesByUser(userID primitive.ObjectID) ([]models.FileMetadata, error) {
	var files []models.FileMetadata
	cursor, err := r.collection.Find(context.Background(), bson.M{"uploaded_by": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var file models.FileMetadata
		if err := cursor.Decode(&file); err != nil {
			return nil, err
		}
		files = append(files, file)
	}
	return files, nil
}

func (r *FileRepository) DeleteFile(id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(context.Background(), bson.M{"_id": id})
	return err
}
