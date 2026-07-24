// This file passes the database to the repositories, the repositories to the services, and the services to the handlers.
package app

import (
	"mkluxe-backend/internal/config"
	"mkluxe-backend/internal/handler"
	"mkluxe-backend/internal/repository"
	"mkluxe-backend/internal/routes"
	"mkluxe-backend/internal/service"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// BuildApp now accepts both the database and the configuration object
func BuildApp(db *mongo.Database, cfg *config.Config) *gin.Engine {
	// 1. Repositories
	userRepo := repository.NewUserRepository(db)
	catRepo := repository.NewCategoryRepository(db)
	prodRepo := repository.NewProductRepository(db)
	inqRepo := repository.NewInquiryRepository(db)
	statsRepo := repository.NewStatsRepository(db)

	// 2. Services
	authSvc := service.NewAuthService(userRepo, cfg)
	catSvc := service.NewCategoryService(catRepo)
	prodSvc := service.NewProductService(prodRepo, catRepo)
	inqSvc := service.NewInquiryService(inqRepo, prodRepo)
	statsSvc := service.NewStatsService(statsRepo, catRepo, inqRepo, prodRepo)
	r2Svc, err := service.NewR2Service(cfg)
	if err != nil {
		panic("failed to initialize R2 service: " + err.Error())
	}

	// 3. Handlers
	handlers := routes.AppHandlers{
		// 💡 Pass the loaded config into the AuthHandler to power the cookie logic
		Auth:     handler.NewAuthHandler(authSvc, cfg),
		Category: handler.NewCategoryHandler(catSvc),
		Product:  handler.NewProductHandler(prodSvc),
		Inquiry:  handler.NewInquiryHandler(inqSvc),
		Upload:   handler.NewUploadHandler(r2Svc),
		Stats:    handler.NewStatsHandler(statsSvc),
	}

	// 4. Router
	return routes.SetupRouter(handlers, cfg)
}
