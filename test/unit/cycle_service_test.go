package unit

import (
	"context"
	"testing"
	"time"

	"github.com/bowe99/phone-usage-service/internal/application/dtos"
	"github.com/bowe99/phone-usage-service/internal/application/service"
	"github.com/bowe99/phone-usage-service/internal/domain/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCycleRepository struct {
	mock.Mock
}

func (m *MockCycleRepository) Create(ctx context.Context, cycle *model.Cycle) error {
	args := m.Called(ctx, cycle)
	return args.Error(0)
}

func (m *MockCycleRepository) GetByID(ctx context.Context, id string) (*model.Cycle, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Cycle), args.Error(1)
}

func (m *MockCycleRepository) GetByMDN(ctx context.Context, mdn string) ([]*model.Cycle, error) {
	args := m.Called(ctx, mdn)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Cycle), args.Error(1)
}

func (m *MockCycleRepository) GetByUserID(ctx context.Context, userID string) ([]*model.Cycle, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Cycle), args.Error(1)
}

func (m *MockCycleRepository) GetCurrentCycle(ctx context.Context, userID, mdn string, currentDate time.Time) (*model.Cycle, error) {
	args := m.Called(ctx, userID, mdn, currentDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Cycle), args.Error(1)
}

func TestCycleService_GetCycleHistory(t *testing.T) {
	mockRepo := new(MockCycleRepository)
	cycleService := service.SetupCycleService(mockRepo)

	req := dto.GetCycleHistoryRequest{
		UserID: "user123",
		MDN:    "5551234567",
	}

	expectedCycles := []*model.Cycle{
		{
			ID:        "cycle1",
			MDN:       "5551234567",
			StartDate: time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 11, 30, 23, 59, 59, 0, time.UTC),
			UserID:    "user123",
		},
		{
			ID:        "cycle2",
			MDN:       "5551234567",
			StartDate: time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 10, 31, 23, 59, 59, 0, time.UTC),
			UserID:    "user123",
		},
	}

	mockRepo.On("GetByMDN", mock.Anything, req.MDN).Return(expectedCycles, nil)

	result, err := cycleService.GetCycleHistory(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, "cycle1", result[0].CycleID)
	assert.Equal(t, "cycle2", result[1].CycleID)
	mockRepo.AssertExpectations(t)
}

func TestCycleService_GetCycleHistory_InvalidInput(t *testing.T) {
	mockRepo := new(MockCycleRepository)
	cycleService := service.SetupCycleService(mockRepo)

	req := dto.GetCycleHistoryRequest{
		UserID: "",
		MDN:    "5551234567",
	}

	result, err := cycleService.GetCycleHistory(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "userId is required")
}
