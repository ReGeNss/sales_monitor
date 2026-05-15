package service

import (
	"log"
	"sales_monitor/scraper_app/feature/product/domain/entity"
	"sales_monitor/scraper_app/feature/product/domain/repository"
	"sales_monitor/scraper_app/feature/product/service/usecase"
	scraper "sales_monitor/scraper_app/feature/scraper/domain/entity"
	"sales_monitor/scraper_app/utils"
)

type ProductService interface {
	ProcessProducts(map[string]*scraper.ScrapingResult)
}

type productServiceImpl struct {
	resolveCategory       usecase.ResolveCategoryUseCase
	resolveMarketplace    usecase.ResolveMarketplaceUseCase
	assignBrands          usecase.AssignBrandsUseCase
	resolveBrand          usecase.ResolveBrandUseCase
	resolveProduct        usecase.ResolveProductUseCase
	recordPrice           usecase.RecordPriceUseCase
	marketplaceRepository repository.MarketplaceRepository
	eventBus              utils.EventBus
}

func NewProductService(
	resolveCategory usecase.ResolveCategoryUseCase,
	resolveMarketplace usecase.ResolveMarketplaceUseCase,
	assignBrands usecase.AssignBrandsUseCase,
	resolveBrand usecase.ResolveBrandUseCase,
	resolveProduct usecase.ResolveProductUseCase,
	recordPrice usecase.RecordPriceUseCase,
	marketplaceRepository repository.MarketplaceRepository,
	eventBus utils.EventBus,
) ProductService {
	return &productServiceImpl{
		resolveCategory:       resolveCategory,
		resolveMarketplace:    resolveMarketplace,
		assignBrands:          assignBrands,
		resolveBrand:          resolveBrand,
		resolveProduct:        resolveProduct,
		recordPrice:           recordPrice,
		marketplaceRepository: marketplaceRepository,
		eventBus:              eventBus,
	}
}

func (s *productServiceImpl) ProcessProducts(scrapedData map[string]*scraper.ScrapingResult) {
	for categoryName, result := range scrapedData {
		categoryID, err := s.resolveCategory.Execute(categoryName)
		if err != nil {
			log.Printf("could not resolve category %q: %v", categoryName, err)
			continue
		}

		for _, data := range result.ScrapedProducts {
			if len(data.Products) == 0 {
				continue
			}

			marketplaceID, err := s.resolveMarketplace.Execute(data.MarketplaceName)
			if err != nil {
				log.Printf("could not resolve marketplace %q: %v", data.MarketplaceName, err)
				continue
			}

			brandGroups, err := s.assignBrands.Execute(data.Products)
			if err != nil {
				log.Printf("could not assign brands: %v", err)
				continue
			}

			s.processBrandGroups(brandGroups, categoryName, categoryID, marketplaceID, result.ProductDifferentiationEntity)
		}
	}
}

func (s *productServiceImpl) processBrandGroups(
	groups map[string][]*entity.ScrapedProduct,
	categoryName string,
	categoryID int,
	marketplaceID int,
	differentiation *entity.ProductDifferentiationEntity,
) {
	for brandName, products := range groups {
		brandID, err := s.resolveBrand.Execute(brandName)
		if err != nil {
			log.Printf("could not resolve brand %q: %v", brandName, err)
			continue
		}

		knownURLs, err := s.marketplaceRepository.GetLaterScrapedProducts(brandID)
		if err != nil {
			log.Printf("could not get later scraped products: %v", err)
		}

		priceDrops := s.processProducts(products, categoryName, brandID, categoryID, marketplaceID, knownURLs, differentiation)

		if len(priceDrops) > 0 {
			s.eventBus.Publish(&entity.NotificationTask{
				BrandID:   brandID,
				BrandName: brandName,
				Products:  priceDrops,
			})
		}
	}
}

func (s *productServiceImpl) processProducts(
	products []*entity.ScrapedProduct,
	categoryName string,
	brandID int,
	categoryID int,
	marketplaceID int,
	knownURLs entity.LaterScrapedProductsUrls,
	differentiation *entity.ProductDifferentiationEntity,
) []*entity.Product {
	priceDrops := []*entity.Product{}

	for _, product := range products {
		productID, existing, err := s.resolveProduct.Execute(usecase.ResolveProductInput{
			Scraped:         product,
			Fingerprint:     product.GetFingerprint([]string{categoryName}),
			Attributes:      buildAttributes(product),
			BrandID:         brandID,
			CategoryID:      categoryID,
			Differentiation: differentiation,
		})
		if err != nil {
			log.Printf("could not resolve product: %v", err)
			continue
		}

		drop, err := s.recordPrice.Execute(usecase.RecordPriceInput{
			Scraped:              product,
			ProductID:            productID,
			MarketplaceID:        marketplaceID,
			ExistingProduct:      existing,
			KnownMarketplaceURLs: knownURLs,
		})
		if err != nil {
			log.Printf("could not record price: %v", err)
			continue
		}
		if drop != nil {
			priceDrops = append(priceDrops, drop)
		}
	}

	return priceDrops
}

func buildAttributes(product *entity.ScrapedProduct) []*entity.ProductAttribute {
	attributes := []*entity.ProductAttribute{}
	volume := product.Volume()
	if volume != "" {
		attributes = append(attributes, &entity.ProductAttribute{
			Type:  entity.AttributeTypeVolume,
			Value: volume,
		})
	}
	if product.Weight() != "" {
		attributes = append(attributes, &entity.ProductAttribute{
			Type:  entity.AttributeTypeWeight,
			Value: product.Weight(),
		})
	}
	return attributes
}
