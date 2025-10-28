package service

import (
	"context"
	"fmt"
	"time"

	"github.com/bowe99/phone-usage-service/internal/application/dtos"
	"github.com/bowe99/phone-usage-service/internal/domain/model"
	"github.com/bowe99/phone-usage-service/internal/domain/repository"
)

type DailyUsageService struct {
	usageRepo repository.DailyUsageRepository
	cycleRepo repository.CycleRepository
}

func SetupDailyUsageService(usageRepo repository.DailyUsageRepository, cycleRepo repository.CycleRepository) *DailyUsageService {
	return &DailyUsageService{
		usageRepo: usageRepo,
		cycleRepo: cycleRepo,
	}
}

// Algorithm:
// 1. Find the current active cycle for the user and MDN
// 2. Query usage records for the date range of that cycle
// 3. Return list of {date, daily usage}
func (s *DailyUsageService) GetCurrentCycleUsage(ctx context.Context, req dto.GetCurrentCycleUsageRequest) ([]*model.DailyUsageResponse, error) {
	if req.UserID == "" {
		return nil, fmt.Errorf("userId is required")
	}
	if req.MDN == "" {
		return nil, fmt.Errorf("mdn is required")
	}

	currentCycle, err := s.cycleRepo.GetCurrentCycle(ctx, req.UserID, req.MDN, time.Now())
	if err != nil {
		return nil, fmt.Errorf("no active billing cycle found for user %s and MDN %s", req.UserID, req.MDN)
	}

	usageRecords, err := s.usageRepo.GetByDateRange(
		ctx,
		req.UserID,
		req.MDN,
		currentCycle.StartDate,
		currentCycle.EndDate,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get usage records: %w", err)
	}

	responses := make([]*model.DailyUsageResponse, len(usageRecords))
	for i, record := range usageRecords {
		responses[i] = record.ToResponse()
	}

	return responses, nil
}
