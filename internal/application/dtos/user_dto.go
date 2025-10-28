package dto

type CreateUserRequest struct {
	FirstName string `json:"firstName" binding:"required,min=2,max=50"`
	LastName  string `json:"lastName" binding:"required,min=2,max=50"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
}

type UpdateUserRequest struct {
	FirstName string `json:"firstName" binding:"omitempty,min=2,max=50"`
	LastName  string `json:"lastName" binding:"omitempty,min=2,max=50"`
	Email     string `json:"email" binding:"omitempty,email"`
	Password  string `json:"password" binding:"omitempty,min=8"`
}