package service

import (
	"context"
	"errors"

	"mkluxe-backend/internal/constants"
	"mkluxe-backend/internal/domain"
	"mkluxe-backend/internal/dto"
	"mkluxe-backend/internal/repository"
	"mkluxe-backend/internal/validation"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InquiryService struct {
	inqRepo     *repository.InquiryRepository
	productRepo *repository.ProductRepository
}

func NewInquiryService(iRepo *repository.InquiryRepository, pRepo *repository.ProductRepository) *InquiryService {
	return &InquiryService{inqRepo: iRepo, productRepo: pRepo}
}

func (s *InquiryService) CreateInquiry(ctx context.Context, req *dto.CreateInquiryRequest) (*domain.Inquiry, error) {
	if err := validation.ValidateInquiry(req); err != nil {
		return nil, err
	}

	productID, _ := primitive.ObjectIDFromHex(req.ProductID)
	_, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, errors.New("referenced product not found")
	}

	inq := &domain.Inquiry{
		ProductID:    productID,
		CustomerName: req.CustomerName,
		Phone:        req.Phone,
		Message:      req.Message,
		Status:       constants.InquiryStatusPending,
	}

	if err := s.inqRepo.Create(ctx, inq); err != nil {
		return nil, err
	}
	return inq, nil
}

func (s *InquiryService) UpdateStatus(ctx context.Context, inquiryIDStr, status string) error {
	id, err := primitive.ObjectIDFromHex(inquiryIDStr)
	if err != nil {
		return errors.New("invalid inquiry ID")
	}
	return s.inqRepo.UpdateStatus(ctx, id, status)
}