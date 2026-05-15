package storage

import (
	"sales_monitor/scraper_app/feature/storage/data/repository"
	domainRepository "sales_monitor/scraper_app/feature/storage/domain/repository"
)

func NewResultStorage(folder string) domainRepository.ResultStorageRepository {
	return repository.NewFileResultStorage(folder)
}
