package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/bowe99/phone-usage-service/internal/domain/model"
	"github.com/bowe99/phone-usage-service/internal/domain/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user with this email already exists")
)

type mongoUserRepository struct {
	collection *mongo.Collection
}

func SetupUserRepository(db *mongo.Database) repository.UserRepository {
	return &mongoUserRepository{
		collection: db.Collection("users"),
	}
}

func (m *mongoUserRepository) Create(ctx context.Context, user *model.User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	result, err := m.collection.InsertOne(ctx, user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return ErrUserAlreadyExists
		}
		return fmt.Errorf("failed to create user: %w", err)
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		user.ID = oid.Hex()
	}
	return nil
}

func (m *mongoUserRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrUserNotFound
	}

	result, err := m.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if result.DeletedCount == 0 {
		return ErrUserNotFound
	}

	return nil}

func (m *mongoUserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := m.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &user, nil
}

func (m *mongoUserRepository) GetByID(ctx context.Context, id string) (*model.User, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, ErrUserNotFound
	}

	var user model.User
	err = m.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil{
		if err == mongo.ErrNoDocuments {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (m *mongoUserRepository) Update(ctx context.Context, user *model.User) error {
	objectID, err := primitive.ObjectIDFromHex(user.ID)
	if err != nil {
		return ErrUserNotFound
	}

	user.UpdatedAt = time.Now()

	update := bson.M{
		"$set": bson.M{
			"firstName": user.FirstName,
			"lastName":  user.LastName,
			"email":     user.Email,
			"updatedAt": user.UpdatedAt,
		},
	}

	result, err := m.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return ErrUserAlreadyExists
		}
		return fmt.Errorf("failed to update user: %w", err)
	}

	if result.MatchedCount == 0 {
		return ErrUserNotFound
	}

	return nil
}
