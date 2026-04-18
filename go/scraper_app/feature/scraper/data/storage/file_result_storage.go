package storage

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sales_monitor/scraper_app/feature/scraper/domain/entity"
	"sales_monitor/scraper_app/feature/scraper/domain/gateway"
	"strings"
	"time"
)

type fileResultStorage struct {
	folder string
}

func NewFileResultStorage(folder string) gateway.ResultStorage {
	return &fileResultStorage{folder: folder}
}

func (s *fileResultStorage) Save(results map[string]*entity.ScrapingResult, categories []string) {
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
	err = json.NewEncoder(file).Encode(&results)
	if err != nil {
		log.Printf("could not save scraping result: %v", err)
	}
}
