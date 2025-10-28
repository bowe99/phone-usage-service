package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
)

type MongoDB struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func Connect(uri, database string, timeout time.Duration) (*MongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	clientOptions := options.Client().
		ApplyURI(uri).
		SetMaxPoolSize(100).
		SetMinPoolSize(10).
		SetMaxConnIdleTime(30 * time.Second)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	db := client.Database(database)

	mongodb := &MongoDB{
		Client:   client,
		Database: db,
	}

	if err := mongodb.createIndexes(ctx); err != nil {
		return nil, fmt.Errorf("failed to create indexes: %w", err)
	}

	return mongodb, nil
}

func (m *MongoDB) createIndexes(ctx context.Context) error {
	userIndexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	}
	if _, err := m.Database.Collection("users").Indexes().CreateMany(ctx, userIndexes); err != nil {
		return fmt.Errorf("failed to create user indexes: %w", err)
	}

	cycleIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "userId", Value: 1},
				{Key: "startDate", Value: -1},
			},
		},
		{
			Keys: bson.D{
				{Key: "mdn", Value: 1},
				{Key: "startDate", Value: -1},
			},
		},
		{
			Keys: bson.D{
				{Key: "userId", Value: 1},
				{Key: "mdn", Value: 1},
				{Key: "startDate", Value: 1},
				{Key: "endDate", Value: 1},
			},
		},
	}
	if _, err := m.Database.Collection("cycles").Indexes().CreateMany(ctx, cycleIndexes); err != nil {
		return fmt.Errorf("failed to create cycle indexes: %w", err)
	}

	usageIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "userId", Value: 1},
				{Key: "mdn", Value: 1},
				{Key: "usageDate", Value: -1},
			},
		},
		{
			Keys: bson.D{{Key: "usageDate", Value: 1}},
		},
	}
	if _, err := m.Database.Collection("daily_usage").Indexes().CreateMany(ctx, usageIndexes); err != nil {
		return fmt.Errorf("failed to create usage indexes: %w", err)
	}

	return nil
}

func (m *MongoDB) Disconnect(ctx context.Context) error {
	return m.Client.Disconnect(ctx)
}

func (m *MongoDB) HealthCheck(ctx context.Context) error {
	return m.Client.Ping(ctx, nil)
}
