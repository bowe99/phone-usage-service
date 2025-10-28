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

func TestCycleRepository_Create(t *testing.T) {
	ctx := context.Background()

	mongoContainer, err := mongodb.Run(ctx, "mongo:6")
	require.NoError(t, err)
	defer func() {
		if err := mongoContainer.Terminate(ctx); err != nil {
			t.Logf("Failed to terminate container: %v", err)
		}
	}()

	connStr, err := mongoContainer.ConnectionString(ctx)
	require.NoError(t, err)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connStr))
	require.NoError(t, err)
	defer client.Disconnect(ctx)

	db := client.Database("test_db")
	repo := repository.SetupCycleRepository(db)

	cycle := &model.Cycle{
		MDN:       "5551234567",
		StartDate: time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 11, 30, 23, 59, 59, 0, time.UTC),
		UserID:    "user123",
	}

	err = repo.Create(ctx, cycle)
	assert.NoError(t, err)
	assert.NotEmpty(t, cycle.ID)

	retrieved, err := repo.GetByID(ctx, cycle.ID)
	assert.NoError(t, err)
	assert.Equal(t, cycle.MDN, retrieved.MDN)
	assert.Equal(t, cycle.UserID, retrieved.UserID)
}

func TestCycleRepository_GetByMDN(t *testing.T) {
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
	repo := repository.SetupCycleRepository(db)

	cycle1 := &model.Cycle{
		MDN:       "5551234567",
		StartDate: time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 10, 31, 23, 59, 59, 0, time.UTC),
		UserID:    "user123",
	}
	cycle2 := &model.Cycle{
		MDN:       "5551234567",
		StartDate: time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 11, 30, 23, 59, 59, 0, time.UTC),
		UserID:    "user456",
	}

	err = repo.Create(ctx, cycle1)
	require.NoError(t, err)
	err = repo.Create(ctx, cycle2)
	require.NoError(t, err)

	cycles, err := repo.GetByMDN(ctx, "5551234567")
	assert.NoError(t, err)
	assert.Len(t, cycles, 2)

	assert.Equal(t, "user456", cycles[0].UserID) 
	assert.Equal(t, "user123", cycles[1].UserID)
}

func TestCycleRepository_GetCurrentCycle(t *testing.T) {
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
	repo := repository.SetupCycleRepository(db)

	pastCycle := &model.Cycle{
		MDN:       "5551234567",
		StartDate: time.Date(2024, 9, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 9, 30, 23, 59, 59, 0, time.UTC),
		UserID:    "user123",
	}
	currentCycle := &model.Cycle{
		MDN:       "5551234567",
		StartDate: time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 10, 31, 23, 59, 59, 0, time.UTC),
		UserID:    "user123",
	}
	futureCycle := &model.Cycle{
		MDN:       "5551234567",
		StartDate: time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 11, 30, 23, 59, 59, 0, time.UTC),
		UserID:    "user123",
	}

	repo.Create(ctx, pastCycle)
	repo.Create(ctx, currentCycle)
	repo.Create(ctx, futureCycle)

	testDate := time.Date(2024, 10, 15, 12, 0, 0, 0, time.UTC)
	cycle, err := repo.GetCurrentCycle(ctx, "user123", "5551234567", testDate)

	assert.NoError(t, err)
	assert.NotNil(t, cycle)
	assert.Equal(t, currentCycle.ID, cycle.ID)
	assert.Equal(t, time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC), cycle.StartDate)
}
