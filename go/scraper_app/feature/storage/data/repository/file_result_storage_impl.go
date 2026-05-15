package repository

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"sales_monitor/scraper_app/feature/storage/domain/repository"
)

type fileResultStorage struct {
	folder string
}

func NewFileResultStorage(folder string) repository.ResultStorageRepository {
	return &fileResultStorage{folder: folder}
}

func (s *fileResultStorage) Save(payload any, categories []string) {
	filename := fmt.Sprintf("%s/scraped_%s_%s.json",
		s.folder,
		strings.Join(categories, " "),
		time.Now().Format(time.DateTime),
	)

	if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
		log.Printf("could not create directory: %v", err)
		return
	}

	file, err := os.Create(filename)
	if err != nil {
		log.Printf("could not create file: %v", err)
		return
	}
	defer file.Close()

	if err := json.NewEncoder(file).Encode(payload); err != nil {
		log.Printf("could not save scraping result: %v", err)
	}
}
