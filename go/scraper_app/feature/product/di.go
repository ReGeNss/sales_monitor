package product

import (
	"sales_monitor/scraper_app/feature/product/data/repository"
	domainService "sales_monitor/scraper_app/feature/product/domain/service"
	"sales_monitor/scraper_app/feature/product/service"
	"sales_monitor/scraper_app/feature/product/service/usecase"
	"sales_monitor/scraper_app/utils"

	"gorm.io/gorm"
)

func NewProductService(db *gorm.DB, eventBus utils.EventBus) service.ProductService {
	brandRepository := repository.NewBrandRepository(db)
	categoryRepository := repository.NewCategoryRepository(db)
	marketplaceRepository := repository.NewMarketplaceRepository(db)
	priceRepository := repository.NewPriceRepository(db)
	productRepository := repository.NewProductRepository(db)

	matcher := domainService.NewProductMatcher()

	return service.NewProductService(
		usecase.NewResolveCategoryUseCase(categoryRepository),
		usecase.NewResolveMarketplaceUseCase(marketplaceRepository),
		usecase.NewAssignBrandsUseCase(brandRepository),
		usecase.NewResolveBrandUseCase(brandRepository),
		usecase.NewResolveProductUseCase(productRepository, matcher),
		usecase.NewRecordPriceUseCase(marketplaceRepository, priceRepository),
		marketplaceRepository,
		eventBus,
	)
}
