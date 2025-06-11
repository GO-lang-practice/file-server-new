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

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		collection: database.Database.Collection("users"),
	}
}

func (r *UserRepository) Create(user *models.User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.IsActive = true

	result, err := r.collection.InsertOne(context.Background(), user)
	if err != nil {
		return err
	}
	user.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.collection.FindOne(context.Background(), bson.M{"username": username}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Update(user *models.User) error {
	user.UpdatedAt = time.Now()
	_, err := r.collection.UpdateOne(
		context.Background(),
		bson.M{"_id": user.ID},
		bson.M{"$set": user},
	)
	return err
}

func (r *UserRepository) GetAllUsers() ([]models.User, error) {
	var users []models.User
	cursor, err := r.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
