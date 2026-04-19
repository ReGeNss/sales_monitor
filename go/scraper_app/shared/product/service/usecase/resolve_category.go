package usecase

import (
	"sales_monitor/scraper_app/shared/product/domain/entity"
	"sales_monitor/scraper_app/shared/product/domain/repository"
)

type ResolveCategoryUseCase interface {
	Execute(name string) (int, error)
}

type resolveCategoryUseCase struct {
	repo repository.CategoryRepository
}

func NewResolveCategoryUseCase(repo repository.CategoryRepository) ResolveCategoryUseCase {
	return &resolveCategoryUseCase{repo: repo}
}

func (u *resolveCategoryUseCase) Execute(name string) (int, error) {
	category, err := u.repo.GetCategoryByName(name)
	if err == nil {
		return category.ID, nil
	}

	category = &entity.Category{Name: name}
	if _, err := u.repo.CreateCategory(category); err != nil {
		return 0, err
	}
	return category.ID, nil
}
