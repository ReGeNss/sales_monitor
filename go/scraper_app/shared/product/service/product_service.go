package service

import (
	"log"
	scraper "sales_monitor/scraper_app/feature/scraper/domain/entity"
	"sales_monitor/scraper_app/shared/product/domain/entity"
	"sales_monitor/scraper_app/shared/product/domain/gateway"
	"sales_monitor/scraper_app/shared/product/domain/repository"
	domainservice "sales_monitor/scraper_app/shared/product/domain/service"
	"strings"
)

type ProductService interface {
	ProcessProducts(map[string]*scraper.ScrapingResult)
}

type productServiceImpl struct {
	productRepository     repository.ProductRepository
	categoryRepository    repository.CategoryRepository
	brandRepository       repository.BrandRepository
	marketplaceRepository repository.MarketplaceRepository
	priceRepository       repository.PriceRepository
	notificationPublisher gateway.NotificationPublisher
	productMatcher        domainservice.ProductMatcher
}

func NewProductService(
	productRepository repository.ProductRepository,
	categoryRepository repository.CategoryRepository,
	brandRepository repository.BrandRepository,
	marketplaceRepository repository.MarketplaceRepository,
	priceRepository repository.PriceRepository,
	notificationPublisher gateway.NotificationPublisher,
	productMatcher domainservice.ProductMatcher,
) ProductService {
	return &productServiceImpl{
		productRepository:     productRepository,
		categoryRepository:    categoryRepository,
		brandRepository:       brandRepository,
		marketplaceRepository: marketplaceRepository,
		priceRepository:       priceRepository,
		notificationPublisher: notificationPublisher,
		productMatcher:        productMatcher,
	}
}

