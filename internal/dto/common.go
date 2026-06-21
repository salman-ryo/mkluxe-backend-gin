package dto

type PaginationRequest struct {
	Page  int    `form:"page" binding:"omitempty,min=1"`
	Limit int    `form:"limit" binding:"omitempty,min=1,max=100"`
	Sort  string `form:"sort"`
}

type FilterRequest struct {
	Search     string `form:"search"`
	CategoryID string `form:"category_id"`
	Status     string `form:"status"`
}
