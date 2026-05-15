package repository

import (
	"fmt"
	"sales_monitor/internal/models"
	"sales_monitor/scraper_app/feature/product/data/mapper"
	"sales_monitor/scraper_app/feature/product/domain/entity"
	"sales_monitor/scraper_app/feature/product/domain/exception"
	"sales_monitor/scraper_app/feature/product/domain/repository"

	"gorm.io/gorm"
)

type productRepositoryImpl struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) repository.ProductRepository {
	return &productRepositoryImpl{db: db}
}

func (p *productRepositoryImpl) GetProductByFingerprint(fingerprint *string, brandID int, categoryID int, attributes []*entity.ProductAttribute) (*entity.Product, exception.IDomainError) {
	var product models.Product
	query := p.db.Model(&models.Product{}).Table("products as p").
		Where("p.name_fingerprint = ? AND p.brand_id = ? AND p.category_id = ?", fingerprint, brandID, categoryID)
	query = attributesToQuery(attributes, query)

	if err := query.First(&product).Error; err != nil {
		return nil, err
	}
	entity, err := mapper.ProductToEntity(&product)

	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (p *productRepositoryImpl) FindSimilarCandidates(fingerprint *string, attributes []*entity.ProductAttribute, brandID int, categoryID int) ([]*entity.Product, exception.IDomainError) {
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
		e, err := mapper.ProductToEntity(&product)
		if err != nil {
			return nil, err
		}
		return []*entity.Product{e}, nil
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
		e, err := mapper.ProductToEntity(&products[i])
		if err != nil {
			return nil, err
		}
		candidates = append(candidates, e)
	}
	return candidates, nil
}

func (p *productRepositoryImpl) CreateProduct(product *entity.Product, attributes []*entity.ProductAttribute) (uint, exception.IDomainError) {
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

func (p *productRepositoryImpl) CreateProductAttribute(attribute *entity.ProductAttribute) exception.IDomainError {
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