func (s *productServiceImpl) ProcessProducts(scrapedData map[string]*scraper.ScrapingResult) {
	for categoryName, scrapedData := range scrapedData {

		var categoryID int
		category, err := s.categoryRepository.GetCategoryByName(categoryName)
		if err != nil {
			category = &entity.Category{
				Name: categoryName,
			}
			s.categoryRepository.CreateCategory(category)
			categoryID = category.ID
		} else {
			categoryID = category.ID
		}

		for _, data := range scrapedData.ScrapedProducts {
			if len(data.Products) == 0 {
				continue
			}

			var marketplaceID int
			marketplace, err := s.marketplaceRepository.GetMarketplaceByName(data.MarketplaceName)
			if err != nil {
				marketplace = &entity.Marketplace{
					Name: data.MarketplaceName,
				}
				s.marketplaceRepository.CreateMarketplace(marketplace)
				marketplaceID = marketplace.ID
			} else {
				marketplaceID = marketplace.ID
			}

			brandProducts, unknownBrandProducts := groupProductsByBrand(data.Products)

			if len(unknownBrandProducts) > 0 {
				allBrands, err := s.brandRepository.GetAllBrands()
				if err != nil {
					log.Printf("could not get all brands: %v", err)
					continue
				}
				brandProducts = getBrandsFromProductName(unknownBrandProducts, brandProducts, allBrands)

			}

			for brandName, products := range brandProducts {
				var brandID int
				brand, err := s.brandRepository.GetBrandByName(brandName)
				if err != nil {
					id, err := s.brandRepository.CreateBrand(&entity.Brand{
						Name: brandName,
					})
					if err != nil {
						log.Printf("could not create brand: %v", err)
						continue
					}
					brandID = int(id)
				} else {
					brandID = brand.ID
				}

				laterScrapedProductsUrls, err := s.marketplaceRepository.GetLaterScrapedProducts(brandID)
				if err != nil {
					log.Printf("could not get later scraped products: %v", err)
				}

				notificationProducts := []*entity.Product{}

				for _, product := range products {
					if id, ok := laterScrapedProductsUrls[product.URL]; ok {
						s.marketplaceRepository.AddPriceToMarketplaceProductID(
							id,
							product.RegularPrice,
							&product.SpecialPrice,
						)
					}

					fingerprint := product.GetFingerprint([]string{categoryName})
					attributes := []*entity.ProductAttribute{}

					if product.Volume != "" {
						attributes = append(attributes, &entity.ProductAttribute{
							Type:  entity.AttributeTypeVolume,
							Value: product.Volume,
						})
					}
					if product.Weight != "" {
						attributes = append(attributes, &entity.ProductAttribute{
							Type:  entity.AttributeTypeWeight,
							Value: product.Weight,
						})
					}

					var productID int
					var foundProduct *entity.Product

					existingProduct, err := s.productRepository.GetProductByFingerprint(fingerprint, brandID, categoryID, attributes)

					if err != nil {
						var matchedProduct *entity.Product
						candidates, candErr := s.productRepository.FindSimilarCandidates(fingerprint, attributes, brandID, categoryID)
						if candErr == nil && fingerprint != nil {
							if m, ok := s.productMatcher.PickBestMatch(*fingerprint, candidates, scrapedData.ProductDifferentiationEntity); ok {
								matchedProduct = m
							}
						} else if candErr == nil && fingerprint == nil && len(candidates) > 0 {
							matchedProduct = candidates[0]
						}

						if matchedProduct == nil {
							id, err := s.productRepository.CreateProduct(&entity.Product{
								Name:            product.Name,
								NameFingerprint: fingerprint,
								ImageURL:        product.Image,
								BrandID:         brandID,
								CategoryID:      categoryID,
							}, attributes)

							if err != nil {
								log.Printf("could not create product: %v", err)
								continue
							}
							productID = int(id)
						} else {
							foundProduct = matchedProduct
							productID = matchedProduct.ID
						}
					} else {
						foundProduct = existingProduct
						productID = existingProduct.ID
					}

					if foundProduct != nil {
						laterPrice, err := s.priceRepository.GetLatestProductPrice(productID)
						if err != nil {
							log.Printf("could not get latest product price: %v", err)
						} else if *laterPrice.SpecialPrice < product.SpecialPrice {
							notificationProducts = append(notificationProducts, foundProduct)
						}
					}

					s.marketplaceRepository.AddPriceToMarketplaceProduct(
						productID,
						marketplaceID,
						product.URL,
						product.RegularPrice,
						&product.SpecialPrice,
					)
				}

				if len(notificationProducts) > 0 {
					if err := s.notificationPublisher.SendNotification(&entity.NotificationTask{
						BrandID:   brandID,
						BrandName: brandName,
						Products:  notificationProducts,
					}); err != nil {
						log.Printf("could not send notification: %v", err)
					}
				}
			}
		}
	}
}

func groupProductsByBrand(products []*entity.ScrapedProduct) (map[string][]*entity.ScrapedProduct, []*entity.ScrapedProduct) {
	brandProducts := make(map[string][]*entity.ScrapedProduct)
	unknownBrandProducts := []*entity.ScrapedProduct{}
	for _, product := range products {
		if product == nil {
			continue
		}
		if product.BrandName == "" {
			unknownBrandProducts = append(unknownBrandProducts, product)
			continue
		}

		brandProducts[product.BrandName] = append(brandProducts[product.BrandName], product)
	}
	return brandProducts, unknownBrandProducts
}

func getBrandsFromProductName(unknownBrandProducts []*entity.ScrapedProduct, brandProducts map[string][]*entity.ScrapedProduct, allBrands []*entity.Brand) map[string][]*entity.ScrapedProduct {
	brands := map[string]interface{}{}

	for brandName := range brandProducts {
		brands[brandName] = nil
	}

	for _, brand := range allBrands {
		brands[brand.Name] = nil
	}

	for _, product := range unknownBrandProducts {
		for brandName := range brands {
			if strings.Contains(product.Name, brandName) {
				brandProducts[brandName] = append(brandProducts[brandName], product)
				break
			}
		}
	}

	return brandProducts
}
