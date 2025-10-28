package repository

import (
	"context"
	"time"

	"github.com/bowe99/phone-usage-service/internal/domain/model"
)

type DailyUsageRepository interface {
	Create(ctx context.Context, usage *model.DailyUsage) error
	GetByDateRange(ctx context.Context, userId, mdn string, startDate, endDate time.Time) ([]*model.DailyUsage, error)
	Update(ctx context.Context, usage *model.DailyUsage) error
}