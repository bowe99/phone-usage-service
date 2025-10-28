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

type MockDailyUsageRepository struct {
	mock.Mock
}

func (m *MockDailyUsageRepository) Create(ctx context.Context, usage *model.DailyUsage) error {
	args := m.Called(ctx, usage)
	return args.Error(0)
}

func (m *MockDailyUsageRepository) GetByDateRange(ctx context.Context, userID, mdn string, startDate, endDate time.Time) ([]*model.DailyUsage, error) {
	args := m.Called(ctx, userID, mdn, startDate, endDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.DailyUsage), args.Error(1)
}

func (m *MockDailyUsageRepository) Update(ctx context.Context, usage *model.DailyUsage) error {
	args := m.Called(ctx, usage)
	return args.Error(0)
}

func TestDailyUsageService_GetCurrentCycleUsage(t *testing.T) {
	// Arrange
	mockUsageRepo := new(MockDailyUsageRepository)
	mockCycleRepo := new(MockCycleRepository)
	usageService := service.SetupDailyUsageService(mockUsageRepo, mockCycleRepo)

	req := dto.GetCurrentCycleUsageRequest{
		UserID: "user123",
		MDN:    "5551234567",
	}

	currentCycle := &model.Cycle{
		ID:        "cycle1",
		MDN:       "5551234567",
		StartDate: time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2024, 11, 30, 23, 59, 59, 0, time.UTC),
		UserID:    "user123",
	}

	expectedUsage := []*model.DailyUsage{
		{
			ID:        "usage1",
			MDN:       "5551234567",
			UserID:    "user123",
			UsageDate: time.Date(2024, 11, 1, 0, 0, 0, 0, time.UTC),
			UsedInMB:  250.5,
		},
		{
			ID:        "usage2",
			MDN:       "5551234567",
			UserID:    "user123",
			UsageDate: time.Date(2024, 11, 2, 0, 0, 0, 0, time.UTC),
			UsedInMB:  180.3,
		},
	}

	// Mock: GetCurrentCycle returns active cycle
	mockCycleRepo.On("GetCurrentCycle", mock.Anything, req.UserID, req.MDN, mock.AnythingOfType("time.Time")).
		Return(currentCycle, nil)

	// Mock: GetByDateRange returns usage records
	mockUsageRepo.On("GetByDateRange", mock.Anything, req.UserID, req.MDN, currentCycle.StartDate, currentCycle.EndDate).
		Return(expectedUsage, nil)

	// Act
	result, err := usageService.GetCurrentCycleUsage(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, 250.5, result[0].Usage)
	assert.Equal(t, 180.3, result[1].Usage)
	mockCycleRepo.AssertExpectations(t)
	mockUsageRepo.AssertExpectations(t)
}

func TestDailyUsageService_GetCurrentCycleUsage_NoCycleFound(t *testing.T) {
	// Arrange
	mockUsageRepo := new(MockDailyUsageRepository)
	mockCycleRepo := new(MockCycleRepository)
	usageService := service.SetupDailyUsageService(mockUsageRepo, mockCycleRepo)

	req := dto.GetCurrentCycleUsageRequest{
		UserID: "user123",
		MDN:    "5551234567",
	}

	// Mock: No current cycle found
	mockCycleRepo.On("GetCurrentCycle", mock.Anything, req.UserID, req.MDN, mock.AnythingOfType("time.Time")).
		Return(nil, assert.AnError)

	// Act
	result, err := usageService.GetCurrentCycleUsage(context.Background(), req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "no active billing cycle found")
	mockCycleRepo.AssertExpectations(t)
}

func TestDailyUsageService_GetCurrentCycleUsage_InvalidInput(t *testing.T) {
	mockUsageRepo := new(MockDailyUsageRepository)
	mockCycleRepo := new(MockCycleRepository)
	usageService := service.SetupDailyUsageService(mockUsageRepo, mockCycleRepo)

	// Test missing userId
	req := dto.GetCurrentCycleUsageRequest{
		UserID: "",
		MDN:    "5551234567",
	}

	result, err := usageService.GetCurrentCycleUsage(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "userId is required")
}
