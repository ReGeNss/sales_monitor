package repository

import (
	"fmt"
	"sales_monitor/internal/models"
	"sales_monitor/scraper_app/core/api"
	"sales_monitor/scraper_app/shared/product/domain/entity"
	"sales_monitor/scraper_app/shared/product/domain/repository"
	"sales_monitor/scraper_app/shared/product/utils"

	fuzzy "github.com/paul-mannino/go-fuzzywuzzy"
	"gorm.io/gorm"
)

type productRepositoryImpl struct {
	db         *gorm.DB
	httpClient api.HTTPClient
}

func NewProductRepository(db *gorm.DB, httpClient api.HTTPClient) repository.ProductRepository {
	return &productRepositoryImpl{
		db:         db,
		httpClient: httpClient,
	}
}

func (p *productRepositoryImpl) CreateBrand(brand *models.Brand) (uint, error) {
	p.db.Model(&models.Brand{}).Create(brand)

	return uint(brand.BrandID), p.db.Error
}

func (p *productRepositoryImpl) CreateCategory(category *models.Category) (uint, error) {
	p.db.Model(&models.Category{}).Create(category)

	return uint(category.CategoryID), p.db.Error
}

func (p *productRepositoryImpl) CreateMarketplace(marketplace *models.Marketplace) (uint, error) {
	p.db.Model(&models.Marketplace{}).Create(marketplace)

	return uint(marketplace.MarketplaceID), p.db.Error
}

func (p *productRepositoryImpl) GetBrandByName(name string) (*models.Brand, error) {
	var brand models.Brand
	err := p.db.Model(&models.Brand{}).Where("name = ?", name).First(&brand).Error
	if err != nil {
		return nil, err
	}
	return &brand, nil
}

func (p *productRepositoryImpl) GetCategoryByName(name string) (*models.Category, error) {
	var category models.Category
	err := p.db.Model(&models.Category{}).Where("name = ?", name).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (p *productRepositoryImpl) GetMarketplaceByName(name string) (*models.Marketplace, error) {
	var marketplace models.Marketplace
	err := p.db.Model(&models.Marketplace{}).Where("name = ?", name).First(&marketplace).Error
	if err != nil {
		return nil, err
	}
	return &marketplace, nil
}

func (p *productRepositoryImpl) AddPriceToMarketplaceProduct(productID int, marketplaceID int, url string, regularPrice float64, discountPrice *float64) error {
	marketplaceProduct := models.MarketplaceProduct{
		MarketplaceID: marketplaceID,
		ProductID:     productID,
		URL:           url,
	}

	if err := p.db.Model(&models.MarketplaceProduct{}).
		Where("marketplace_id = ? AND product_id = ? AND url = ?", marketplaceID, productID, url).
		FirstOrCreate(&marketplaceProduct).Error; err != nil {
		return err
	}

	price := models.Price{
		MarketplaceProductID: marketplaceProduct.MarketplaceProductID,
		RegularPrice:         regularPrice,
		DiscountPrice:        discountPrice,
	}

	return p.db.Model(&models.Price{}).Create(&price).Error
}

func (p *productRepositoryImpl) AddPriceToMarketplaceProductID(marketplaceProductID int, regularPrice float64, discountPrice *float64) error {
	price := models.Price{
		MarketplaceProductID: marketplaceProductID,
		RegularPrice:         regularPrice,
		DiscountPrice:        discountPrice,
	}

	return p.db.Model(&models.Price{}).Create(&price).Error
}

func (p *productRepositoryImpl) CreateProduct(product *models.Product, attributes []*models.ProductAttribute) (uint, error) {
	var productAttributes []models.ProductAttribute

	err := p.db.Model(&models.Product{}).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.Product{}).Create(product).Error; err != nil {
			return err
		}

		for _, attr := range attributes {
			var existingAttr models.ProductAttribute
			err := tx.Model(&models.ProductAttribute{}).
				Where("attribute_type = ? AND value = ?", attr.AttributeType, attr.Value).
				First(&existingAttr).Error

			if err == gorm.ErrRecordNotFound {
				if err := tx.Model(&models.ProductAttribute{}).Create(attr).Error; err != nil {
					return err
				}
				productAttributes = append(productAttributes, *attr)
			} else if err != nil {
				return err
			} else {
				productAttributes = append(productAttributes, existingAttr)
			}
		}

		if len(productAttributes) > 0 {
			if err := tx.Model(product).Association("Attributes").Append(productAttributes); err != nil {
				return err
			}
		}

		return nil
	})

	return uint(product.ProductID), err
}

