package model

import "time"

type Cycle struct {
	ID        string    `bson:"_id,omitempty" json:"id"`
	MDN       string    `bson:"mdn" json:"mdn"`      
	StartDate time.Time `bson:"startDate" json:"startDate"`
	EndDate   time.Time `bson:"endDate" json:"endDate"`
	UserID    string    `bson:"userId" json:"userId"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
}

type CycleResponse struct {
	CycleID   string    `json:"cycleId"`
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
}

func (c *Cycle) ToResponse() *CycleResponse {
	return &CycleResponse{
		CycleID:   c.ID,
		StartDate: c.StartDate,
		EndDate:   c.EndDate,
	}
}