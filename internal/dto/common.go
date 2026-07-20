package dto

type PaginationRequest struct {
	Page  int    `form:"page" binding:"omitempty,min=1"`
	Limit int    `form:"limit" binding:"omitempty,min=1,max=100"`
	Sort  string `form:"sort"`
}

// FilterRequest is used for querying products with dynamic filters.
type FilterRequest struct {
	Status       string `json:"status"`
	CategorySlug string `json:"category_slug"`
	Search       string `json:"search"`
	IsFeatured   *bool  `json:"is_featured"`
	IsMostSold   *bool  `json:"is_most_sold"`
}
