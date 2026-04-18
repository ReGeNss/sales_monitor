package utils

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

func SaveToJsonFile(data interface{}, filename string) error {
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Printf("Error creating directory: %v", err)
		return err
	}

	file, err := os.Create(filename)
	if err != nil {
		log.Printf("Error creating file: %v", err)
		return err
	}
	defer file.Close()
	return json.NewEncoder(file).Encode(data)
}
