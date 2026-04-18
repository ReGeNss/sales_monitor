package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"sales_monitor/internal/models"
	"sales_monitor/scraper_app/core/env"
	"sales_monitor/scraper_app/shared/product/data/mapper"
	"sales_monitor/scraper_app/shared/product/domain/entity"
	"sales_monitor/scraper_app/shared/product/domain/repository"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type productRepositoryImpl struct {
	db          *gorm.DB
	redisClient *redis.Client
}

func (p *productRepositoryImpl) GetLatestProductPrice(productID int) (*entity.Price, error) {
	var price models.Price
	err := p.db.Model(&models.Price{}).
		Joins("JOIN marketplace_products mp ON mp.marketplace_product_id = prices.marketplace_product_id").
		Where("mp.product_id = ?", productID).
		Order("prices.created_at DESC, prices.price_id DESC").
		First(&price).Error
	if err != nil {
		return nil, err
	}
	return mapper.PriceToEntity(&price), nil
}

func (p *productRepositoryImpl) SendNotification(notificationTask *entity.NotificationTask) error {
	if notificationTask == nil {
		return fmt.Errorf("notification task is nil")
	}

	if p.redisClient == nil {
		return fmt.Errorf("redis client is not initialized")
	}

	payload, err := json.Marshal(mapper.NotificationTaskToModel(notificationTask))
	if err != nil {
		return err
	}
	return p.redisClient.LPush(context.Background(), env.GetNotificationQueueKey(), string(payload)).Err()
}

func NewProductRepository(db *gorm.DB, redisClient *redis.Client) repository.ProductRepository {
	return &productRepositoryImpl{
		db:          db,
		redisClient: redisClient,
	}
}

func (p *productRepositoryImpl) CreateBrand(brand *entity.Brand) (uint, error) {
	m := mapper.BrandToModel(brand)
	if err := p.db.Model(&models.Brand{}).Create(m).Error; err != nil {
		return 0, err
	}
	brand.ID = m.BrandID
	return uint(m.BrandID), nil
}

func (p *productRepositoryImpl) CreateCategory(category *entity.Category) (uint, error) {
	m := mapper.CategoryToModel(category)
	if err := p.db.Model(&models.Category{}).Create(m).Error; err != nil {
		return 0, err
	}
	category.ID = m.CategoryID
	return uint(m.CategoryID), nil
}

func (p *productRepositoryImpl) CreateMarketplace(marketplace *entity.Marketplace) (uint, error) {
	m := mapper.MarketplaceToModel(marketplace)
	if err := p.db.Model(&models.Marketplace{}).Create(m).Error; err != nil {
		return 0, err
	}
	marketplace.ID = m.MarketplaceID
	return uint(m.MarketplaceID), nil
}

func (p *productRepositoryImpl) GetBrandByName(name string) (*entity.Brand, error) {
	var brand models.Brand
	err := p.db.Model(&models.Brand{}).Where("name = ?", name).First(&brand).Error
	if err != nil {
		return nil, err
	}
	return mapper.BrandToEntity(&brand), nil
}

func (p *productRepositoryImpl) GetCategoryByName(name string) (*entity.Category, error) {
	var category models.Category
	err := p.db.Model(&models.Category{}).Where("name = ?", name).First(&category).Error
	if err != nil {
		return nil, err
	}
	return mapper.CategoryToEntity(&category), nil
}

func (p *productRepositoryImpl) GetMarketplaceByName(name string) (*entity.Marketplace, error) {
	var marketplace models.Marketplace
	err := p.db.Model(&models.Marketplace{}).Where("name = ?", name).First(&marketplace).Error
	if err != nil {
		return nil, err
	}
	return mapper.MarketplaceToEntity(&marketplace), nil
}

func (p *productRepositoryImpl) AddPriceToMarketplaceProduct(productID int, marketplaceID int, url string, regularPrice float64, specialPrice *float64) error {
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
		SpecialPrice:         specialPrice,
	}

	return p.db.Model(&models.Price{}).Create(&price).Error
}

func (p *productRepositoryImpl) AddPriceToMarketplaceProductID(marketplaceProductID int, regularPrice float64, specialPrice *float64) error {
	price := models.Price{
		MarketplaceProductID: marketplaceProductID,
		RegularPrice:         regularPrice,
		SpecialPrice:         specialPrice,
	}

	return p.db.Model(&models.Price{}).Create(&price).Error
}

