package service

import (
	"sales_monitor/internal/models"
	"sales_monitor/scraper_app/shared/product/domain/entity"
	"sales_monitor/scraper_app/shared/product/domain/repository"
	"sales_monitor/scraper_app/shared/product/utils"
)

type ProductService interface {
	ProcessProducts(scrapedData []*entity.ScrapedProducts)
}

type productServiceImpl struct {
	productRepository repository.ProductRepository
}

func NewProductService(productRepository repository.ProductRepository) ProductService {
	return &productServiceImpl{
		productRepository: productRepository,
	}
}

func (s *productServiceImpl) ProcessProducts(scrapedData []*entity.ScrapedProducts) {
	for _, data := range scrapedData {
		if len(data.Products) == 0 {
			continue
		}

		var marketplaceID uint
		marketplace, err := s.productRepository.GetMarketplaceByName(data.MarketplaceName)
		if err != nil {
			marketplace = &models.Marketplace{
				Name: data.MarketplaceName,
			}
			s.productRepository.CreateMarketplace(marketplace)
			marketplaceID = uint(marketplace.MarketplaceID)
		} else {
			marketplaceID = uint(marketplace.MarketplaceID)
		}

		for _, product := range data.Products {
			fingerprint := utils.NormalizeProductName(product.Name)

			var productID uint

			existingProduct, err := s.productRepository.GetProductByFingerprint(fingerprint)

			if err != nil {
				matchedProductID, err := s.productRepository.GetMostSimilarProductID(fingerprint)
				if err != nil {
					s.productRepository.CreateProduct(&models.Product{
						Name:            product.Name,
						NameFingerprint: fingerprint,
						ImageURL:        product.Image,
						BrandID:         1,
						CategoryID:      1,
					})
				}
				productID = matchedProductID
			} else {
				productID = uint(existingProduct.ProductID)
			}

			s.productRepository.AddPriceToProduct(&models.Price{
				ProductID:     int(productID),
				MarketplaceID: int(marketplaceID),
				RegularPrice:  product.RegularPrice,
				DiscountPrice: &product.DiscountedPrice,
				URL:           product.Image,
			})
		}
	}
}
