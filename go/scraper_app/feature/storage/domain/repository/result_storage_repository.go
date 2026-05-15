package repository

type ResultStorageRepository interface {
	Save(payload any, categories []string)
}
