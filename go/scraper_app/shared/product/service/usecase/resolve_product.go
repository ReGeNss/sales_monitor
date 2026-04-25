package usecase

import (
	"sales_monitor/scraper_app/shared/product/domain/entity"
	"sales_monitor/scraper_app/shared/product/domain/exception"
	"sales_monitor/scraper_app/shared/product/domain/repository"
	"sales_monitor/scraper_app/shared/product/domain/service"
)

type ResolveProductInput struct {
	Scraped         *entity.ScrapedProduct
	Fingerprint     *string
	Attributes      []*entity.ProductAttribute
	BrandID         int
	CategoryID      int
	Differentiation *entity.ProductDifferentiationEntity
}

type ResolveProductUseCase interface {
	Execute(input ResolveProductInput) (productID int, existing *entity.Product, err exception.IDomainError)
}

type resolveProductUseCase struct {
	productRepository repository.ProductRepository
	matcher           service.ProductMatcher
}

func NewResolveProductUseCase(productRepository repository.ProductRepository, matcher service.ProductMatcher) ResolveProductUseCase {
	return &resolveProductUseCase{productRepository: productRepository, matcher: matcher}
}

func (u *resolveProductUseCase) Execute(input ResolveProductInput) (int, *entity.Product, exception.IDomainError) {
	if existing, err := u.productRepository.GetProductByFingerprint(input.Fingerprint, input.BrandID, input.CategoryID, input.Attributes); err == nil {
		return existing.ID, existing, nil
	}

	if matched := u.pickMatch(input); matched != nil {
		return matched.ID, matched, nil
	}

	id, err := u.productRepository.CreateProduct(&entity.Product{
		Name:            input.Scraped.Name,
		NameFingerprint: input.Fingerprint,
		ImageURL:        input.Scraped.Image,
		BrandID:         input.BrandID,
		CategoryID:      input.CategoryID,
	}, input.Attributes)
	if err != nil {
		return 0, nil, err
	}
	return int(id), nil, nil
}

func (u *resolveProductUseCase) pickMatch(input ResolveProductInput) *entity.Product {
	candidates, err := u.productRepository.FindSimilarCandidates(input.Fingerprint, input.Attributes, input.BrandID, input.CategoryID)
	if err != nil {
		return nil
	}

	if input.Fingerprint == nil {
		if len(candidates) > 0 {
			return candidates[0]
		}
		return nil
	}

	if matched, ok := u.matcher.PickBestMatch(*input.Fingerprint, candidates, input.Differentiation); ok {
		return matched
	}
	return nil
}
