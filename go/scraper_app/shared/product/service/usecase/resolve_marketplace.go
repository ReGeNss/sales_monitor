package usecase

import (
	"sales_monitor/scraper_app/shared/product/domain/entity"
	"sales_monitor/scraper_app/shared/product/domain/repository"
)

type ResolveMarketplaceUseCase interface {
	Execute(name string) (int, error)
}

type resolveMarketplaceUseCase struct {
	repo repository.MarketplaceRepository
}

func NewResolveMarketplaceUseCase(repo repository.MarketplaceRepository) ResolveMarketplaceUseCase {
	return &resolveMarketplaceUseCase{repo: repo}
}

func (u *resolveMarketplaceUseCase) Execute(name string) (int, error) {
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
