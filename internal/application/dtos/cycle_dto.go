package dto

type GetCycleHistoryRequest struct {
	UserID string `json:"userId" binding:"required"`
	MDN    string `json:"mdn" binding:"required,len=10"` // US phone numbers are 10 digits
}