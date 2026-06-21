package dto

type CreateInquiryRequest struct {
	ProductID    string `json:"product_id" binding:"required"`
	CustomerName string `json:"customer_name" binding:"required,min=2"`
	Phone        string `json:"phone" binding:"required"`
	Message      string `json:"message" binding:"required"`
}

type UpdateInquiryStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=pending reviewed resolved"`
}
