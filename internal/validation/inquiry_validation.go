package validation

import (
	"errors"
	"mkluxe-backend/internal/dto"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ValidateInquiry(req *dto.CreateInquiryRequest) error {
	if _, err := primitive.ObjectIDFromHex(req.ProductID); err != nil {
		return errors.New("invalid product_id format")
	}
	if len(req.Phone) < 7 {
		return errors.New("phone number provided is too short")
	}
	return nil
}
