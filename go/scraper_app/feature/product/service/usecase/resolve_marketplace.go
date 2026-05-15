package usecase

import (
	"sales_monitor/scraper_app/feature/product/domain/entity"
	"sales_monitor/scraper_app/feature/product/domain/exception"
	"sales_monitor/scraper_app/feature/product/domain/repository"
)

type ResolveMarketplaceUseCase interface {
	Execute(name string) (int, exception.IDomainError)
}

type resolveMarketplaceUseCase struct {
	repo repository.MarketplaceRepository
}

func NewResolveMarketplaceUseCase(repo repository.MarketplaceRepository) ResolveMarketplaceUseCase {
	return &resolveMarketplaceUseCase{repo: repo}
}

func (u *resolveMarketplaceUseCase) Execute(name string) (int, exception.IDomainError) {
	marketplace, err := u.repo.GetMarketplaceByName(name)
	if err == nil {
		return marketplace.ID, nil
	}

	marketplace = &entity.Marketplace{Name: name}
	if _, err := u.repo.CreateMarketplace(marketplace); err != nil {
		return 0, err
	}
	return marketplace.ID, nil
}
