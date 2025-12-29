package service

import (
	"sales_monitor/internal/models"
	"sales_monitor/scraper_app/shared/product/domain/entity"
	"sales_monitor/scraper_app/shared/product/domain/repository"
	"sales_monitor/scraper_app/shared/product/utils"
)

type ProductService interface {
	ProcessProducts(products []*entity.ScrapedProduct)
}

type productServiceImpl struct {
	productRepository repository.ProductRepository
}

func NewProductService(productRepository repository.ProductRepository) ProductService {
	return &productServiceImpl{
		productRepository: productRepository,
	}
}

func (s *productServiceImpl) ProcessProducts(products []*entity.ScrapedProduct) {
	for _, product := range products {
		fingerprint := utils.NormalizeProductName(product.Name)

		var productID uint

		existingProduct, err := s.productRepository.GetProductByFingerprint(fingerprint)

		if err != nil {
			matchedProductID, err := s.productRepository.GetMostSimilarProductID(fingerprint)
			if err != nil {
				s.productRepository.CreateProduct(&models.Product{
					Name: product.Name,
					NameFingerprint: fingerprint,
					ImageURL: product.Image,
					BrandID: 1,
					CategoryID: 1,
				})
			}
			productID = matchedProductID
		} else {
			productID = uint(existingProduct.ProductID)
		}

		s.productRepository.AddPriceToProduct(&models.Price{
			ProductID: int(productID),
			MarketplaceID: 1,
			RegularPrice: product.RegularPrice,
			DiscountPrice: &product.DiscountedPrice,
			URL: product.Image,
		})
	}
}
