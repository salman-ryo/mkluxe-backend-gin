package dto

import "time"

type RecentInquiryItem struct {
	ID           string    `json:"id"`
	ProductName  string    `json:"product_name"`
	CustomerName string    `json:"customer_name"`
	Phone        string    `json:"phone"`
	Message      string    `json:"message"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
}

type CategoryStatItem struct {
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	ProductCount int64  `json:"product_count"`
}

type DashboardStatsResponse struct {
	TotalProducts           int64               `json:"total_products"`
	TotalCategories         int64               `json:"total_categories"`
	TotalInquiries          int64               `json:"total_inquiries"`
	ProductStatusCounts     map[string]int64    `json:"product_status_counts"`
	InquiryStatusCounts     map[string]int64    `json:"inquiry_status_counts"`
	CategoryStats           []CategoryStatItem  `json:"category_stats"`
	TotalStock              int64               `json:"total_stock"`
	OutOfStockCount         int64               `json:"out_of_stock_count"`
	FeaturedProductsCount   int64               `json:"featured_products_count"`
	MostSoldProductsCount   int64               `json:"most_sold_products_count"`
	RecentInquiries         []RecentInquiryItem `json:"recent_inquiries"`
}
