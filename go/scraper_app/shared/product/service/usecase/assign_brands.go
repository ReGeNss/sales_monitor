package usecase

import (
	"sales_monitor/scraper_app/shared/product/domain/entity"
	"sales_monitor/scraper_app/shared/product/domain/repository"
	"strings"
)

type AssignBrandsUseCase interface {
	Execute(products []*entity.ScrapedProduct) (map[string][]*entity.ScrapedProduct, error)
}

type assignBrandsUseCase struct {
	brandRepository repository.BrandRepository
}

func NewAssignBrandsUseCase(brandRepository repository.BrandRepository) AssignBrandsUseCase {
	return &assignBrandsUseCase{brandRepository: brandRepository}
}

func (u *assignBrandsUseCase) Execute(products []*entity.ScrapedProduct) (map[string][]*entity.ScrapedProduct, error) {
	grouped, unknown := groupByBrand(products)

	if len(unknown) == 0 {
		return grouped, nil
	}

	allBrands, err := u.brandRepository.GetAllBrands()
	if err != nil {
		return nil, err
	}

	return matchUnknownByName(unknown, grouped, allBrands), nil
}

func groupByBrand(products []*entity.ScrapedProduct) (map[string][]*entity.ScrapedProduct, []*entity.ScrapedProduct) {
	grouped := make(map[string][]*entity.ScrapedProduct)
	unknown := []*entity.ScrapedProduct{}
	for _, product := range products {
		if product == nil {
			continue
		}
		if product.BrandName == "" {
			unknown = append(unknown, product)
			continue
		}
		grouped[product.BrandName] = append(grouped[product.BrandName], product)
	}
	return grouped, unknown
}

func matchUnknownByName(unknown []*entity.ScrapedProduct, grouped map[string][]*entity.ScrapedProduct, allBrands []*entity.Brand) map[string][]*entity.ScrapedProduct {
	known := map[string]struct{}{}
	for brandName := range grouped {
		known[brandName] = struct{}{}
	}
	for _, brand := range allBrands {
		known[brand.Name] = struct{}{}
	}

	for _, product := range unknown {
		for brandName := range known {
			if strings.Contains(product.Name, brandName) {
				grouped[brandName] = append(grouped[brandName], product)
				break
			}
		}
	}
	return grouped
}
