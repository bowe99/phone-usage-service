package model

import "time"

type DailyUsage struct {
	ID        string    `bson:"_id,omitempty" json:"id"`
	MDN       string    `bson:"mdn" json:"mdn"`
	UserID    string    `bson:"userId" json:"userId"`
	UsageDate time.Time `bson:"usageDate" json:"usageDate"`
	UsedInMB  float64   `bson:"usedInMb" json:"usedInMb"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}

type DailyUsageResponse struct {
	Date  time.Time `json:"date"`
	Usage float64   `json:"dailyUsage"`
}

func (d *DailyUsage) ToResponse() *DailyUsageResponse {
	return &DailyUsageResponse{
		Date:  d.UsageDate,
		Usage: d.UsedInMB,
	}
}