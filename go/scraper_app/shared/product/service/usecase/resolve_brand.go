package usecase

import (
	"sales_monitor/scraper_app/shared/product/domain/entity"
	"sales_monitor/scraper_app/shared/product/domain/exception"
	"sales_monitor/scraper_app/shared/product/domain/repository"
)

type ResolveBrandUseCase interface {
	Execute(name string) (int, exception.IDomainError)
}

type resolveBrandUseCase struct {
	repo repository.BrandRepository
}

func NewResolveBrandUseCase(repo repository.BrandRepository) ResolveBrandUseCase {
	return &resolveBrandUseCase{repo: repo}
}

func (u *resolveBrandUseCase) Execute(name string) (int, exception.IDomainError) {
	brand, err := u.repo.GetBrandByName(name)
	if err == nil {
		return brand.ID, nil
	}

	id, err := u.repo.CreateBrand(&entity.Brand{Name: name})
	if err != nil {
		return 0, err
	}
	return int(id), nil
}