func (p *productRepositoryImpl) CreateProduct(product *entity.Product, attributes []*entity.ProductAttribute) (uint, error) {
	productModel := mapper.ProductToModel(product)
	var productAttributes []models.ProductAttribute

	err := p.db.Model(&models.Product{}).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.Product{}).Create(productModel).Error; err != nil {
			return err
		}

		for _, attr := range attributes {
			attrModel := mapper.ProductAttributeToModel(attr)
			var existingAttr models.ProductAttribute
			err := tx.Model(&models.ProductAttribute{}).
				Where("attribute_type = ? AND value = ?", attrModel.AttributeType, attrModel.Value).
				First(&existingAttr).Error

			if err == gorm.ErrRecordNotFound {
				if err := tx.Model(&models.ProductAttribute{}).Create(attrModel).Error; err != nil {
					return err
				}
				attr.ID = attrModel.AttributeID
				productAttributes = append(productAttributes, *attrModel)
			} else if err != nil {
				return err
			} else {
				attr.ID = existingAttr.AttributeID
				productAttributes = append(productAttributes, existingAttr)
			}
		}

		if len(productAttributes) > 0 {
			if err := tx.Model(productModel).Association("Attributes").Append(productAttributes); err != nil {
				return err
			}
		}

		return nil
	})

	product.ID = productModel.ProductID
	return uint(productModel.ProductID), err
}

func (p *productRepositoryImpl) FindSimilarCandidates(fingerprint *string, attributes []*entity.ProductAttribute, brandID int, categoryID int) ([]*entity.Product, error) {
	query := p.db.Model(&models.Product{}).Table("products as p").
		Select("p.product_id, p.name_fingerprint").
		Where("category_id = ? AND brand_id = ?", categoryID, brandID)

	query = attributesToQuery(attributes, query)

	if fingerprint == nil {
		var product models.Product
		err := query.Where("p.name_fingerprint IS NULL").First(&product).Error
		if err != nil {
			return nil, err
		}
		return []*entity.Product{mapper.ProductToEntity(&product)}, nil
	}

	var products []models.Product
	err := query.Where("MATCH(p.name_fingerprint) AGAINST(? IN NATURAL LANGUAGE MODE) > 0", *fingerprint).
		Order(fmt.Sprintf("MATCH(p.name_fingerprint) AGAINST('%s' IN NATURAL LANGUAGE MODE) DESC", *fingerprint)).
		Limit(4).
		Find(&products).Error
	if err != nil {
		return nil, err
	}

	candidates := make([]*entity.Product, 0, len(products))
	for i := range products {
		candidates = append(candidates, mapper.ProductToEntity(&products[i]))
	}
	return candidates, nil
}

func (p *productRepositoryImpl) GetProductByFingerprint(fingerprint *string, brandID int, categoryID int, attributes []*entity.ProductAttribute) (*entity.Product, error) {
	var product models.Product
	query := p.db.Model(&models.Product{}).Table("products as p").Where("p.name_fingerprint = ? AND p.brand_id = ? AND p.category_id = ?", fingerprint, brandID, categoryID)
	query = attributesToQuery(attributes, query)

	if err := query.First(&product).Error; err != nil {
		return nil, err
	}
	return mapper.ProductToEntity(&product), nil
}

func (p *productRepositoryImpl) CreateProductAttribute(attribute *entity.ProductAttribute) error {
	m := mapper.ProductAttributeToModel(attribute)
	if err := p.db.Model(&models.ProductAttribute{}).Create(m).Error; err != nil {
		return err
	}
	attribute.ID = m.AttributeID
	return nil
}

func attributesToQuery(attributes []*entity.ProductAttribute, query *gorm.DB) *gorm.DB {
	for i, attribute := range attributes {
		query = query.Joins(fmt.Sprintf("JOIN product_attributes pa%[1]d ON pa%[1]d.product_id = p.product_id", i)).
			Joins(fmt.Sprintf("JOIN attributes attr%[1]d ON attr%[1]d.attribute_id = pa%[1]d.attribute_id", i)).
			Where(fmt.Sprintf("attr%[1]d.attribute_type = ? AND attr%[1]d.value = ?", i),
				attribute.Type,
				attribute.Value,
			)
	}
	return query
}

func (p *productRepositoryImpl) GetAllBrands() ([]*entity.Brand, error) {
	var brands []models.Brand
	err := p.db.Model(&models.Brand{}).Find(&brands).Error
	if err != nil {
		return nil, err
	}
	return mapper.BrandsToEntities(brands), nil
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
