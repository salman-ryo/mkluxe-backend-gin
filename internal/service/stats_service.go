package service

import (
	"context"
	"log"
	"mkluxe-backend/internal/dto"
	"mkluxe-backend/internal/repository"
)

type StatsService struct {
	statsRepo *repository.StatsRepository
	catRepo   *repository.CategoryRepository
	inqRepo   *repository.InquiryRepository
	prodRepo  *repository.ProductRepository
}

func NewStatsService(
	statsRepo *repository.StatsRepository,
	catRepo *repository.CategoryRepository,
	inqRepo *repository.InquiryRepository,
	prodRepo *repository.ProductRepository,
) *StatsService {
	return &StatsService{
		statsRepo: statsRepo,
		catRepo:   catRepo,
		inqRepo:   inqRepo,
		prodRepo:  prodRepo,
	}
}

func (s *StatsService) GetDashboardStats(ctx context.Context) (*dto.DashboardStatsResponse, error) {
	// 1. Fetch total counts
	totalProds, err := s.statsRepo.CountProducts(ctx)
	if err != nil {
		log.Printf("Error counting products: %v", err)
	}

	totalCats, err := s.statsRepo.CountCategories(ctx)
	if err != nil {
		log.Printf("Error counting categories: %v", err)
	}

	totalInqs, err := s.statsRepo.CountInquiries(ctx)
	if err != nil {
		log.Printf("Error counting inquiries: %v", err)
	}

	// 2. Fetch counts grouped by status
	prodStatusCounts, err := s.statsRepo.GetProductStatusCounts(ctx)
	if err != nil {
		log.Printf("Error getting product status counts: %v", err)
		prodStatusCounts = make(map[string]int64)
	}

	inqStatusCounts, err := s.statsRepo.GetInquiryStatusCounts(ctx)
	if err != nil {
		log.Printf("Error getting inquiry status counts: %v", err)
		inqStatusCounts = make(map[string]int64)
	}

	// 3. Featured & Most Sold Counts
	featuredProdsCount, err := s.statsRepo.CountFeaturedProducts(ctx)
	if err != nil {
		log.Printf("Error counting featured products: %v", err)
	}

	mostSoldProdsCount, err := s.statsRepo.CountMostSoldProducts(ctx)
	if err != nil {
		log.Printf("Error counting most sold products: %v", err)
	}

	// 4. Stock Stats
	totalStock, outOfStockCount, err := s.statsRepo.GetStockStats(ctx)
	if err != nil {
		log.Printf("Error getting stock stats: %v", err)
	}

	// 5. Category stats (with counts)
	categories, err := s.catRepo.ListAll(ctx, nil)
	var categoryStats []dto.CategoryStatItem
	if err == nil {
		catCounts, err := s.statsRepo.GetProductCategoryCounts(ctx)
		if err != nil {
			log.Printf("Error getting product category counts: %v", err)
			catCounts = make(map[string]int64)
		}

		categoryStats = make([]dto.CategoryStatItem, 0, len(categories))
		for _, cat := range categories {
			count := catCounts[cat.Slug] // defaults to 0 if not present
			categoryStats = append(categoryStats, dto.CategoryStatItem{
				Name:         cat.Name,
				Slug:         cat.Slug,
				ProductCount: count,
			})
		}
	} else {
		log.Printf("Error listing categories: %v", err)
		categoryStats = []dto.CategoryStatItem{}
	}

	// 6. Recent Inquiries (top 5) with product names
	inquiries, _, err := s.inqRepo.List(ctx, 1, 5)
	var recentInquiries []dto.RecentInquiryItem
	if err == nil {
		recentInquiries = make([]dto.RecentInquiryItem, 0, len(inquiries))
		for _, inq := range inquiries {
			prodName := "Unknown Product"
			prod, err := s.prodRepo.GetByID(ctx, inq.ProductID)
			if err == nil && prod != nil {
				prodName = prod.Name
			}

			recentInquiries = append(recentInquiries, dto.RecentInquiryItem{
				ID:           inq.ID.Hex(),
				ProductName:  prodName,
				CustomerName: inq.CustomerName,
				Phone:        inq.Phone,
				Message:      inq.Message,
				Status:       inq.Status,
				CreatedAt:    inq.CreatedAt,
			})
		}
	} else {
		log.Printf("Error listing recent inquiries: %v", err)
		recentInquiries = []dto.RecentInquiryItem{}
	}

	return &dto.DashboardStatsResponse{
		TotalProducts:         totalProds,
		TotalCategories:       totalCats,
		TotalInquiries:        totalInqs,
		ProductStatusCounts:   prodStatusCounts,
		InquiryStatusCounts:   inqStatusCounts,
		CategoryStats:         categoryStats,
		TotalStock:            totalStock,
		OutOfStockCount:       outOfStockCount,
		FeaturedProductsCount: featuredProdsCount,
		MostSoldProductsCount: mostSoldProdsCount,
		RecentInquiries:       recentInquiries,
	}, nil
}
