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
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrCycleNotFound = errors.New("cycle not found")
	ErrNoCycleActive = errors.New("no active cycle found")
)

type mongoCycleRepository struct {
	collection *mongo.Collection
}

func SetupCycleRepository(db *mongo.Database) repository.CycleRepository {
	return &mongoCycleRepository{
		collection: db.Collection("cycles"),
	}
}

func (m *mongoCycleRepository) Create(ctx context.Context, cycle *model.Cycle) error {
	cycle.CreatedAt = time.Now()

	result, err := m.collection.InsertOne(ctx, cycle)
	if err != nil {
		return fmt.Errorf("failed to create cycle: %w", err)
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		cycle.ID = oid.Hex()
	}

	return nil
}

func (m *mongoCycleRepository) GetByID(ctx context.Context, id string) (*model.Cycle, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, ErrCycleNotFound
	}

	var cycle model.Cycle
	err = m.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&cycle)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrCycleNotFound
		}
		return nil, fmt.Errorf("failed to get cycle: %w", err)
	}

	return &cycle, nil
}

// Note: Query by MDN, not userId, because MDNs can transfer between users
func (m *mongoCycleRepository) GetByMDN(ctx context.Context, mdn string) ([]*model.Cycle, error) {
	opts := options.Find().SetSort(bson.D{{Key: "startDate", Value: -1}})

	cursor, err := m.collection.Find(ctx, bson.M{"mdn": mdn}, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get cycles by MDN: %w", err)
	}
	defer cursor.Close(ctx)

	var cycles []*model.Cycle
	if err := cursor.All(ctx, &cycles); err != nil {
		return nil, fmt.Errorf("failed to decode cycles: %w", err)
	}

	return cycles, nil
}

func (m *mongoCycleRepository) GetByUserID(ctx context.Context, userID string) ([]*model.Cycle, error) {
	opts := options.Find().SetSort(bson.D{{Key: "startDate", Value: -1}})

	cursor, err := m.collection.Find(ctx, bson.M{"userId": userID}, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get cycles by user ID: %w", err)
	}
	defer cursor.Close(ctx)

	var cycles []*model.Cycle
	if err := cursor.All(ctx, &cycles); err != nil {
		return nil, fmt.Errorf("failed to decode cycles: %w", err)
	}

	return cycles, nil
}

func (r *mongoCycleRepository) GetCurrentCycle(ctx context.Context, userID, mdn string, currentDate time.Time) (*model.Cycle, error) {
	filter := bson.M{
		"userId":    userID,
		"mdn":       mdn,
		"startDate": bson.M{"$lte": currentDate},
		"endDate":   bson.M{"$gte": currentDate},
	}

	var cycle model.Cycle
	err := r.collection.FindOne(ctx, filter).Decode(&cycle)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrNoCycleActive
		}
		return nil, fmt.Errorf("failed to get current cycle: %w", err)
	}

	return &cycle, nil
}
