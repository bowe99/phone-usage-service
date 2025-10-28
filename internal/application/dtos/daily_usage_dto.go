package dto

type GetCurrentCycleUsageRequest struct {
	UserID string `json:"userId" binding:"required"`
	MDN    string `json:"mdn" binding:"required,len=10"`
}
