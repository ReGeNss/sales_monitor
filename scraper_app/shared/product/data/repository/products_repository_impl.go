package repository

import (
	"fmt"
	"sales_monitor/internal/models"
	"sales_monitor/scraper_app/shared/product/domain/repository"
	fuzzy "github.com/paul-mannino/go-fuzzywuzzy"
	"gorm.io/gorm"
)

type productRepositoryImpl struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) repository.ProductRepository {
	return &productRepositoryImpl{
		db: db,
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

func (p *productRepositoryImpl) AddPriceToProduct(price *models.Price) error {
	return p.db.Model(&models.Price{}).Create(price).Error
}

func (p *productRepositoryImpl) CreateProduct(product *models.Product, attributes []*models.ProductAttribute) (uint, error) {
	p.db.Model(&models.Product{}).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.Product{}).Create(product).Error; err != nil {
			return err
		}

		for _, attribute := range attributes {
			attribute.ProductID = product.ProductID
			if err := tx.Model(&models.ProductAttribute{}).Create(attribute).Error; err != nil {
				return err
			}
		}

		return nil
	})

	return uint(product.ProductID), p.db.Error
}

func (p *productRepositoryImpl) GetMostSimilarProductID(fingerprint string, attributes []*models.ProductAttribute, brandID int, categoryID int) (uint, error) {
	var products []models.Product

	query := p.db.Model(&models.Product{}).Table("Product as p").
		Select("p.product_id, p.name_fingerprint").
		Where("category_id = ? AND brand_id = ?", categoryID, brandID)

	query = attributesToQuery(attributes, query)

	err := query.Where("MATCH(p.name_fingerprint) AGAINST(? IN NATURAL LANGUAGE MODE) > 0", fingerprint).
		Order(fmt.Sprintf("MATCH(p.name_fingerprint) AGAINST('%s' IN NATURAL LANGUAGE MODE) DESC", fingerprint)).
		Limit(4).
		Find(&products).Error

	if err != nil {
		return 0, err
	}

	bestSimilarity := 0
	bestProductID := 0

	for _, product := range products {
		similarity := fuzzy.Ratio(fingerprint, product.NameFingerprint)
		if similarity > bestSimilarity {
			bestSimilarity = similarity
			bestProductID = product.ProductID
		}
	}

	if bestSimilarity >= 95 {
		return uint(bestProductID), nil
	}

	return 0, fmt.Errorf("no similar product found")
}

func (p *productRepositoryImpl) GetProductByFingerprint(fingerprint string, brandID int, categoryID int, attributes []*models.ProductAttribute) (*models.Product, error) {
	var product models.Product
	query := p.db.Model(&models.Product{}).Table("Product as p").Where("p.name_fingerprint = ? AND p.brand_id = ? AND p.category_id = ?", fingerprint, brandID, categoryID)
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
		query = query.Joins(fmt.Sprintf("JOIN product_attributes attr%[1]d ON attr%[1]d.product_id = p.product_id", i)).
			Where(fmt.Sprintf("attr%[1]d.attribute_type = ? AND attr%[1]d.value = ?", i),
				attribute.AttributeType,
				attribute.Value,
			)
	}
	return query
}