package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/bowe99/phone-usage-service/internal/domain/model"
	"github.com/bowe99/phone-usage-service/internal/domain/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoDailyUsageRepository struct {
	collection *mongo.Collection
}

func SetupDailyUsageRepository(db *mongo.Database) repository.DailyUsageRepository {
	return &mongoDailyUsageRepository{
		collection: db.Collection("daily_usage"),
	}
}

func (m *mongoDailyUsageRepository) Create(ctx context.Context, usage *model.DailyUsage) error {
	usage.CreatedAt = time.Now()
	usage.UpdatedAt = time.Now()

	result, err := m.collection.InsertOne(ctx, usage)
	if err != nil {
		return fmt.Errorf("failed to create usage: %w", err)
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		usage.ID = oid.Hex()
	}

	return nil
}

func (m *mongoDailyUsageRepository) GetByDateRange(ctx context.Context, userID, mdn string, startDate, endDate time.Time) ([]*model.DailyUsage, error) {
	filter := bson.M{
		"userId": userID,
		"mdn":    mdn,
		"usageDate": bson.M{
			"$gte": startDate,
			"$lte": endDate,
		},
	}

	opts := options.Find().SetSort(bson.D{{Key: "usageDate", Value: 1}})

	cursor, err := m.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get usage by date range: %w", err)
	}
	defer cursor.Close(ctx)

	var usageRecords []*model.DailyUsage
	if err := cursor.All(ctx, &usageRecords); err != nil {
		return nil, fmt.Errorf("failed to decode usage records: %w", err)
	}

	return usageRecords, nil
}

func (m *mongoDailyUsageRepository) Update(ctx context.Context, usage *model.DailyUsage) error {
	objectID, err := primitive.ObjectIDFromHex(usage.ID)
	if err != nil {
		return fmt.Errorf("invalid usage ID: %w", err)
	}

	usage.UpdatedAt = time.Now()

	update := bson.M{
		"$set": bson.M{
			"usedInMb":  usage.UsedInMB,
			"updatedAt": usage.UpdatedAt,
		},
	}

	result, err := m.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return fmt.Errorf("failed to update usage: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("usage record not found")
	}

	return nil
}
