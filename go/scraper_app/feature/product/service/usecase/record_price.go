package usecase

import (
	"sales_monitor/scraper_app/feature/product/domain/entity"
	"sales_monitor/scraper_app/feature/product/domain/exception"
	"sales_monitor/scraper_app/feature/product/domain/repository"
)

type RecordPriceInput struct {
	Scraped              *entity.ScrapedProduct
	ProductID            int
	MarketplaceID        int
	ExistingProduct      *entity.Product
	KnownMarketplaceURLs entity.LaterScrapedProductsUrls
}

type RecordPriceUseCase interface {
	Execute(input RecordPriceInput) (priceDrop *entity.Product, err exception.IDomainError)
}

type recordPriceUseCase struct {
	marketplaceRepository repository.MarketplaceRepository
	priceRepository       repository.PriceRepository
}

func NewRecordPriceUseCase(
	marketplaceRepository repository.MarketplaceRepository,
	priceRepository repository.PriceRepository,
) RecordPriceUseCase {
	return &recordPriceUseCase{
		marketplaceRepository: marketplaceRepository,
		priceRepository:       priceRepository,
	}
}

func (u *recordPriceUseCase) Execute(input RecordPriceInput) (*entity.Product, exception.IDomainError) {
	if marketplaceProductID, ok := input.KnownMarketplaceURLs[input.Scraped.URL()]; ok {
		knownSpecialPrice := input.Scraped.SpecialPrice()
		u.marketplaceRepository.AddPriceToMarketplaceProductID(
			marketplaceProductID,
			input.Scraped.RegularPrice(),
			&knownSpecialPrice,
		)
	}

	priceDrop := u.detectPriceDrop(input)

	specialPrice := input.Scraped.SpecialPrice()
	if err := u.marketplaceRepository.AddPriceToMarketplaceProduct(
		input.ProductID,
		input.MarketplaceID,
		input.Scraped.URL(),
		input.Scraped.RegularPrice(),
		&specialPrice,
	); err != nil {
		return nil, err
	}
	return priceDrop, nil
}

func (u *recordPriceUseCase) detectPriceDrop(input RecordPriceInput) *entity.Product {
	if input.ExistingProduct == nil {
		return nil
	}
	latest, err := u.priceRepository.GetLatestProductPrice(input.ProductID)
	if err != nil || latest == nil || latest.SpecialPrice == nil {
		return nil
	}
	if *latest.SpecialPrice < input.Scraped.SpecialPrice() {
		return input.ExistingProduct
	}
	return nil
}
