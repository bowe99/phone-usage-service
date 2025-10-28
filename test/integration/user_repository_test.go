package integration

import (
	"context"
	"testing"

	"github.com/bowe99/phone-usage-service/internal/domain/model"
	"github.com/bowe99/phone-usage-service/internal/infra/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestUserRepository_Create(t *testing.T) {
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
	repo := repository.SetupUserRepository(db)

	// Test Create
	user := &model.User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Password:  "hashedpassword",
	}

	err = repo.Create(ctx, user)
	assert.NoError(t, err)
	assert.NotEmpty(t, user.ID)

	// Test GetByID
	retrieved, err := repo.GetByID(ctx, user.ID)
	assert.NoError(t, err)
	assert.Equal(t, user.FirstName, retrieved.FirstName)
	assert.Equal(t, user.Email, retrieved.Email)
}

func TestUserRepository_Update(t *testing.T) {
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
	repo := repository.SetupUserRepository(db)

	// Create user first
	user := &model.User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Password:  "hashedpassword",
	}
	err = repo.Create(ctx, user)
	require.NoError(t, err)

	// Update user
	user.FirstName = "Jane"
	user.Email = "jane.doe@example.com"
	err = repo.Update(ctx, user)
	assert.NoError(t, err)

	// Verify update
	updated, err := repo.GetByID(ctx, user.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Jane", updated.FirstName)
	assert.Equal(t, "jane.doe@example.com", updated.Email)
}
