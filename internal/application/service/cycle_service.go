package service

import (
	"context"
	"fmt"
	"time"

	dto "github.com/bowe99/phone-usage-service/internal/application/dtos"
	"github.com/bowe99/phone-usage-service/internal/domain/model"
	"github.com/bowe99/phone-usage-service/internal/domain/repository"
)

type CycleService struct {
	cycleRepo repository.CycleRepository
}

func SetupCycleService(cycleRepo repository.CycleRepository) *CycleService {
	return &CycleService{
		cycleRepo: cycleRepo,
	}
}

// Note: Query by MDN, not just userId, because MDNs can be transferred between users
// This ensures we return the full history of the phone number, regardless of ownership changes
func (s *CycleService) GetCycleHistory(ctx context.Context, req dto.GetCycleHistoryRequest) ([]*model.CycleResponse, error) {
	if req.UserID == "" {
		return nil, fmt.Errorf("userId is required")
	}
	if req.MDN == "" {
		return nil, fmt.Errorf("mdn is required")
	}

	cycles, err := s.cycleRepo.GetByMDN(ctx, req.MDN)
	if err != nil {
		return nil, fmt.Errorf("failed to get cycle history: %w", err)
	}

	responses := make([]*model.CycleResponse, len(cycles))
	for i, cycle := range cycles {
		responses[i] = cycle.ToResponse()
	}

	return responses, nil
}

func (s *CycleService) GetCurrentCycle(ctx context.Context, userID, mdn string) (*model.Cycle, error) {
	cycle, err := s.cycleRepo.GetCurrentCycle(ctx, userID, mdn, time.Now())
	if err != nil {
		return nil, err
	}

	return cycle, nil
}