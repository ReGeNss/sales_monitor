package service

import (
	"log"
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

		var marketplaceID int
		marketplace, err := s.productRepository.GetMarketplaceByName(data.MarketplaceName)
		if err != nil {
			marketplace = &models.Marketplace{
				Name: data.MarketplaceName,
			}
			s.productRepository.CreateMarketplace(marketplace)
			marketplaceID = marketplace.MarketplaceID
		} else {
			marketplaceID = marketplace.MarketplaceID
		}

		var categoryID int
		category, err := s.productRepository.GetCategoryByName(data.Category)
		if err != nil {
			category = &models.Category{
				Name: data.Category,
			}
			s.productRepository.CreateCategory(category)
			categoryID = category.CategoryID
		} else {
			categoryID = category.CategoryID
		}

		brandProducts := groupProductsByBrand(data.Products)

		for brandName, products := range brandProducts {
			var brandID int
			brand, err := s.productRepository.GetBrandByName(brandName)
			if err != nil {
				id, err := s.productRepository.CreateBrand(&models.Brand{
					Name: brandName,
				})
				if err != nil {
					log.Printf("could not create brand: %v", err)
					continue
				}
				brandID = int(id)
			} else {
				brandID = brand.BrandID
			}

			for _, product := range products {
				fingerprint := utils.NormalizeProductName(product.Name)

				var productID int

				existingProduct, err := s.productRepository.GetProductByFingerprint(fingerprint)

				if err != nil {
					matchedProductID, err := s.productRepository.GetMostSimilarProductID(fingerprint)
					if err != nil {
						s.productRepository.CreateProduct(&models.Product{
							Name:            product.Name,
							NameFingerprint: fingerprint,
							ImageURL:        product.Image,
							BrandID:         brandID,
							CategoryID:      categoryID,
						})
					}
					productID = int(matchedProductID)
				} else {
					productID = existingProduct.ProductID
				}

				s.productRepository.AddPriceToProduct(&models.Price{
					ProductID:     productID,
					MarketplaceID: marketplaceID,
					RegularPrice:  product.RegularPrice,
					DiscountPrice: &product.DiscountedPrice,
					URL:           product.URL,
				})
			}
		}
	}
}

func groupProductsByBrand(products []*entity.ScrapedProduct) map[string][]*entity.ScrapedProduct {
	brandProducts := make(map[string][]*entity.ScrapedProduct)
	for _, product := range products {
		brandProducts[product.BrandName] = append(brandProducts[product.BrandName], product)
	}
	return brandProducts
}
