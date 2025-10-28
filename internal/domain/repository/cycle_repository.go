package repository

import (
	"context"
	"time"

	"github.com/bowe99/phone-usage-service/internal/domain/model"
)

type CycleRepository interface {
	Create(ctx context.Context, cycle *model.Cycle) error
	GetByID(ctx context.Context, id string) (*model.Cycle, error)
	GetByMDN(ctx context.Context, mdn string) ([]*model.Cycle, error)
	GetByUserID(ctx context.Context, userID string) ([]*model.Cycle, error)
	GetCurrentCycle(ctx context.Context, userID, mdn string, currentDate time.Time) (*model.Cycle, error)
}