func (p *productRepositoryImpl) GetMostSimilarProductID(fingerprint *string, attributes []*models.ProductAttribute, productDifferentiationEntity *entity.ProductDifferentiationEntity, brandID int, categoryID int, currentMarketplaceID int) (uint, error) {
	var products []models.Product

	query := p.db.Model(&models.Product{}).Table("products as p").
		Select("p.product_id, p.name_fingerprint").
		Where("category_id = ? AND brand_id = ?", categoryID, brandID)

	query = attributesToQuery(attributes, query)

	if fingerprint == nil {
		var product models.Product
		err := query.Where("p.name_fingerprint IS NULL").First(&product).Error

		if err != nil {
			return 0, err
		}
		return uint(product.ProductID), nil
	}

	err := query.Where("MATCH(p.name_fingerprint) AGAINST(? IN NATURAL LANGUAGE MODE) > 0", *fingerprint).
		Order(fmt.Sprintf("MATCH(p.name_fingerprint) AGAINST('%s' IN NATURAL LANGUAGE MODE) DESC", *fingerprint)).
		Limit(4).
		Find(&products).Error

	if err != nil {
		return 0, err
	}

	bestSimilarity := 0
	var bestProduct models.Product

	for _, product := range products {
		similarity := fuzzy.TokenSortRatio(*fingerprint, *product.NameFingerprint)
		if similarity > bestSimilarity {
			bestSimilarity = similarity
			bestProduct = product
		}
	}

	if bestSimilarity >= 91 && utils.ProductDifferentiator(*fingerprint, *bestProduct.NameFingerprint, productDifferentiationEntity) {
		return uint(bestProduct.ProductID), nil
	}

	return 0, fmt.Errorf("no similar product found")
}

func (p *productRepositoryImpl) GetProductByFingerprint(fingerprint *string, brandID int, categoryID int, attributes []*models.ProductAttribute) (*models.Product, error) {
	var product models.Product
	query := p.db.Model(&models.Product{}).Table("products as p").Where("p.name_fingerprint = ? AND p.brand_id = ? AND p.category_id = ?", fingerprint, brandID, categoryID)
	query = attributesToQuery(attributes, query)

	if err := query.First(&product).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (p *productRepositoryImpl) CreateProductAttribute(attribute *models.ProductAttribute) error {
	return p.db.Model(&models.ProductAttribute{}).Create(attribute).Error
}

func attributesToQuery(attributes []*models.ProductAttribute, query *gorm.DB) *gorm.DB {
	for i, attribute := range attributes {
		query = query.Joins(fmt.Sprintf("JOIN product_attributes pa%[1]d ON pa%[1]d.product_id = p.product_id", i)).
			Joins(fmt.Sprintf("JOIN attributes attr%[1]d ON attr%[1]d.attribute_id = pa%[1]d.attribute_id", i)).
			Where(fmt.Sprintf("attr%[1]d.attribute_type = ? AND attr%[1]d.value = ?", i),
				attribute.AttributeType,
				attribute.Value,
			)
	}
	return query
}

func (p *productRepositoryImpl) GetAllBrands() ([]models.Brand, error) {
	var brands []models.Brand
	err := p.db.Model(&models.Brand{}).Find(&brands).Error
	if err != nil {
		return nil, err
	}
	return brands, nil
}

func (p *productRepositoryImpl) GetLaterScrapedProducts(brandID int) (entity.LaterScrapedProductsUrls, error) {
	var marketplaceProducts []models.MarketplaceProduct
	err := p.db.Model(&models.MarketplaceProduct{}).Table("marketplace_products as mp").
		Joins("JOIN products p ON p.product_id = mp.product_id").
		Where("p.brand_id = ?", brandID).
		Find(&marketplaceProducts).Error
	if err != nil {
		return nil, err
	}

	laterScrapedProductsUrls := make(entity.LaterScrapedProductsUrls)
	for _, marketplaceProduct := range marketplaceProducts {
		laterScrapedProductsUrls[marketplaceProduct.URL] = marketplaceProduct.MarketplaceProductID
	}
	return laterScrapedProductsUrls, nil
}
