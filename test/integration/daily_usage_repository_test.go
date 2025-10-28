package integration

import (
	"context"
	"testing"
	"time"

	"github.com/bowe99/phone-usage-service/internal/domain/model"
	"github.com/bowe99/phone-usage-service/internal/infra/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestDailyUsageRepository_Create(t *testing.T) {
	ctx := context.Background()

	mongoContainer, err := mongodb.Run(ctx, "mongo:6")
	require.NoError(t, err)
	defer mongoContainer.Terminate(ctx)

	connStr, err := mongoContainer.ConnectionString(ctx)
	require.NoError(t, err)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connStr))
	require.NoError(t, err)
	defer client.Disconnect(ctx)

	db := client.Database("test_db")
	repo := repository.SetupDailyUsageRepository(db)

	usage := &model.DailyUsage{
		MDN:       "5551234567",
		UserID:    "user123",
		UsageDate: time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC),
		UsedInMB:  250.5,
	}

	err = repo.Create(ctx, usage)
	assert.NoError(t, err)
	assert.NotEmpty(t, usage.ID)
}

func TestDailyUsageRepository_GetByDateRange(t *testing.T) {
	ctx := context.Background()

	mongoContainer, err := mongodb.Run(ctx, "mongo:6")
	require.NoError(t, err)
	defer mongoContainer.Terminate(ctx)

	connStr, err := mongoContainer.ConnectionString(ctx)
	require.NoError(t, err)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connStr))
	require.NoError(t, err)
	defer client.Disconnect(ctx)

	db := client.Database("test_db")
	repo := repository.SetupDailyUsageRepository(db)

	usageRecords := []*model.DailyUsage{
		{
			MDN:       "5551234567",
			UserID:    "user123",
			UsageDate: time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC),
			UsedInMB:  250.5,
		},
		{
			MDN:       "5551234567",
			UserID:    "user123",
			UsageDate: time.Date(2024, 11, 2, 0, 0, 0, 0, time.UTC),
			UsedInMB:  180.3,
		},
		{
			MDN:       "5551234567",
			UserID:    "user123",
			UsageDate: time.Date(2024, 11, 3, 0, 0, 0, 0, time.UTC),
			UsedInMB:  320.7,
		},
		{
			MDN:       "5551234567",
			UserID:    "user123",
			UsageDate: time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC), // Outside range
			UsedInMB:  400.0,
		},
	}

	for _, record := range usageRecords {
		err := repo.Create(ctx, record)
		require.NoError(t, err)
	}

	startDate := time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 11, 30, 23, 59, 59, 0, time.UTC)

	results, err := repo.GetByDateRange(ctx, "user123", "5551234567", startDate, endDate)
	assert.NoError(t, err)
	assert.Len(t, results, 3) // Should exclude December record

	assert.Equal(t, time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC), results[0].UsageDate)
	assert.Equal(t, time.Date(2024, 11, 2, 0, 0, 0, 0, time.UTC), results[1].UsageDate)
	assert.Equal(t, time.Date(2024, 11, 3, 0, 0, 0, 0, time.UTC), results[2].UsageDate)

	assert.Equal(t, 250.5, results[0].UsedInMB)
	assert.Equal(t, 180.3, results[1].UsedInMB)
	assert.Equal(t, 320.7, results[2].UsedInMB)
}

func TestDailyUsageRepository_Update(t *testing.T) {
	ctx := context.Background()

	mongoContainer, err := mongodb.Run(ctx, "mongo:6")
	require.NoError(t, err)
	defer mongoContainer.Terminate(ctx)

	connStr, err := mongoContainer.ConnectionString(ctx)
	require.NoError(t, err)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connStr))
	require.NoError(t, err)
	defer client.Disconnect(ctx)

	db := client.Database("test_db")
	repo := repository.SetupDailyUsageRepository(db)

	usage := &model.DailyUsage{
		MDN:       "5551234567",
		UserID:    "user123",
		UsageDate: time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC),
		UsedInMB:  250.5,
	}

	err = repo.Create(ctx, usage)
	require.NoError(t, err)

	usage.UsedInMB = 275.8
	err = repo.Update(ctx, usage)
	assert.NoError(t, err)

	results, err := repo.GetByDateRange(
		ctx,
		"user123",
		"5551234567",
		time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 11, 1, 23, 59, 59, 0, time.UTC),
	)
	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, 275.8, results[0].UsedInMB)
}
